package dbcodegen

import "fmt"

func (g *Generator) buildTableStructString(table DbTable, importList []string) (string, []string) {
	var tableStruct string
	tableStruct = ""

	// PublicAuditLog struct wraps around the audit_log struct to allow custom functions to be written
	// This file is only created once and is never edited by the generator Audit Log of activities happening inside the system
	// Use this space to create custom functions for Audit Log of activities happening inside the system or the Dao

	tableStruct += fmt.Sprintf("// %v struct wraps around the %v struct to allow custom functions to be written\n",
		table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	tableStruct += fmt.Sprintf("// NOTE: This file is only created once and is never edited by the generator\n")
	tableStruct += fmt.Sprintf("// IMPORTANT: Do not remove the embedding of the base struct or give it a custom name\n")
	tableStruct += fmt.Sprintf("//     Doing so will cause the generated base file to misbehave\n")
	tableStruct += fmt.Sprintf("// Use this file to create custom functions for %v or %v\n",
		table.fullyQualifiedStructName(), table.fullyQualifiedDaoName())
	tableStruct += fmt.Sprintf("type %s struct {\n", table.fullyQualifiedStructName())
	tableStruct += fmt.Sprintf("%v\n", table.fullyQualifiedBaseStructName())
	tableStruct += fmt.Sprintf("}\n")
	tableStruct += fmt.Sprintf("// File ends here\n")
	return tableStruct, importList
}

func (g *Generator) buildTableBaseStructString(table DbTable, importList []string) (string, []string) {
	var tableStruct string
	tableStruct = ""

	// tableName := getGoName(table.Name)
	tableStruct += fmt.Sprintf("// %v struct corresponds to the %v table in %v schema of the DB\n",
		table.fullyQualifiedBaseStructName(), table.Name, table.Schema)
	tableStruct += fmt.Sprintf("// Table Comment: %v\n", table.commentForStruct())
	tableStruct += fmt.Sprintf("type %s struct {\n", table.fullyQualifiedBaseStructName())

	colNames := table.ColumnList
	if g.Config.ColumnOrderAlphabetic {
		colNames = table.ColumnListA2z
	}
	for _, columnName := range colNames {
		column, columnFound := table.ColumnMap[columnName]
		if !columnFound {
			panic(fmt.Sprintf("P#1OL11R - Column %v not found in table %v of schema %v", columnName, table.Name, table.Schema))
		}

		// fmt.Printf("Table: %v | Column: %v | PG DataType: %v | Nullable: %v | Go DataType: %v\n",
		//	table.Name, column.Name, column.DataType, column.Nullable, column.GoDataType)

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
