package dbCodegen

import "fmt"

func (g *Generator) buildInitCode(importList []string) (string, []string) {
	initCode := ""
	initCode += "\n"
	initCode += fmt.Sprintf("var %v db\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += fmt.Sprintf("var %vReader db\n", upperFirstChar(g.Config.DbModelPackageName))
	initCode += "\n"
	initCode += fmt.Sprintf("// This piece of code initializes the DB connectors")
	initCode += "func init() {\n"
	initCode += "var error err\n"
	initCode += fmt.Sprintf("%v.DB, err = sqlx.Connect(\"pgx\", \"%v\"\n", upperFirstChar(g.Config.DbModelPackageName), g.Config.PgDbUrl)
	initCode += "if err != nil {\n"
	initCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not connect to the database! Error: %v\", err)\n"
	initCode += "fmt.Println(errMsg)"
	initCode += "}\n"

	if g.readerEnabled() {
		initCode += fmt.Sprintf("%vReader.DB, err = sqlx.Connect(\"pgx\", \"%v\"\n", upperFirstChar(g.Config.DbModelPackageName), g.Config.PgReaderDbUrl)
		initCode += "if err != nil {\n"
		initCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not connect to the database! Error: %v\", err)\n"
		initCode += "fmt.Println(errMsg)"
		initCode += "}\n"
	} else {
		initCode += fmt.Sprintf("%vReader = db{\n", upperFirstChar(g.Config.DbModelPackageName))
		initCode += fmt.Sprintf("DB: %v.DB", upperFirstChar(g.Config.DbModelPackageName))
		initCode += "}\n"
	}

	initCode += "}\n"

	return initCode, importList
}
