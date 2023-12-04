package dbCodegen

import (
	"database/sql"
	"fmt"
	"go/format"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/gertd/go-pluralize"

	"github.com/techrail/ground/typs/appError"
)

// DbSchema represents the schema in the database
type DbSchema struct {
	Name   string             // Name of the schema
	GoName string             // The name of table as a go variable that we would use
	Tables map[string]DbTable // Map of name of table to their DbTable struct values
}

// DbTable represents a table in the database
type DbTable struct {
	Name           string              // Table name as got from the query
	GoName         string              // The name of table as a go variable that we would use
	GoNameSingular string              // Singular form of GoName
	GoNamePlural   string              // Plural form of GoName
	Schema         string              // The name of the schema where this table resides
	Comment        string              // Comment on the table
	ColumnMap      map[string]DbColumn // List of columns as a map from the name of the column to the DbColumn type
	ColumnList     []DbColumn          // List of columns (ordinal)
	ColumnListA2z  []string            // List of column names (alphabetical)
	PkColumnList   []DbColumn          // List of columns that make the primary key (slice because order matters)
	IndexList      []DbIndex           // List of indexes on this table
	FKeyMap        map[string]DbFkInfo // List of foreign keys in table as map from constraint name to DbFkInfo type
}

func (table *DbTable) fullyQualifiedTableName() string {
	return table.Schema + "." + table.Name
}

