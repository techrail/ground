package dbCodegen

import (
	`database/sql`
	`fmt`

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	`github.com/techrail/ground/typs/appError`
)

type CodegenConfig struct {
	TablePackageName string // Name of the package under which the package for table related code will be placed
	TablePackagePath string // Full path of the directory where the generated code for table will be placed
	PgDbUrl          string // DB URL string for PostgreSQL database to which we have to connect
}

type Generator struct {
	// Fields to be decided
	Config CodegenConfig
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

// DbColumn is the column representation for the generator
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

	for i, c := range columns {
		fmt.Printf("Row %v, Col: %v", i, c)
	}

	return appError.BlankError
}
