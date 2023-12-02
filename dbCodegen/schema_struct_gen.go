package dbCodegen

import (
	"fmt"
	"strings"
)

func (g *Generator) buildSchemaStructString(schemaName string, importList []string) (string, []string) {
	schema, schemaFound := g.Schemas[schemaName]
	if !schemaFound {
		panic(fmt.Sprintf("P#1OBETJ - Schema %v not found", schemaName))
	}

	schemaStruct := ""

	schemaStruct += fmt.Sprintf("// %v struct corresponds to the %v schema of the DB\n",
		schema.GoName, schema.Name)
	schemaStruct += fmt.Sprintf("type %s struct {\n", schema.GoName)
	for _, table := range schema.Tables {
		fmt.Printf("Schema: %v | Table: %v \n",
			schema.Name, table.Name)
		tableComment := ""
		if table.Comment != "" {
			tableComment = "// " + strings.ReplaceAll(table.Comment, "\n", "")
		}
		schemaStruct += fmt.Sprintf("\t%s %s %v\n",
			table.GoNameSingular, table.fullyQualifiedStructName(), tableComment)
	}
	schemaStruct += "}\n"

	return schemaStruct, importList
}
