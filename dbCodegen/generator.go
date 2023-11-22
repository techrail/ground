package dbCodegen

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/techrail/ground/typs/appError"
)

// DbSchema represents the schema in the database
type DbSchema struct {
	Name   string             // Name of the schema
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
	ColumnList     []DbColumn          // List of columns as in array
	PkColumnList   []DbColumn          // List of columns that make the primary key of this table
	IndexList      []DbIndex           // List of indexes on this table
	FkList         []DbFkInfo          // List of foreign keys in this table
}

// DbColumn is the column representation of a table in the database for the generator
type DbColumn struct {
	Name            string // Column name
	GoName          string // Name we want to use for Golang code that will be generated
	GoNameSingular  string // Singular form of the name
	GoNamePlural    string // Plural form of the name
	DataType        string // Data type we get from db
	GoDataType      string // Data type we want to use in go program
	Comment         string // Column comment
	CharacterLength int    // Length in case it is varchar
	Nullable        bool   // NOT NULL means it is false
	HasDefaultValue bool   // Does the column have a default value?
	DefaultValue    string // If column has default value then what is it
}

// DbIndex represents an index inside a table
type DbIndex struct {
	Name       string     // Name of the index
	IsUnique   bool       // Is this a "unique" index?
	IsPrimary  bool       // Does this index correspond to the primary key of the table
	ColumnList []DbColumn // List of columns of this index (ordered the same way as the columns are defined in the index)
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

// CodegenConfig contains the values and rules using which the code is to be generated
type CodegenConfig struct {
	PgDbUrl             string // DB URL string for PostgreSQL database to which we have to connect
	TablePackageName    string // Name of the package under which the package for table related code will be placed
	TablePackagePath    string // Full path of the directory where the generated code for table will be placed
	MagicComment        string // Magic comment which allows the generator to update only the generated portion of code
	ColCommentSeparator string // The string after which we can place the Properties JSON in column comment
}

// Generator is the structure we return to a client which needs a generator.
// It is supposed to contain all the information needed for performing the code generation
type Generator struct {
	// More fields to be decided
	Config     CodegenConfig       // The configuration for this generator
	Schemas    map[string]DbSchema // The schemas in the database (will in turn contain tables)
	sync.Mutex                     // To prevent parallel runs
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

// For storing the result we get from the DB about column data
var columns []rawCol

func NewCodeGenerator(config CodegenConfig) (Generator, appError.Typ) {
	return Generator{Config: config}, appError.BlankError
}

func (g *Generator) Connect() appError.Typ {
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

	// We need to iterate over the list of columns and create DbTable instances,
	for _, columnDetail := range columns {
		// If the schema is not yet built, build it.
		// If the table in that schema is not yet built then build it.
		// Column won't have been built for sure. So build that anyway.
	}

	return appError.BlankError
}
