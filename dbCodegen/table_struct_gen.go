package dbCodegen

import "fmt"

func (g *Generator) buildTableStructString(table DbTable, importList []string) (string, []string) {
	var tableStruct string
	tableStruct = ""

	// tableName := getGoName(table.Name)
	tableStruct += fmt.Sprintf("// %v struct corresponds to the %v table in %v schema of the DB\n",
		table.fullyQualifiedStructName(), table.Name, table.Schema)
	tableStruct += fmt.Sprintf("// Table Comment: %v\n", table.commentForStruct())
	tableStruct += fmt.Sprintf("type %s struct {\n", table.fullyQualifiedStructName())

	colNames := table.ColumnList
	if g.Config.ColumnOrderAlphabetic {
		colNames = table.ColumnListA2z
	}
	for _, columnName := range colNames {
		column, columnFound := table.ColumnMap[columnName]
		if !columnFound {
			panic(fmt.Sprintf("P#1OL11R - Column %v not found in table %v of schema %v", columnName, table.Name, table.Schema))
		}

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