func (table *DbTable) fullyQualifiedStructName() string {
	return getGoName(table.Schema) + "_" + table.GoNameSingular
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

// Properties that would be expressed as json in the column comment **after** the actual comment
// For example, the comment on a `email` column with `(^_^)` as the comment separator can be something like this:
//
//	Email address of the user (^_^) {"minStrLen":8}
//
// It should generate the validator which would ensure that the email field at least has 8 characters in it!
type dbColumnProperty struct {
	MinStrLen            int       `json:"minStrLen"`            // for string columns - minimum length
	MaxStrLen            int       `json:"maxStrLen"`            // for string columns - maximum length
	MinIntVal            int64     `json:"minIntVal"`            // for integer columns - minimum value
	MaxIntVal            int64     `json:"maxIntVal"`            // for integer columns - maximum value
	MinTimestampVal      time.Time `json:"minTimestampVal"`      // For timestamp without time zone columns
	MaxTimestampVal      time.Time `json:"maxTimestampVal"`      // For timestamp without time zone columns
	StrValidateAs        string    `json:"strValidateAs"`        // Validate String Data as what? (Email? URL? Name? Regex?)
	HideFromNetwork      bool      `json:"hideFromNetwork"`      // Should this field be hidden in network response
	StrConversionViaType string    `json:"strConversionViaType"` // To be used for enumerated fields that need to be represented as string in network responses
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

// DbFkInfo represents a single foreign key in a table
type DbFkInfo struct {
	FromSchema     string            // The schema name of the table from which the foreign key reference is being made
	FromTable      string            // The table name which is referencing the target table
	ToSchema       string            // The schema name of the table whose column is being referenced
	ToTable        string            // The table name of the table which is being referenced
	References     map[string]string // The reference map ([from_column]to_column format)
	ConstraintName string            // Name of the foreign key constraint
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

// CodegenConfig contains the values and rules using which the code is to be generated
type CodegenConfig struct {
	PgDbUrl                  string // DB URL string for PostgreSQL database to which we have to connect
	DbModelPackageName       string // Name of the package under which the package for db related code will be placed
	DbModelPackagePath       string // Full path of the directory where the generated code for db will be placed
	MagicComment             string // Magic comment which allows the generator to update only the generated portion of code
	ColCommentSeparator      string // The string after which we can place the Properties JSON in column comment
	InsertCreatedAtInCode    bool   // Should the code for inserting `created_at` timestamp column be generated by the code generator?
	InsertUpdatedAtInCode    bool   // Should the code for inserting `updated_at` timestamp column be generated by the code generator?
	UpdateUpdatedAtInCode    bool   // Should the code for updating `updated_at` timestamp column be generated by the code generator?
	BuildUpdateByUniqueIndex bool   // Should we generate the update functions for unique indexes?
}

// Generator is the structure we return to a client which needs a generator.
// It is supposed to contain all the information needed for performing the code generation
type Generator struct {
	// More fields to be decided
	Config       CodegenConfig       // The configuration for this generator
	Schemas      map[string]DbSchema // The schemas in the database (will in turn contain tables)
	pluralClient *pluralize.Client   // Pluralization client
	sync.Mutex                       // To prevent parallel runs
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
	}

	if strings.TrimSpace(config.MagicComment) == "" {
		g.Config.MagicComment = "// MAGIC COMMENT (DO NOT EDIT): Please write any custom code only below this line.\n"
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
	err = db.Select(&columns, tableInfoQuery)
	if err != nil {
		panic(err)
	}
	var tables map[string]DbTable = map[string]DbTable{}
	// We need to iterate over the list of columns and create DbTable instances,
	for _, columnDetail := range columns {
		// If the schema is not yet built, build it.
		// If the table in that schema is not yet built then build it.
		// Column won't have been built for sure. So build that anyway.
		if !columnDetail.Schema.Valid || (columnDetail.Schema.Valid && columnDetail.Schema.String == "") {
			return appError.NewError(appError.Error, "1NXLYE", fmt.Sprintf("Not possible for column %v in table %v to not have a schema", columnDetail.ColumnName.String))
		}

		if table, tableOk := tables[columnDetail.Schema.String+"."+columnDetail.TableName.String]; tableOk {
			dbColProp, colComment := getCommentAndPropertyFromComment(columnDetail.ColumnComment.String)
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
			table.ColumnList = append(table.ColumnList, dbCol)
			table.ColumnListA2z = append(table.ColumnListA2z, dbCol.Name)
			tables[columnDetail.Schema.String+"."+columnDetail.TableName.String] = table
		} else {
			dbColProp, colComment := getCommentAndPropertyFromComment(columnDetail.ColumnComment.String)
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
				ColumnMap: map[string]DbColumn{
					columnDetail.ColumnName.String: dbCol,
				},
				FKeyMap: map[string]DbFkInfo{},
			}
			table.ColumnList = append(table.ColumnList, dbCol)
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
			fmt.Println("I#1O4FCO - ", schema.Name, ".", table.Name, "==>", strings.Join(pkColumnNames, ","))

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
					//colObj, colFindErr := getColumnFromListByName(col, table.ColumnList)
					//if colFindErr != nil {
					//
					//}
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
		fmt.Println("ERROR FKS -->-->-->-->-->-->", err)
	}

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
				FromSchema: fkInf.FromSchema,
				FromTable:  fkInf.FromTable,
				ToSchema:   fkInf.ToSchema,
				ToTable:    fkInf.ToTable,
				References: map[string]string{
					fkInf.FromColumn: fkInf.ToColumn,
				},
				ConstraintName: fkInf.ConstraintName,
			}
		} else {
			g.Schemas[fkInf.FromSchema].Tables[fkInf.FromTable].FKeyMap[fkInf.ConstraintName].References[fkInf.FromColumn] = fkInf.ToColumn
		}
	}

	// MARKER: Start processing for Schema structs
	schemaFileTemplate := `
//{{PACKAGE_NAME}}

//{{IMPORT_LIST}}

//{{SCHEMA_STRUCT}}

//{{MAGIC_COMMENT}}
//{{FIRST_TIME_FILE_CONTENT}}
`
	var fileAlreadyExists bool
	var customCodeInFile []string

	var tableStructStr string
	var tableBaseFuncsStr string
	var importsString string
	var importList []string

	for _, schema := range g.Schemas {
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
		fileContent := schemaFileTemplate
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
			fmt.Println("E#1OBP5N -", err)
		}

		outputFile, err := os.Create(fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
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
			return appError.NewError(appError.Error, "1OBPCE", err.Error())
		}
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
			importList = []string{}
			importsString = ""
			tableStructStr, importList = g.buildTableStructString(table, importList)
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

			fileContent := tableStructFileTemplate
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
				fmt.Println("E#1OFXLN -", err)
			}

			outputFile, err := os.Create(fmt.Sprintf("%s/%s", g.Config.DbModelPackagePath, outputFileName))
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
	case "numeric":
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
		return "jsonObject.Typ", "JsonObject.Typ"
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

// This function tries to read the comment and separate the comment and the column properties json and return the
// properties object and the comment separately.
// TODO: Implement later
func getCommentAndPropertyFromComment(comment string) (dbColumnProperty, string) {
	return dbColumnProperty{}, comment
}
