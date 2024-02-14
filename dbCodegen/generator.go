package dbCodegen

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/techrail/ground/typs/set"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/gertd/go-pluralize"

	"github.com/techrail/ground/typs/appError"
)

const noSuchEnumErrCode = "1RD8QB"

var goKeywords []string

func init() {
	goKeywords = []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var"}
}

const DefaultMagicComment = "// MAGIC COMMENT (DO NOT EDIT): Please write any custom code only below this line.\n"

// DbSchema represents the schema in the database
type DbSchema struct {
	Name      string             // Name of the schema
	GoName    string             // The name of table as a go variable that we would use
	Tables    map[string]DbTable // Map of name of table to their DbTable struct values
	TablesA2z []string           // List of Tables in alphabetical order
}

// DbTable represents a table in the database
type DbTable struct {
	Name           string                 // Table name as got from the query
	GoName         string                 // The name of table as a go variable that we would use
	GoNameSingular string                 // Singular form of GoName
	GoNamePlural   string                 // Plural form of GoName
	Schema         string                 // The name of the schema where this table resides
	Comment        string                 // Comment on the table
	ColumnMap      map[string]DbColumn    // List of columns as a map from the name of the column to the DbColumn type
	ColumnList     []string               // List of columns (ordinal)
	ColumnListA2z  []string               // List of column names (alphabetical)
	PkColumnList   []DbColumn             // List of columns that make the primary key (slice because order matters)
	IndexList      []DbIndex              // List of indexes on this table
	FKeyMap        map[string]DbFkInfo    // List of foreign keys in table as map from constraint name to DbFkInfo type
	RevFKeyMap     map[string]DbRevFkInfo // List of reverse reference in table as map from constraint name to DbRevFkInfo type
}

func (table *DbTable) fullyQualifiedTableName() string {
	return table.Schema + "." + table.Name
}

func (table *DbTable) fullyQualifiedStructName() string {
	// return getGoName(table.Schema) + "_" + table.GoNameSingular
	return getGoName(table.Schema) + table.GoNameSingular
}

func (table *DbTable) fullyQualifiedDaoName() string {
	return table.fullyQualifiedStructName() + "Dao"
}

func (table *DbTable) fullyQualifiedVariableName() string {
	return lowerFirstChar(getGoName(table.Schema) + "_" + table.GoNameSingular)
}

func (table *DbTable) variableName() string {
	return lowerFirstChar(table.GoNameSingular)
}

func (table *DbTable) commentForStruct() string {
	return strings.ReplaceAll(table.Comment, "\n", "\n// ")
}

func (table *DbTable) isColumnPrimaryKey(input string) bool {
	for _, col := range table.PkColumnList {
		if col.Name == input {
			return true
		}
	}
	return false
}

func (table *DbTable) FindIndexByColumnNames(colNames []string) DbIndex {
	slices.Sort(colNames)

	for _, idx := range table.IndexList {
		// Check if the index has the same number of columns or not
		if len(idx.ColumnList) == len(colNames) {
			// Is this the index we are looking for?
			thisIsTheIndex := true
			// Now check that the columns mentioned in the input match with the ones in this index
			sortedIdxColNames := idx.GetSortedColumnNames()
			for i := 0; i < len(idx.ColumnList); i++ {
				if colNames[i] != sortedIdxColNames[i] {
					// At least this is not the index
					thisIsTheIndex = false
					break
				}
			}

			if thisIsTheIndex {
				return idx
			}
		}
	}

	return DbIndex{}
}

// DbColumn is the column representation of a table in the database for the generator
type DbColumn struct {
	Schema            string           // Schema name in which this column's table resides
	Table             string           // Table name of the table in which this column is
	Name              string           // Column name
	GoName            string           // Name we want to use for Golang code that will be generated
	GoNameSingular    string           // Singular form of the name
	GoNamePlural      string           // Plural form of the name
	DataType          string           // Data type we get from db
	GoDataType        string           // Data type we want to use in go program
	NetworkDataType   string           // Data type we want to use for the network model
	Comment           string           // Column comment
	CharacterLength   int              // Length in case it is varchar
	Nullable          bool             // NOT NULL means it is false
	HasDefaultValue   bool             // Does the column have a default value?
	DefaultValue      string           // If column has default value then what is it
	CommentProperties dbColumnProperty // Properties that will control mostly column value validations
}

func (col *DbColumn) newlineEscapedComment() string {
	return strings.ReplaceAll(col.Comment, "\n", " (nwln) ")
}

func (col *DbColumn) fullyQualifiedColumnName() string {
	return col.Schema + "." + col.Table + "." + col.Name
}

func (col *DbColumn) asSafeVariableName() string {
	n := lowerFirstChar(col.GoNameSingular)
	// Should not be a keyword
	if isGoKeyword(n) {
		// repeat the last character
		return n + n[len(n)-1:]
	}
	return n
}

func isGoKeyword(word string) bool {
	isKeyword := false
	for _, w := range goKeywords {
		if w == word {
			isKeyword = true
			break
		}
	}
	return isKeyword
}

// Properties that would be expressed as json in the column comment **after** the actual comment
// For example, the comment on a `email` column with `(^_^)` as the comment separator can be something like this:
//
//	Email address of the user (^_^) {"minStrLen":8}
//
// It should generate the validator which would ensure that the email field at least has 8 characters in it!
type dbColumnProperty struct {
	Disabled             bool      `json:"disabled"`             // Makes generator behave as if no column property set
	MinStrLen            int       `json:"minStrLen"`            // for string columns - minimum length
	MaxStrLen            int       `json:"maxStrLen"`            // for string columns - maximum length
	MinIntVal            int64     `json:"minIntVal"`            // for integer columns - minimum value
	MaxIntVal            int64     `json:"maxIntVal"`            // for integer columns - maximum value
	MinTimestampVal      time.Time `json:"minTimestampVal"`      // For timestamp without time zone columns
	MaxTimestampVal      time.Time `json:"maxTimestampVal"`      // For timestamp without time zone columns
	StrValidateAs        string    `json:"strValidateAs"`        // Validate String Data as what? (Email? URL? Name? Regex?)
	HideFromNetwork      bool      `json:"hideFromNetwork"`      // Should this field be hidden in network response
	StrConversionViaEnum string    `json:"strConversionViaEnum"` // To be used for enumerated fields that need to be represented as string in network responses. The enums must be one being generated.
	// StrConversionViaType string    `json:"strConversionViaType"` // To be used for enumerated fields that need to be represented as string in network responses. The type must be in the
}

