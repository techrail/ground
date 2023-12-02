package dbCodegen

import "fmt"

func (g *Generator) buildTableStructString(table DbTable, importList []string) (string, []string) {
	var tableStruct string
	tableStruct = ""

	// tableName := getGoName(table.Name)
	tableStruct += fmt.Sprintf("// %v struct corresponds to the %v table in %v schema of the DB\n",
		table.GoNameSingular, table.Name, table.Schema)
	tableStruct += fmt.Sprintf("// Table Comment: %v\n", table.commentForStruct())
	tableStruct += fmt.Sprintf("type %s struct {\n", table.fullyQualifiedStructName())
	for _, column := range table.ColumnMap {
		fmt.Printf("Table: %v | Column: %v | PG DataType: %v | Nullable: %v | Go DataType: %v\n",
			table.Name, column.Name, column.DataType, column.Nullable, column.GoDataType)
		columnComment := ""
		if column.Comment != "" {
			columnComment = "// " + column.newlineEscapedComment()
		}
		tableStruct += fmt.Sprintf("\t%s %s `db:\"%s\"` %v\n",
			column.GoName, column.GoDataType, column.Name, columnComment)
		importList = g.addToImports(g.getGoImportForDataType(column.DataType, column.Nullable), importList)
	}
	tableStruct += "}\n"

	return tableStruct, importList
}
