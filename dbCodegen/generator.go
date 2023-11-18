package dbCodegen

import (
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
	return appError.BlankError
}
