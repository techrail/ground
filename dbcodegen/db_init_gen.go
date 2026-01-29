package dbcodegen

import (
	"fmt"
	"strings"
)

func (g *Generator) buildInitCode(importList []string, tables map[string]DbTable) (string, []string) {
	initCode := ""

	// Let us make the structure which will hold the enumerations
	initCode += "// Here are the enumerations for this DB\n"
	initCode += "type enums struct {"
	for _, enum := range g.Enums {
		// enumTypeName = enum.goNameSingular

		// IMPORTANT: In case you are ready to generate the abstractions only for the exported ones, enale the code below
		// if !enum.Exported {
		//     continue
		// }

		initCode += fmt.Sprintf("\n%v %v", enum.goNameSingular, "enum"+enum.goNameSingular+".Typ")
		importList = g.addToImports(g.Config.ModelsContainerPackage+"/"+g.Config.DbModelPackageName+"/enum"+enum.goNameSingular, importList)
	}
	initCode += "}\n"
	initCode += "var Enums enums \n"

	initCode += "\n"
	initCode += fmt.Sprintf("var %v db\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += fmt.Sprintf("var %vReader db\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += "\n"

	initCode += "\n"
	initCode += "// This piece of code initializes the DB connectors\n"
	initCode += "func init() {\n"
	initCode += "var err error\n\n"

	initCode += fmt.Sprintf("dbUrl := os.Getenv(\"DATABASE_%v_URL\")\n", strings.ToUpper(g.DbName))
	initCode += "if dbUrl == \"\" {\n"
	initCode += fmt.Sprintf("dbUrl = \"%v\"\n", g.Config.PgDbUrl)
	initCode += "}\n\n"

	initCode += fmt.Sprintf("%v.DB, err = sqlx.Connect(\"pgx\", dbUrl)\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += "if err != nil {\n"
	initCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not connect to the database! Error: %v\", err)\n"
	initCode += "fmt.Println(errMsg)"

	initCode += "}else{\n"
	initCode += fmt.Sprintf("%v.DB.SetMaxOpenConns(20) // Maximum number of open connections\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += fmt.Sprintf("%v.DB.SetMaxIdleConns(10) // Maximum number of idle connections\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += fmt.Sprintf("%v.DB.SetConnMaxIdleTime(5 * time.Minute) // How long an idle connection lives\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += fmt.Sprintf("%v.DB.SetConnMaxLifetime(30 * time.Minute) // Max lifetime of a connection\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += "}\n"
	importList = g.addToImports("time", importList)

	if g.readerEnabled() {
		initCode += fmt.Sprintf("%vReader.DB, err = sqlx.Connect(\"pgx\", \"%v\"\n", upperFirstChar(g.Config.DbModelPackageName), g.Config.PgReaderDbUrl)
		initCode += "if err != nil {\n"
		initCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not connect to the database! Error: %v\", err)\n"
		initCode += "fmt.Println(errMsg)"
		initCode += "}else{\n"
		initCode += fmt.Sprintf("%vReader.DB.SetMaxOpenConns(20) // Maximum number of open connections\n", upperFirstChar(g.Config.DbModelPackageName))
		initCode += fmt.Sprintf("%vReader.DB.SetMaxIdleConns(10) // Maximum number of idle connections\n", upperFirstChar(g.Config.DbModelPackageName))
		initCode += fmt.Sprintf("%vReader.DB.SetConnMaxIdleTime(5 * time.Minute) // How long an idle connection lives\n", upperFirstChar(g.Config.DbModelPackageName))
		initCode += fmt.Sprintf("%vReader.DB.SetConnMaxLifetime(30 * time.Minute) // Max lifetime of a connection\n", upperFirstChar(g.Config.DbModelPackageName))
		initCode += "}\n\n"
	} else {
		initCode += fmt.Sprintf("%vReader = db{\n", upperFirstChar(g.Config.DbModelPackageName))
		initCode += fmt.Sprintf("DB: %v.DB,\n", upperFirstChar(g.Config.DbModelPackageName))
		initCode += "}\n\n"
	}

	initCode += "}\n"

	return initCode, importList
}