// DbIndex represents an index inside a table
type DbIndex struct {
	Name       string     // Name of the index
	IsUnique   bool       // Is this a "unique" index?
	IsPrimary  bool       // Does this index correspond to the primary key of the table
	ColumnList []DbColumn // List of columns of this index (ordered the same way as the columns are defined in the index)
}

func (idx *DbIndex) GetFuncNamePart() string {
	funcNamePart := ""
	for i, col := range idx.ColumnList {
		if i == 0 {
			funcNamePart = col.GoName
		} else {
			funcNamePart = funcNamePart + col.GoNameSingular
		}
	}
	return funcNamePart
}

func (idx *DbIndex) GetSortedColumnNames() []string {
	colNames := []string{}
	for _, col := range idx.ColumnList {
		colNames = append(colNames, col.Name)
	}
	slices.Sort(colNames)
	return colNames
}

// DbFkInfo represents a single foreign key in a table
type DbFkInfo struct {
	FromSchema     string            // The schema name of the table from which the foreign key reference is being made
	FromTable      string            // The table name which is referencing the target table
	ToSchema       string            // The schema name of the table whose column is being referenced
	ToTable        string            // The table name of the table which is being referenced
	FromColOrder   []string          // The order which the columns appear in the FromTable
	References     map[string]string // The reference map ([from_column]to_column format)
	ConstraintName string            // Name of the foreign key constraint
}

func (fki *DbFkInfo) GetReverseRefName() string {
	return fki.FromSchema + "." + fki.FromTable + "." + fki.ConstraintName
}

type DbRevFkInfo struct {
	DbFkInfo
	UniqueIndex bool // Is there a unique index on the column set pointing to this column
}

type fkInfoFromDb struct {
	FromSchema      string `db:"from_schema"`
	FromTable       string `db:"from_table"`
	FromColumn      string `db:"from_column"`
	ToSchema        string `db:"to_schema"`
	ToTable         string `db:"to_table"`
	ToColumn        string `db:"to_column"`
	OrdinalPosition int    `db:"ordinal_position"`
	ConstraintName  string `db:"constraint_name"`
}

// EnumDefinition defines an enumeration in code which would ideally be saved in the DB
type EnumDefinition struct {
	Name              string           // Name of this enum
	Exported          bool             // Enum to be used outside the DB package
	IsDbType          bool             // Is this enum supposed to be used in the DB?
	Mappings          map[string]int16 // List of enumerations
	DisableGeneration bool             // Disable the generation/update of this type (temporarily?)
	goName            string           // Enum name for use in go code
	goNameSingular    string           // Enum name in singular form for go code
	goNamePlural      string           // Enum name in Plural form for go code
	goTypeName        string           // Enum type name form for go struct
}

// CodegenConfig contains the values and rules using which the code is to be generated
type CodegenConfig struct {
	PgDbUrl                  string // DB URL string for PostgreSQL database to which we have to connect
	ModelsContainerPackage   string // Full package name under which the DB and Network packages will fall
	DbModelPackageName       string // Name of the package under which the db related code will be placed
	DbModelPackagePath       string // Full path of the directory where the generated code for db will be placed
	NetworkPackageName       string // Name of the package under which the networking related code for DB tables is gonna be placed
	NetworkPackagePath       string // Full path of the directory where the networking related code for DB tables is gonna be placed
	MagicComment             string // Magic comment which allows the generator to update only the generated portion of code
	ColCommentSeparator      string // The string after which we can place the Properties JSON in column comment
	InsertCreatedAtInCode    bool   // Should the code for inserting `created_at` timestamp column be generated by the code generator?
	InsertUpdatedAtInCode    bool   // Should the code for inserting `updated_at` timestamp column be generated by the code generator?
	UpdateUpdatedAtInCode    bool   // Should the code for updating `updated_at` timestamp column be generated by the code generator?
	BuildUpdateByUniqueIndex bool   // Should we generate the update functions for unique indexes?
	ColumnOrderAlphabetic    bool   // Column order in generated code will be alphabetic if this is set to true, ordinal otherwise
	Enumerations             map[string]EnumDefinition
}

// Generator is the structure we return to a client which needs a generator.
// It is supposed to contain all the information needed for performing the code generation
type Generator struct {
	// More fields to be decided
	Config       CodegenConfig             // The configuration for this generator
	Schemas      map[string]DbSchema       // The schemas in the database (will in turn contain tables)
	Enums        map[string]EnumDefinition // The enumerations to be built
	pluralClient *pluralize.Client         // Pluralization client
	sync.Mutex                             // To prevent parallel runs
}

// The type to get the column info for all the tables in all the schemas
type rawCol struct {
	Schema         sql.NullString `db:"table_schema"`
	TableName      sql.NullString `db:"table_name"`
	TableComment   sql.NullString `db:"table_comment"`
	ColumnName     sql.NullString `db:"column_name"`
	ColumnDefault  sql.NullString `db:"column_default"`
	ColumnComment  sql.NullString `db:"column_comment"`
	ColumnDataType sql.NullString `db:"column_data_type"`
	CharLength     sql.NullInt32  `db:"char_len"`
	NumericLength  sql.NullString `db:"numeric_length"`
	ColumnNullable sql.NullBool   `db:"nullable"`
}

