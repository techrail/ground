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

	schemaStruct += fmt.Sprintf("// %vSchema struct corresponds to the %v schema of the DB\n",
		schema.GoName, schema.Name)
	schemaStruct += fmt.Sprintf("type %sSchema struct {\n", schema.GoName)
	for _, table := range schema.Tables {
		fmt.Printf("Schema: %v | Table: %v \n",
			schema.Name, table.Name)
		tableComment := ""
		if table.Comment != "" {
			tableComment = "// " + strings.ReplaceAll(table.Comment, "\n", "")
		} else {
			tableComment = "// _No comment on table_"
		}
		schemaStruct += fmt.Sprintf("\t%s %s %v\n",
			table.GoNameSingular, table.fullyQualifiedStructName(), tableComment)
		schemaStruct += fmt.Sprintf("\t%sDao *%sDao // Dao for %v\n",
			table.GoNameSingular, table.fullyQualifiedStructName(), table.Name)
	}
	schemaStruct += "}\n"
	schemaStruct += fmt.Sprintf("var %v %vSchema\n\n", schema.GoName, schema.GoName)
	schemaStruct += "func init() {\n"
	schemaStruct += fmt.Sprintf("%v = %vSchema{\n", schema.GoName, schema.GoName)
	for _, table := range schema.Tables {
		schemaStruct += fmt.Sprintf("\t%sDao: New%sDao(), // Dao for %v\n",
			table.GoNameSingular, table.fullyQualifiedStructName(), table.Name)
	}
	schemaStruct += "}\n"
	schemaStruct += "}\n"

	return schemaStruct, importList
}
