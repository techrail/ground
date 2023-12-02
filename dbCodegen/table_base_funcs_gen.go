package dbCodegen

import (
	"fmt"
	"strings"
)

func (g *Generator) buildTableBaseFuncs(table DbTable, importList []string) (string, []string) {
	tableInsertionValidationFuncStr := ""
	tableInsertionValidationFuncStr, importList = g.buildTableInsertionValidation(table, importList)
	tableInsertionFuncStr := ""
	tableInsertionFuncStr, importList = g.buildTableInsertMethod(table, importList)

	tableBaseFuncStr := tableInsertionValidationFuncStr + tableInsertionFuncStr
	return tableBaseFuncStr, importList
}

func (g *Generator) buildTableInsertionValidation(table DbTable, importList []string) (string, []string) {
	tabInsertValidation := ""
	tabInsertValidation += fmt.Sprintf("func (%v *%v) validateForInsertion() error {\n",
		table.variableName(), table.fullyQualifiedStructName())
	tabInsertValidation += fmt.Sprintf("err := %v.commonValidation()\n", table.variableName())
	tabInsertValidation += "if err != nil{\nreturn err\n}\n"
	tabInsertValidation += "// More code to be written here for validation\n"
	tabInsertValidation += "return nil\n"
	tabInsertValidation += "}\n\n"

	return tabInsertValidation, importList
}

func (g *Generator) buildTableInsertMethod(table DbTable, importList []string) (string, []string) {
	insertCode := ""
	insertCode += fmt.Sprintf("func (%v *%v) insert() error {\n",
		table.variableName(), table.fullyQualifiedStructName())
	insertCode += "var err error\n"
	insertCode += fmt.Sprintf("err = %v.validateForInsertion()\n", table.variableName())
	insertCode += "if err != nil{\nreturn err\n}\n"

	insertCode += fmt.Sprintf("insertQuery := `INSERT INTO %v.%v (\n", table.Schema, table.Name)

	colNameSlice := []string{}
	argPositionSlice := []string{}
	pkeyColumnNameSlice := []string{}
	pkeyAmpInsertedColumnNameSlice := []string{}
	returningColsSlice := []DbColumn{}
	goColumnNameSlice := []string{}
	i := 0
	for _, column := range table.ColumnList {
		if !(column.Name == "created_at" && g.Config.InsertCreatedAtInCode == false && // Created at timestamps might not need to be created in code
			(column.GoDataType == "time.Time" || column.GoDataType == "sql.NullTime")) &&
			!(column.Name == "updated_at" && g.Config.InsertUpdatedAtInCode == false && // Updated at timestamps might not need to be created in code
				(column.GoDataType == "time.Time" || column.GoDataType == "sql.NullTime")) {
			if column.HasDefaultValue && isColumnInList(column.Name, table.PkColumnList) {
				// Column already has default value for a primary column. Do not include this column in list
				continue
			}
			i += 1
			colNameSlice = append(colNameSlice, `"`+column.Name+`"`)
			if column.DataType == "json" || column.DataType == "jsonb" {
				goColumnNameSlice = append(goColumnNameSlice, fmt.Sprintf("%v.%v.String()", table.variableName(), column.GoName))
			} else {
				goColumnNameSlice = append(goColumnNameSlice, fmt.Sprintf("%v.%v", table.variableName(), column.GoName))
			}
			argPositionSlice = append(argPositionSlice, fmt.Sprintf("$%v", i))
		}
	}

	for _, column := range table.PkColumnList {
		pkeyColumnNameSlice = append(pkeyColumnNameSlice, `"`+column.Name+`"`)
		pkeyAmpInsertedColumnNameSlice = append(pkeyAmpInsertedColumnNameSlice, "&inserted"+column.GoName)
	}

	for _, column := range table.ColumnList {
		if (column.Name == "created_at" && g.Config.InsertCreatedAtInCode == false &&
			(column.GoDataType == "time.Time" || column.GoDataType == "sql.NullTime")) ||
			(column.Name == "updated_at" && g.Config.InsertUpdatedAtInCode == false &&
				(column.GoDataType == "time.Time" || column.GoDataType == "sql.NullTime")) {
			returningColsSlice = append(returningColsSlice, column)
		}
	}

	insertCode += "\t\t\t"
	insertCode += strings.Join(colNameSlice, ", \n\t\t\t")
	insertCode += "\n\t\t) VALUES (\n\t\t\t"
	insertCode += strings.Join(argPositionSlice, ", \n\t\t\t")
	insertCode += "\n\t\t) "
	insertCode += "RETURNING "
	insertCode += strings.Join(pkeyColumnNameSlice, ", ")
	for _, col := range returningColsSlice {
		insertCode += `, "` + col.Name + `"`
	}
	insertCode += "`;\n\n"

	//
	insertCode += fmt.Sprintf("resultRow := %v.QueryRowx(insertQuery,\n", upperFirstChar(g.Config.DbModelPackageName))
	insertCode += strings.Join(goColumnNameSlice, ", \n\t\t\t")
	insertCode += ",\n)\n\n"
	insertCode += "if resultRow.Err() != nil {\n"
	insertCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not insert into database: %v\", resultRow.Err())\n"
	importList = append(importList, "fmt", "errors")
	insertCode += "logger.Println(errMsg)\n"
	importList = append(importList, "github.com/techrail/ground/logger")
	insertCode += "return errors.New(errMsg)\n"
	insertCode += "}\n\n"

	for _, pkcol := range table.PkColumnList {
		insertCode += fmt.Sprintf("var inserted%v %v\n", pkcol.GoName, pkcol.GoDataType)
	}
	for _, col := range returningColsSlice {
		insertCode += fmt.Sprintf("var inserted%v %v\n", col.GoName, col.GoDataType)
	}

	returningColList := strings.Join(pkeyAmpInsertedColumnNameSlice, ",")
	for i, col := range returningColsSlice {
		if i == 0 && len(pkeyAmpInsertedColumnNameSlice) == 0 {
			// In case this table has no primary key, then returningColList must not begin with a `,`
			returningColList += `&inserted` + col.GoName
		} else {
			returningColList += `, &inserted` + col.GoName
		}
	}

	insertCode += fmt.Sprintf("\nerr = resultRow.Scan(%v)\n", returningColList)

	insertCode += "if err != nil {\n"
	insertCode += "return fmt.Errorf(\"E#" + newUniqueLmid() + " - Scan failed. Error: %v\", err)\n"
	insertCode += "}\n\n"

	for _, col := range table.PkColumnList {
		insertCode += fmt.Sprintf("%v.%v = inserted%v\n", lowerFirstChar(table.GoNameSingular), col.GoName, col.GoName)
	}
	for _, col := range returningColsSlice {
		insertCode += fmt.Sprintf("%v.%v = inserted%v\n", lowerFirstChar(table.GoNameSingular), col.GoName, col.GoName)
	}

	insertCode += "return nil"
	insertCode += "}\n\n"

	//importList = g.addToImports(g.Config.DbModelPackageName+"/resources", importList)
	importList = g.addToImports("fmt", importList)
	importList = g.addToImports("errors", importList)
	return insertCode, importList
}