type rawIndexInfo struct {
	SchemaName  string `db:"schema_name"`
	TableName   string `db:"table_name"`
	IndexName   string `db:"index_name"`
	IsUnique    bool   `db:"is_unique"`
	IsPrimary   bool   `db:"pkey"`
	ColumnNames string `db:"column_names"`
}

// For storing the result we get from the DB about column data
var columns []rawCol

func NewCodeGenerator(config CodegenConfig) (Generator, appError.Typ) {
	g := Generator{
		Config:       config,
		pluralClient: pluralize.NewClient(),
		Schemas:      map[string]DbSchema{},
		Enums:        map[string]EnumDefinition{},
	}

	if strings.TrimSpace(config.MagicComment) == "" {
		g.Config.MagicComment = DefaultMagicComment
	}

	return g, appError.BlankError
}

func (g *Generator) Generate() appError.Typ {
	// Attempt connecting
	db, err := sqlx.Connect("pgx", g.Config.PgDbUrl)
	if err != nil {
		return appError.NewError(
			appError.Panic,
			"1NPL60",
			fmt.Sprintf("Could not connect. Error: %v", err))
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("E#1NPGEQ - Error when deferring db Close:", err)
		}
	}(db)
	fmt.Println("I#1NPKWR - Looks like we got connected")
	// return appError.BlankError
	// We will first run the query which will fetch the details
	var customCodeInFile []string
	var fileAlreadyExists bool
	var importList []string
	var importsString string
	var outputFile *os.File
	var fileContent string

	// Now let's try to generate the enumerations. Because we might need them later when generating table code
	// Validate the enumerations first
	for key, enum := range g.Config.Enumerations {
		if enum.DisableGeneration {
			continue
		}
		// enumCopy := enum
		// Make sure that key name is same as the enum name
		if key != enum.Name {
			panic(fmt.Sprintf("E#1PM3T2 - Key name %v of the enum does not match its name %v", key, enum.Name))
		}
		// Make sure that integer values are all unique and non-negative
		valSet := set.New[int16]()
		newMapping := map[string]int16{}
		for txtVal, intVal := range enum.Mappings {
			if intVal < 0 {
				panic(fmt.Sprintf("P#1PL54Y - Negative value (%v) for %v in enum %v. Enumerations cannot have negative values.", intVal, txtVal, enum.Name))
			}
			oldTxtVal, found := valSet.Value(intVal)
			if found {
				// This value is already there!
				panic(fmt.Sprintf("P#1PL4EK - Integer value %v has already been set with %v. It is also being set against %v. Please fix.", intVal, oldTxtVal, txtVal))
			}

			if strings.TrimSpace(strings.ToLower(txtVal)) == "undefined" {
				panic(fmt.Sprintf("E#1POGSI - Enum named %v contains a value named %v while also enabling the option to handle undefined values. This is not allowed.", enum.Name, txtVal))
			}

			valSet.Add(intVal, txtVal)

			newMapping[getGoName(txtVal)] = intVal
		}
		enum.Mappings = newMapping

		// Get the names created
		enum.goName = getGoName(enum.Name)
		enum.goNameSingular = g.pluralClient.Singular(enum.goName)
		enum.goNamePlural = g.pluralClient.Plural(enum.goName)
		enum.goTypeName = enum.goNameSingular
		if !enum.Exported {
			enum.goTypeName = lowerFirstChar(enum.goNameSingular)
		}

		// Now set this to the generator
		g.Enums[enum.Name] = enum
	}

	// Validated. Now generate
	enumImportsStr := ""
	enumImportsList := []string{}
	enumFileContentsStr := ""
	enumTemplate := `
//{{PACKAGE_NAME}}

//{{IMPORT_LIST}}

//{{FILE_CONTENTS}}

//{{MAGIC_COMMENT}}
`

	for _, enum := range g.Enums {
		fileAlreadyExists = true
		enumImportsList = []string{}
		enumImportsStr = ""
		enumFileContentsStr, enumImportsList = g.buildEnumContentString(enum, enumImportsList)

		if len(enumImportsList) > 0 {
			enumImportsStr += "\nimport (\n"
			for _, impo := range enumImportsList {
				enumImportsStr += "\t\"" + impo + "\"\n"
			}
			enumImportsStr += ")\n"
		} else {
			enumImportsStr = ""
		}

		fileContent = enumTemplate
		fileContent = strings.ReplaceAll(fileContent, "//{{PACKAGE_NAME}}", fmt.Sprintf("package %v", g.Config.DbModelPackageName))
		fileContent = strings.ReplaceAll(fileContent, "//{{IMPORT_LIST}}", enumImportsStr)
		fileContent = strings.ReplaceAll(fileContent, "//{{FILE_CONTENTS}}", enumFileContentsStr)
		fileContent = strings.ReplaceAll(fileContent, "//{{MAGIC_COMMENT}}", g.Config.MagicComment)
		fmt.Println("E#1PO4OB - Printing to make sure the variable gets used: ", enum)
		outputFileName := "gen_enum_" + strings.ToLower(enum.Name) + ".go"
		// Check if the file already exists
		existingFileContentBytes, fileErr := os.ReadFile(
			fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
		if fileErr != nil {
			// File does not exist
			fileAlreadyExists = false
		}

		if !fileAlreadyExists {
			// File has to be created.
			fileContent = strings.ReplaceAll(fileContent, "//{{FIRST_TIME_FILE_CONTENT}}",
				"// Make sure code below is valid before running code generator else the generator will fail\n\n")
		} else {
			// file already exists
			fileContent = strings.ReplaceAll(fileContent, "//{{FIRST_TIME_FILE_CONTENT}}", "")
		}

		existingFileContent := string(existingFileContentBytes)

		// Look for the magic comment
		if strings.Contains(existingFileContent, g.Config.MagicComment) {
			allcode := strings.Split(existingFileContent, g.Config.MagicComment)
			for i := 0; i < len(allcode); i++ {
				if i > 0 {
					customCodeInFile = append(customCodeInFile, allcode[i])
				}
			}
		}

		err = os.Mkdir(g.Config.DbModelPackagePath, 0777)
		if err != nil {
			// fmt.Println("E#1OBP5N -", err)
		}

		outputFile, err = os.Create(fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
		if err != nil {
			panic(fmt.Sprintf("P#1OECMC - %v", err))
		}

		fileContentBytes, err := format.Source([]byte(fileContent))
		if err != nil {
			panic(err)
		}

		fileContent = string(fileContentBytes)

		if fileAlreadyExists {
			for _, val := range customCodeInFile {
				fileContent += val + "\n"
			}
		}

		fileContent = g.removeTrailingNewlines(fileContent) + "\n"

		_, err = outputFile.WriteString(fileContent)
		if err != nil {
			return appError.NewError(appError.Error, "1OBPB4", err.Error())
		}

		err = outputFile.Close()
		if err != nil {
			return appError.NewError(appError.Error, "1PYY78", err.Error())
		}
	}

	// ==================================================
	// Enum generation done. Proceed for table generation
	// ==================================================

	err = db.Select(&columns, tableInfoQuery)
	if err != nil {
		panic(err)
	}
	var tables = map[string]DbTable{}
	// We need to iterate over the list of columns and create DbTable instances,
	for _, columnDetail := range columns {
		// If the schema is not yet built, build it.
		// If the table in that schema is not yet built then build it.
		// Column won't have been built for sure. So build that anyway.
		if !columnDetail.Schema.Valid || (columnDetail.Schema.Valid && columnDetail.Schema.String == "") {
			return appError.NewError(appError.Error, "1NXLYE", fmt.Sprintf("Not possible for column %v in table %v to not have a schema", columnDetail.ColumnName.String))
		}

		if table, tableOk := tables[columnDetail.Schema.String+"."+columnDetail.TableName.String]; tableOk {
			dbColProp, colComment, appErr := g.getCommentAndPropertyFromComment(columnDetail.ColumnComment.String)
			if appErr.IsNotBlank() {
				if appErr.Code == noSuchEnumErrCode {
					panic(fmt.Sprintf("P#1RD974 -  Error with enum referenced in the column properties of column %v.%v.%v: %v", table.Schema, table.Name, columnDetail.ColumnName, appErr))
				}
			}
			goDataType, networkDataType := g.getGoType(columnDetail.ColumnDataType.String, columnDetail.ColumnNullable.Bool)
			dbCol := DbColumn{
				Schema:            columnDetail.Schema.String,
				Table:             columnDetail.TableName.String,
				Name:              columnDetail.ColumnName.String,
				GoName:            getGoName(columnDetail.ColumnName.String),
				GoNameSingular:    g.pluralClient.Singular(getGoName(columnDetail.ColumnName.String)),
				GoNamePlural:      g.pluralClient.Plural(getGoName(columnDetail.ColumnName.String)),
				DataType:          columnDetail.ColumnDataType.String,
				GoDataType:        goDataType,
				NetworkDataType:   networkDataType,
				Comment:           colComment,
				HasDefaultValue:   columnDetail.ColumnDefault.Valid,
				DefaultValue:      columnDetail.ColumnDefault.String,
				CharacterLength:   int(columnDetail.CharLength.Int32),
				Nullable:          columnDetail.ColumnNullable.Bool,
				CommentProperties: dbColProp,
			}
			table.ColumnMap[columnDetail.ColumnName.String] = dbCol
			// collist := table.ColumnList
			// collist = append(collist, dbCol.Name)
			// table.ColumnList = collist
			table.ColumnList = append(table.ColumnList, dbCol.Name)
			table.ColumnListA2z = append(table.ColumnListA2z, dbCol.Name)
			tables[columnDetail.Schema.String+"."+columnDetail.TableName.String] = table
		} else {
			dbColProp, colComment, appErr := g.getCommentAndPropertyFromComment(columnDetail.ColumnComment.String)
			if appErr.IsNotBlank() {
				if appErr.Code == noSuchEnumErrCode {
					panic(fmt.Sprintf("P#1RD96Y - Error with enum referenced in the column properties of column %v.%v.%v: %v", table.Schema, table.Name, columnDetail.ColumnName, appErr))
				}
			}

			goDataType, networkDataType := g.getGoType(columnDetail.ColumnDataType.String, columnDetail.ColumnNullable.Bool)
			dbCol := DbColumn{
				Schema:            columnDetail.Schema.String,
				Table:             columnDetail.TableName.String,
				Name:              columnDetail.ColumnName.String,
				GoName:            getGoName(columnDetail.ColumnName.String),
				GoNameSingular:    g.pluralClient.Singular(getGoName(columnDetail.ColumnName.String)),
				GoNamePlural:      g.pluralClient.Plural(getGoName(columnDetail.ColumnName.String)),
				DataType:          columnDetail.ColumnDataType.String,
				GoDataType:        goDataType,
				NetworkDataType:   networkDataType,
				Comment:           colComment,
				HasDefaultValue:   columnDetail.ColumnDefault.Valid,
				DefaultValue:      columnDetail.ColumnDefault.String,
				CharacterLength:   int(columnDetail.CharLength.Int32),
				Nullable:          columnDetail.ColumnNullable.Bool,
				CommentProperties: dbColProp,
			}
			table = DbTable{
				Name:           columnDetail.TableName.String,
				GoName:         getGoName(columnDetail.TableName.String),
				GoNameSingular: g.pluralClient.Singular(getGoName(columnDetail.TableName.String)),
				GoNamePlural:   g.pluralClient.Plural(getGoName(columnDetail.TableName.String)),
				Schema:         columnDetail.Schema.String,
				Comment:        columnDetail.TableComment.String,
				ColumnMap:      map[string]DbColumn{columnDetail.ColumnName.String: dbCol},
				FKeyMap:        map[string]DbFkInfo{},
				ColumnList:     []string{},
				ColumnListA2z:  []string{},
			}
			// collist := table.ColumnList
			// collist = append(collist, dbCol.Name)
			// table.ColumnList = collist
			table.ColumnList = append(table.ColumnList, dbCol.Name)
			table.ColumnListA2z = append(table.ColumnListA2z, dbCol.Name)
			tables[columnDetail.Schema.String+"."+columnDetail.TableName.String] = table
		}
	}

	// Tables collected. Now sort into schemas
	for _, table := range tables {
		if _, schemaOk := g.Schemas[table.Schema]; schemaOk {
			g.Schemas[table.Schema].Tables[table.Name] = table
		} else {
			s := DbSchema{
				Name:   table.Schema,
				GoName: getGoName(table.Schema),
				Tables: map[string]DbTable{table.Name: table},
			}
			g.Schemas[table.Schema] = s
		}
	}

	// Now we need to create a new struct for each schema.
	var pkColumnNames []string
	for _, schema := range g.Schemas {
		// MARKER : Find the primary keys for each table
		for tableName, table := range schema.Tables {
			pkColumnNames = []string{}
			queryFormat := primaryKeyInfoQuery
			query := fmt.Sprintf(queryFormat, schema.Name+"."+table.Name)

			err = db.Select(&pkColumnNames, query)
			if err != nil {
				fmt.Println("..............", err)
			}
			// fmt.Println("I#1O4FCO - ", schema.Name, ".", table.Name, "==>", strings.Join(pkColumnNames, ","))

			for _, colname := range pkColumnNames {
				table.PkColumnList = append(table.PkColumnList, table.ColumnMap[colname])
			}
			g.Schemas[schema.Name].Tables[tableName] = table
		}
	}

	var indexes []rawIndexInfo
	for _, schema := range g.Schemas {
		for tableName, table := range schema.Tables {
			indexes = []rawIndexInfo{}
			queryFormat := tableIndexQuery
			query := fmt.Sprintf(queryFormat, schema.Name, table.Name)

			err = db.Select(&indexes, query)
			if err != nil {
				fmt.Println("E#1O4AM4 - Error in getting indexes: ", err)
			}

			for _, index := range indexes {
				colList := []DbColumn{}
				cols := strings.Split(index.ColumnNames, ",")
				for _, col := range cols {
					colObj, colExists := table.ColumnMap[col]
					if !colExists {
						panic("P#1O4CSW - Expected the column to be there.")
					}
					// colObj, colFindErr := getColumnFromListByName(col, table.ColumnList)
					// if colFindErr != nil {
					//
					// }
					if colObj.Name != "" {
						colList = append(colList, colObj)
					}
				}
				// build index struct
				i := DbIndex{
					Name:       index.IndexName,
					IsUnique:   index.IsUnique,
					IsPrimary:  index.IsPrimary,
					ColumnList: colList,
				}
				table.IndexList = append(table.IndexList, i)
			}
			g.Schemas[table.Schema].Tables[tableName] = table
		}
	}

	// Foreign keys
	var fkInfoArr []fkInfoFromDb
	queryFormat := tableForeignKeyQuery

	query := fmt.Sprintf(queryFormat)

	err = db.Select(&fkInfoArr, query)
	if err != nil {
		fmt.Printf("E#1OMHRT - %v\n", err)
	}

	// Forward references
	for _, fkInf := range fkInfoArr {
		_, schemaFound := g.Schemas[fkInf.FromSchema]
		if !schemaFound {
			panic(fmt.Sprintf("P#1OAJ30 - Expected to find schema %v but was not found", fkInf.FromSchema))
		}

		_, tableFound := g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable]
		if !tableFound {
			panic(fmt.Sprintf("P#1OAKXZ - Expected to find table %v in schema %v but was not found", fkInf.FromTable, fkInf.FromSchema))
		}

		_, fkInfoFound := g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable].FKeyMap[fkInf.ConstraintName]
		if !fkInfoFound {
			// g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable].FKeyMap = map[string]DbFkInfo{}

			g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable].FKeyMap[fkInf.ConstraintName] = DbFkInfo{
				FromSchema:   fkInf.FromSchema,
				FromTable:    fkInf.FromTable,
				ToSchema:     fkInf.ToSchema,
				ToTable:      fkInf.ToTable,
				FromColOrder: []string{},
				References: map[string]string{
					fkInf.FromColumn: fkInf.ToColumn,
				},
				ConstraintName: fkInf.ConstraintName,
			}
		} else {
			g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable].FKeyMap[fkInf.ConstraintName].References[fkInf.FromColumn] = fkInf.ToColumn
		}

		fkeyFromSchemaChain := g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable].FKeyMap[fkInf.ConstraintName]
		fkeyFromSchemaChain.FromColOrder = append(fkeyFromSchemaChain.FromColOrder, fkInf.FromColumn)
		g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable].FKeyMap[fkInf.ConstraintName] = fkeyFromSchemaChain
	}

	// Check for duplicates
	// Important: I was here
	dupLink := set.New[string]()
	for _, schema := range g.Schemas {
		for _, table := range schema.Tables {
			dupLink.Empty()
			for _, fkey := range table.FKeyMap {
				dupLink.Add(fmt.Sprintf("%v.%v->%v.%v", fkey.FromSchema, fkey.FromTable, fkey.ToSchema, fkey.ToTable))
			}
		}
	}

	// Reverse References
	for _, schema := range g.Schemas {
		for _, table := range schema.Tables {
			for _, fkey := range table.FKeyMap {
				toSchema, schemaFound := g.Schemas[fkey.ToSchema]
				if !schemaFound {
					panic(fmt.Sprintf("P#1OU0YK - %v schema expected but not found", fkey.ToSchema))
				}

				toTable, tableFound := toSchema.Tables[fkey.ToTable]
				if !tableFound {
					panic(fmt.Sprintf("P#1OU10C - table %v not found in schema %v", fkey.ToTable, fkey.ToSchema))
				}

				revFkeyMap := toTable.RevFKeyMap
				if revFkeyMap == nil {
					revFkeyMap = map[string]DbRevFkInfo{}
				}

				fromColOrderForRevRef := []string{}
				reverseReferences := map[string]string{}
				for _, fromColName := range fkey.FromColOrder {
					toColName, toColNameFound := fkey.References[fromColName]
					if !toColNameFound {
						panic(fmt.Sprintf("P#1OU44G - FromCol was not found %v", fromColName))
					}
					fromColOrderForRevRef = append(fromColOrderForRevRef, fromColName)
					reverseReferences[fromColName] = toColName
				}

				fromTable, fromTableFound := schema.Tables[fkey.FromTable]
				if !fromTableFound {
					panic(fmt.Sprintf("P#1OUPO8 - fromTable %v not found in schema %v", fkey.FromTable, fkey.FromSchema))
				}

				uniqIdx := false
				revIdx := fromTable.FindIndexByColumnNames(fromColOrderForRevRef)
				if revIdx.IsUnique || revIdx.IsPrimary {
					uniqIdx = true
				}

				revFkeyMap[fkey.GetReverseRefName()] = DbRevFkInfo{
					DbFkInfo: DbFkInfo{
						FromSchema:     fkey.ToSchema,
						FromTable:      fkey.ToTable,
						ToSchema:       fkey.FromSchema,
						ToTable:        fkey.FromTable,
						FromColOrder:   fromColOrderForRevRef,
						References:     reverseReferences,
						ConstraintName: fkey.GetReverseRefName(),
					},
					UniqueIndex: uniqIdx,
				}
				toTable.RevFKeyMap = revFkeyMap

				schema.Tables[fkey.ToTable] = toTable
			}
		}
		g.Schemas[schema.Name] = schema
	}

	// There is some trouble with the reverse references. Let's look at them a bit
	for _, schema := range g.Schemas {
		if schema.Name != "auth" {
			continue
		}

		for _, table := range schema.Tables {
			for revKeyName, r := range table.RevFKeyMap {
				for _, fromColName := range r.FromColOrder {
					fCol, tCol := r.References[fromColName]
					fmt.Printf("I#1P6NNX - From: %v.%v \tTo %v.%v \t %v \n", r.FromTable, fCol, r.ToTable, tCol, revKeyName)
				}
			}
		}
	}
	// Let's sort things alphabetically and ordinarily where possible
	// Sort the tables in schemas
	for _, schema := range g.Schemas {
		tablist := []string{}
		for _, table := range schema.Tables {
			// let's sort the columns on this table
			colList := []string{}
			for _, col := range table.ColumnMap {
				colList = append(colList, col.Name)
			}
			slices.Sort(colList)
			table.ColumnListA2z = colList
			schema.Tables[table.Name] = table

			tablist = append(tablist, table.Name)
		}
		slices.Sort(tablist)
		schema.TablesA2z = tablist

		g.Schemas[schema.Name] = schema
	}

	// MARKER: Start processing for Schema structs
	schemaFileTemplate := `
//{{PACKAGE_NAME}}

//{{IMPORT_LIST}}

//{{SCHEMA_STRUCT}}

//{{MAGIC_COMMENT}}
//{{FIRST_TIME_FILE_CONTENT}}
`

	var tableStructStr string
	var tableValidationStr string
	var tableBaseFuncsStr string

	for _, schema := range g.Schemas {
		fileAlreadyExists = true
		importList = []string{}
		importsString = ""
		tableStructStr, importList = g.buildSchemaStructString(schema.Name, importList)
		if len(importList) > 0 {
			importsString += "\nimport (\n"
			for _, impo := range importList {
				importsString += "\t\"" + impo + "\"\n"
			}
			importsString += ")\n"
		} else {
			importsString = ""
		}
		fileContent = schemaFileTemplate
		fileContent = strings.ReplaceAll(fileContent, "//{{PACKAGE_NAME}}", fmt.Sprintf("package %v", g.Config.DbModelPackageName))
		fileContent = strings.ReplaceAll(fileContent, "//{{IMPORT_LIST}}", importsString)
		fileContent = strings.ReplaceAll(fileContent, "//{{SCHEMA_STRUCT}}", tableStructStr)
		fileContent = strings.ReplaceAll(fileContent, "//{{MAGIC_COMMENT}}", g.Config.MagicComment)

		outputFileName := "gen_schema_" + strings.ToLower(schema.Name) + ".go"
		// Check if the file already exists
		existingFileContentBytes, fileErr := os.ReadFile(
			fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
		if fileErr != nil {
			// File does not exist
			fileAlreadyExists = false
		}

		if !fileAlreadyExists {
			// File has to be created.
			fileContent = strings.ReplaceAll(fileContent, "//{{FIRST_TIME_FILE_CONTENT}}",
				"// Make sure code below is valid before running code generator else the generator will fail\n\n")
		} else {
			// file already exists
			fileContent = strings.ReplaceAll(fileContent, "//{{FIRST_TIME_FILE_CONTENT}}", "")
		}

		existingFileContent := string(existingFileContentBytes)

		// Look for the magic comment
		if strings.Contains(existingFileContent, g.Config.MagicComment) {
			allcode := strings.Split(existingFileContent, g.Config.MagicComment)
			for i := 0; i < len(allcode); i++ {
				if i > 0 {
					customCodeInFile = append(customCodeInFile, allcode[i])
				}
			}
		}

		err = os.Mkdir(g.Config.DbModelPackagePath, 0777)
		if err != nil {
			// fmt.Println("E#1OBP5N -", err)
		}

		outputFile, err = os.Create(fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
		if err != nil {
			panic(fmt.Sprintf("P#1OECMC - %v", err))
		}

		fileContentBytes, err := format.Source([]byte(fileContent))
		if err != nil {
			panic(err)
		}

		fileContent = string(fileContentBytes)

		if fileAlreadyExists {
			for _, val := range customCodeInFile {
				fileContent += val + "\n"
			}
		}

		fileContent = g.removeTrailingNewlines(fileContent) + "\n"

		_, err = outputFile.WriteString(fileContent)
		if err != nil {
			return appError.NewError(appError.Error, "1OBPB4", err.Error())
		}

		err = outputFile.Close()
		if err != nil {
			return appError.NewError(appError.Error, "1PYY7G", err.Error())
		}

		customCodeInFile = []string{}
	}

	// Table struct file template
	tableStructFileTemplate := `
//{{PACKAGE_NAME}}

//{{IMPORT_LIST}}

//{{TABLE_STRUCT}}

//{{TABLE_BASE_FUNCS}}

//{{MAGIC_COMMENT}}
//{{FIRST_TIME_FILE_CONTENT}}
`
	for _, schema := range g.Schemas {
		for _, table := range schema.Tables {
			fileAlreadyExists = true
			importList = []string{}
			importsString = ""
			tableStructStr, importList = g.buildTableStructString(table, importList)
			tableValidationStr, importList = g.buildTableValidationFuncs(table, importList)
			tableBaseFuncsStr, importList = g.buildTableBaseFuncs(table, importList)
			if len(importList) > 0 {
				importsString += "\nimport (\n"
				for _, impo := range importList {
					importsString += "\t\"" + impo + "\"\n"
				}
				importsString += ")\n"
			} else {
				importsString = ""
			}

			fileContent = tableStructFileTemplate
			fileContent = strings.ReplaceAll(fileContent, "//{{PACKAGE_NAME}}", fmt.Sprintf("package %v", g.Config.DbModelPackageName))
			fileContent = strings.ReplaceAll(fileContent, "//{{IMPORT_LIST}}", importsString)
			fileContent = strings.ReplaceAll(fileContent, "//{{TABLE_STRUCT}}", tableStructStr)
			fileContent = strings.ReplaceAll(fileContent, "//{{TABLE_BASE_FUNCS}}", tableBaseFuncsStr)
			fileContent = strings.ReplaceAll(fileContent, "//{{MAGIC_COMMENT}}", g.Config.MagicComment)

			outputFileName := "gen_schema_" + strings.ToLower(schema.Name) + "_" + strings.ToLower(table.Name) + ".go"
			// Check if the file already exists
			existingFileContentBytes, fileErr := os.ReadFile(
				fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
			if fileErr != nil {
				// File does not exist
				fileAlreadyExists = false
			}

			if !fileAlreadyExists {
				// File has to be created.
				fileContent = strings.ReplaceAll(fileContent, "//{{FIRST_TIME_FILE_CONTENT}}",
					"// Make sure code below is valid before running code generator else the generator will fail\n\n"+tableValidationStr)
			} else {
				// file already exists
				fileContent = strings.ReplaceAll(fileContent, "//{{FIRST_TIME_FILE_CONTENT}}", "")
			}

			existingFileContent := string(existingFileContentBytes)

			// Look for the magic comment
			if strings.Contains(existingFileContent, g.Config.MagicComment) {
				allcode := strings.Split(existingFileContent, g.Config.MagicComment)
				for i := 0; i < len(allcode); i++ {
					if i > 0 {
						customCodeInFile = append(customCodeInFile, allcode[i])
					}
				}
			}

			err = os.Mkdir(g.Config.DbModelPackagePath, 0777)
			if err != nil {
				// fmt.Println("E#1OFXLN -", err)
			}

			outputFile, err = os.Create(fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
			if err != nil {
				panic(fmt.Sprintf("P#1OFXLQ - %v", err))
			}

			fileContentBytes, err := format.Source([]byte(fileContent))
			if err != nil {
				panic(fmt.Sprintf("P#1OFXM7 - %v", err))
			}

			fileContent = string(fileContentBytes)

			if fileAlreadyExists {
				for _, val := range customCodeInFile {
					fileContent += val + "\n"
				}
			}

			fileContent = g.removeTrailingNewlines(fileContent) + "\n"

			_, err = outputFile.WriteString(fileContent)
			if err != nil {
				return appError.NewError(appError.Error, "1OBPB4", err.Error())
			}

			err = outputFile.Close()
			if err != nil {
				return appError.NewError(appError.Error, "1OBPCE", err.Error())
			}

			customCodeInFile = []string{}
		}
	}

	networkFileTemplate := `
//{{PACKAGE_NAME}}

//{{IMPORT_LIST}}

//{{NETWORK_STRUCT}}

//{{MAGIC_COMMENT}}

`

	var networkStructStr string

	for _, schema := range g.Schemas {
		for _, table := range schema.Tables {
			fileAlreadyExists = true
			importList = []string{}
			importsString = ""
			networkStructStr, importList = g.buildNetworkStructString(table, importList)
			if len(importList) > 0 {
				importsString += "\nimport (\n"
				for _, impo := range importList {
					importsString += "\t\"" + impo + "\"\n"
				}
				importsString += ")\n"
			} else {
				importsString = ""
			}

			fileContent = networkFileTemplate
			fileContent = strings.ReplaceAll(fileContent, "//{{PACKAGE_NAME}}", fmt.Sprintf("package %v", g.Config.NetworkPackageName))
			fileContent = strings.ReplaceAll(fileContent, "//{{IMPORT_LIST}}", importsString)
			fileContent = strings.ReplaceAll(fileContent, "//{{NETWORK_STRUCT}}", networkStructStr)
			fileContent = strings.ReplaceAll(fileContent, "//{{MAGIC_COMMENT}}", g.Config.MagicComment)

			err = os.Mkdir(g.Config.NetworkPackagePath, 0777)
			if err != nil {
				fmt.Println("E#1RDF4N -", err)
			}

			outputFileName := "gen_schema_" + strings.ToLower(schema.Name) + "_" + strings.ToLower(table.Name) + "_net.go"
			// Check if the file already exists
			existingFileContentBytes, fileErr := os.ReadFile(
				fmt.Sprintf("%s/%s", g.Config.NetworkPackagePath, outputFileName))
			if fileErr != nil {
				// File does not exist
				fileAlreadyExists = false
			}

			existingFileContent := string(existingFileContentBytes)

			// Look for the magic comment
			if strings.Contains(existingFileContent, g.Config.MagicComment) {
				allcode := strings.Split(existingFileContent, g.Config.MagicComment)
				for i := 0; i < len(allcode); i++ {
					if i > 0 {
						customCodeInFile = append(customCodeInFile, allcode[i])
					}
				}
			}

			err = os.Mkdir(g.Config.DbModelPackagePath, 0777)
			if err != nil {
				// fmt.Println("E#1OFXLN -", err)
			}

			outputFile, err = os.Create(fmt.Sprintf("%s/%s", g.Config.NetworkPackagePath, outputFileName))
			if err != nil {
				panic(fmt.Sprintf("P#1OFXLQ - %v", err))
			}

			fileContentBytes, err := format.Source([]byte(fileContent))
			if err != nil {
				panic(fmt.Sprintf("P#1OFXM7 - %v", err))
			}

			fileContent = string(fileContentBytes)

			if fileAlreadyExists {
				for _, val := range customCodeInFile {
					fileContent += val + "\n"
				}
			}

			fileContent = g.removeTrailingNewlines(fileContent) + "\n"

			_, err = outputFile.WriteString(fileContent)
			if err != nil {
				return appError.NewError(appError.Error, "1OBPB4", err.Error())
			}

			err = outputFile.Close()
			if err != nil {
				return appError.NewError(appError.Error, "1OBPCE", err.Error())
			}

			customCodeInFile = []string{}
		}
	}

	return appError.BlankError
}

func (g *Generator) getColumnFromListByName(colName string, colList []DbColumn) (DbColumn, error) {
	for _, col := range colList {
		if col.Name == colName {
			return col, nil
		}
	}
	return DbColumn{}, fmt.Errorf("E#1O4CIS - No such column")
}

// Function to get the Go type for DB and network for a given PostgreSQL data type
func (g *Generator) getGoType(datatype string, nullable bool) (string, string) {
	switch datatype {
	case "bigint":
		if nullable {
			return "sql.NullInt64", "*int64"
		}
		return "int64", "int64"
	case "integer":
		if nullable {
			return "sql.NullInt32", "*int32"
		}
		return "int32", "int32"
	case "smallint":
		if nullable {
			return "sql.NullInt16", "*int16"
		}
		return "int16", "int16"
	case "numeric", "double precision":
		if nullable {
			return "sql.NullFloat64", "*float64"
		}
		return "float64", "float64"
	case "boolean":
		if nullable {
			return "sql.NullBool", "*bool"
		}
		return "bool", "bool"
	case "character varying", "text", "uuid":
		if nullable {
			return "sql.NullString", "*string"
		}
		return "string", "string"
	case "jsonb":
		return "jsonObject.Typ", "jsonObject.Typ"
	case "timestamp without time zone", "timestamp", "timestamp with time zone":
		if nullable {
			return "sql.NullTime", "*time.Time"
		}
		return "time.Time", "time.Time"
	default:
		return "any", "any"
	}
}

func (g *Generator) addToImports(str string, impList []string) []string {
	strExists := false
	for _, s := range impList {
		if s == str {
			strExists = true
		}
	}

	if !strExists && str != "" {
		impList = append(impList, str)
	}
	return impList
}

func (g *Generator) getEnumByName(name string) (EnumDefinition, error) {
	enum, found := g.Enums[name]
	if !found {
		return EnumDefinition{}, fmt.Errorf("E#1RDB11 - So such enum")
	}
	return enum, nil
}

// This function tries to read the comment and separate the comment and the column properties json and return the
// properties object and the comment separately.
func (g *Generator) getCommentAndPropertyFromComment(comment string) (dbColumnProperty, string, appError.Typ) {
	dbColProp := dbColumnProperty{}
	// Split by delimiter
	commentParts := strings.Split(comment, g.Config.ColCommentSeparator)
	if len(commentParts) != 2 {
		// More than 2 would mean that the comment was not written properly
		if len(commentParts) > 2 {
			return dbColProp, comment, appError.NewError(appError.Error, "1RD7U6", fmt.Sprintf(
				"Column comment separator %v found %v times while there should have been just one",
				g.Config.ColCommentSeparator, len(commentParts)))
		} else {
			return dbColProp, comment, appError.BlankError
		}
	}

	colComment := commentParts[0]
	err := json.Unmarshal([]byte(commentParts[1]), &dbColProp)
	if err != nil {
		return dbColumnProperty{}, colComment, appError.NewError(appError.Error, "1RD88J",
			fmt.Sprintf("JSON Marshalling failed. Error: %v", err))
	}

	if dbColProp.Disabled {
		return dbColProp, colComment, appError.BlankError
	}

	// Check that the StrConversionViaEnum value exists in case it has been set.
	if strings.TrimSpace(dbColProp.StrConversionViaEnum) != "" {
		// This enum must be present in the list of enums
		if _, enumExists := g.Enums[dbColProp.StrConversionViaEnum]; !enumExists {
			return dbColProp, colComment, appError.NewError(appError.Panic, noSuchEnumErrCode, fmt.Sprintf("No such enum supplied to generator: %v", dbColProp.StrConversionViaEnum))
		}
	}

	return dbColProp, colComment, appError.BlankError
}
