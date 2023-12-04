package dbCodegen

import (
	"fmt"
	"math"
	"strings"
)

func (g *Generator) buildTableBaseFuncs(table DbTable, importList []string) (string, []string) {
	tableBaseValidationFuncStr := ""
	tableBaseValidationFuncStr, importList = g.buildTableBaseValidation(table, importList)
	tableCommonValidationFuncStr := ""
	tableCommonValidationFuncStr, importList = g.buildTableCommonValidation(table, importList)
	tableInsertionValidationFuncStr := ""
	tableInsertionValidationFuncStr, importList = g.buildTableInsertValidation(table, importList)
	tableUpdateValidationFuncStr := ""
	tableUpdateValidationFuncStr, importList = g.buildTableUpdateValidation(table, importList)
	tableInsertionFuncStr := ""
	tableInsertionFuncStr, importList = g.buildTableInsertMethod(table, importList)
	tableUpdateFuncStr := ""
	tableUpdateFuncStr, importList = g.buildTableUpdateMethod(table, importList)

	tableBaseFuncStr := tableBaseValidationFuncStr + tableCommonValidationFuncStr + tableInsertionValidationFuncStr +
		tableUpdateValidationFuncStr + tableInsertionFuncStr + tableUpdateFuncStr
	return tableBaseFuncStr, importList
}

func (g *Generator) buildTableBaseValidation(table DbTable, importList []string) (string, []string) {
	tabCommonValidation := ""
	tabCommonValidation += fmt.Sprintf("func (%v *%v) baseValidation() error {\n",
		table.variableName(), table.fullyQualifiedStructName())

	for _, col := range table.ColumnList {
		switch col.GoDataType {
		case "string":
			// Non nullable string
			maxLen, lenCheckReasonComment := getMaxlenWithReasonCommentForStringColumn(col)
			if maxLen > 0 {
				tabCommonValidation += lenCheckReasonComment
				tabCommonValidation += fmt.Sprintf("if len(%v.%v) > %v {\n",
					table.variableName(), col.GoName, maxLen)
				tabCommonValidation += `return fmt.Errorf("E#` + newUniqueLmid() +
					fmt.Sprintf(` - Invalid length for %v.%v`,
						table.variableName(), col.GoName) +
					` %v", ` + fmt.Sprintf("len(%v.%v))\n",
					table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("}\n")
			}

			minLen := col.CommentProperties.MinStrLen
			if minLen > 0 {
				tabCommonValidation += fmt.Sprintf("if len(%v.%v) < %v {\n",
					table.variableName(), col.GoName, minLen)
				tabCommonValidation += `return fmt.Errorf("E#` + newUniqueLmid() +
					fmt.Sprintf(` - Length too less for %v.%v`,
						table.variableName(), col.GoName) +
					` %v", ` + fmt.Sprintf("len(%v.%v))\n",
					table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("}\n")
			}

			if col.CommentProperties.StrValidateAs == "email" {
				tabCommonValidation += fmt.Sprintf("if !types.IsValidEmail(%v.%v) {\n",
					table.variableName(), col.GoName)
				tabCommonValidation += `return fmt.Errorf("E#` + newUniqueLmid() +
					fmt.Sprintf(` - Value for %v.%v is expected to be a valid email")`,
						table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("}\n")
			}
		case "sql.NullString":
			// Nullable string.
			// Rule: A nullable string can be either null or have a max length specified
			maxLen, lenCheckReasonComment := getMaxlenWithReasonCommentForStringColumn(col)
			if maxLen > 0 {
				tabCommonValidation += fmt.Sprintf("if %v.%v.Valid {\n", table.variableName(), col.GoName)
				tabCommonValidation += lenCheckReasonComment
				tabCommonValidation += fmt.Sprintf("if len(%v.%v.String) > %v {\n",
					table.variableName(), col.GoName, col.CharacterLength)
				tabCommonValidation += `return fmt.Errorf("E#` + newUniqueLmid() +
					fmt.Sprintf(` - Invalid length for %v.%v`,
						table.variableName(), col.GoName) +
					` %v", ` + fmt.Sprintf("len(%v.%v.String))\n",
					table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("}\n")
				tabCommonValidation += fmt.Sprintf("}\n")
			}

			minLen := col.CommentProperties.MinStrLen
			if minLen > 0 {
				tabCommonValidation += fmt.Sprintf("if %v.%v.Valid {\n", table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("if len(%v.%v.String) < %v {\n",
					table.variableName(), col.GoName, minLen)
				tabCommonValidation += `return fmt.Errorf("E#` + newUniqueLmid() +
					fmt.Sprintf(` - Length too less for %v.%v`,
						table.variableName(), col.GoName) +
					` %v", ` + fmt.Sprintf("len(%v.%v.String))\n",
					table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("}\n")
				tabCommonValidation += fmt.Sprintf("}\n")
			}

			if col.CommentProperties.StrValidateAs == "email" {
				tabCommonValidation += fmt.Sprintf("if !types.IsValidEmail(%v.%v) {\n",
					table.variableName(), col.GoName)
				tabCommonValidation += `return fmt.Errorf("E#` + newUniqueLmid() +
					fmt.Sprintf(` - Value for %v.%v is expected to be a valid email")`,
						table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("}\n")
			}
		// TODO: Write cases for int64, sql.NullInt64 and other data types
		case "int64":
			minVal := col.CommentProperties.MinIntVal
			if minVal != 0 && minVal < math.MaxInt64 {
				tabCommonValidation += fmt.Sprintf("if %v.%v < int%v {\n",
					table.variableName(), col.GoName, minVal)
				tabCommonValidation += `return fmt.Errorf("E#` + newUniqueLmid() +
					fmt.Sprintf(` - Value too less for %v.%v`,
						table.variableName(), col.GoName) +
					` %v", ` + fmt.Sprintf("len(%v.%v))\n",
					table.variableName(), col.GoName)
				tabCommonValidation += fmt.Sprintf("}\n")
			}
		}
	}

	tabCommonValidation += "return nil\n"
	tabCommonValidation += "}\n\n"

	return tabCommonValidation, importList
}

func (g *Generator) buildTableCommonValidation(table DbTable, importList []string) (string, []string) {
	tabCommonValidation := ""
	tabCommonValidation += fmt.Sprintf("func (%v *%v) commonValidation() error {\n",
		table.variableName(), table.fullyQualifiedStructName())
	tabCommonValidation += fmt.Sprintf("err := %v.baseValidation()\n", lowerFirstChar(table.GoNameSingular))
	tabCommonValidation += "if err != nil{\nreturn err\n}\n"
	tabCommonValidation += "// More code to be written here for validation\n"
	tabCommonValidation += "return nil\n"
	tabCommonValidation += "}\n\n"

	return tabCommonValidation, importList
}

func (g *Generator) buildTableInsertValidation(table DbTable, importList []string) (string, []string) {
	tabInsertValidation := ""
	tabInsertValidation += fmt.Sprintf("func (%v *%v) validateForInsert() error {\n",
		table.variableName(), table.fullyQualifiedStructName())
	tabInsertValidation += fmt.Sprintf("err := %v.commonValidation()\n", table.variableName())
	tabInsertValidation += "if err != nil{\nreturn err\n}\n"
	tabInsertValidation += "// More code to be written here for validation\n"
	tabInsertValidation += "return nil\n"
	tabInsertValidation += "}\n\n"

	return tabInsertValidation, importList
}

func (g *Generator) buildTableUpdateValidation(table DbTable, importList []string) (string, []string) {
	tabUpdateValidation := ""
	tabUpdateValidation += fmt.Sprintf("func (%v *%v) validateForUpdate() error {\n",
		table.variableName(), table.fullyQualifiedStructName())
	tabUpdateValidation += fmt.Sprintf("err := %v.commonValidation()\n", lowerFirstChar(table.GoNameSingular))
	tabUpdateValidation += "if err != nil{\nreturn err\n}\n"
	tabUpdateValidation += "// More code to be written here for validation\n"
	tabUpdateValidation += "return nil\n"
	tabUpdateValidation += "}\n\n"

	return tabUpdateValidation, importList
}

func (g *Generator) buildTableInsertMethod(table DbTable, importList []string) (string, []string) {
	insertCode := ""
	insertCode += fmt.Sprintf("func (%v *%v) insert() error {\n",
		table.variableName(), table.fullyQualifiedStructName())
	insertCode += "var err error\n"
	insertCode += fmt.Sprintf("err = %v.validateForInsert()\n", table.variableName())
	insertCode += "if err != nil{\nreturn err\n}\n"

	insertCode += fmt.Sprintf("insertQuery := `INSERT INTO %v (\n", table.fullyQualifiedTableName())

	colNameSlice := []string{}
	argPositionSlice := []string{}
	pkeyColumnNameSlice := []string{}
	pkeyAmpInsertedColumnNameSlice := []string{}
	returningColsSlice := []DbColumn{}
	goColumnNameSlice := []string{}
	i := 0
	for _, column := range table.ColumnMap {
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

	for _, column := range table.ColumnMap {
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
		insertCode += fmt.Sprintf("%v.%v = inserted%v\n", table.variableName(), col.GoName, col.GoName)
	}
	for _, col := range returningColsSlice {
		insertCode += fmt.Sprintf("%v.%v = inserted%v\n", table.variableName(), col.GoName, col.GoName)
	}

	insertCode += "return nil"
	insertCode += "}\n\n"

	//importList = g.addToImports(g.Config.DbModelPackageName+"/resources", importList)
	importList = g.addToImports("fmt", importList)
	importList = g.addToImports("errors", importList)
	return insertCode, importList
}

func (g *Generator) buildTableUpdateMethod(table DbTable, importList []string) (string, []string) {
	updateCode := ""

	updateCode += fmt.Sprintf("func (%v *%v) update() error {\n",
		table.variableName(), table.fullyQualifiedStructName())

	if len(table.PkColumnList) == 0 {
		updateCode += "return errors.New(\"E#" + newUniqueLmid() + " - Cannot update " + table.fullyQualifiedTableName() + " because of no primary key\")\n"
		updateCode += "}\n"
		importList = g.addToImports("errors", importList)
		return updateCode, importList
	}

	updateCode += fmt.Sprintf("err := %v.validateForUpdate()\n", table.variableName())
	updateCode += "if err != nil{\nreturn err\n}\n"

	updateCode += fmt.Sprintf("updateQuery := `UPDATE %v SET\n\t\t\t", table.fullyQualifiedTableName())
	// Get the column information
	columnNameArgPositionPairCollection := []string{}
	goColumnNameCollection := []string{}

	i := 0
	for _, column := range table.ColumnList {
		if !columnInList(column.Name, table.PkColumnList) {
			if !(column.Name == "created_at" && // Created at timestamps are never to be updated.
				(column.GoDataType == "time.Time" || column.GoDataType == "sql.NullTime")) &&
				!(column.Name == "updated_at" && g.Config.UpdateUpdatedAtInCode == false &&
					(column.GoDataType == "time.Time" || column.GoDataType == "sql.NullTime")) {
				i += 1
				columnNameArgPositionPairCollection = append(columnNameArgPositionPairCollection, fmt.Sprintf(`"%v" = $%v`, column.Name, i))
				if column.DataType == "json" || column.DataType == "jsonb" {
					goColumnNameCollection = append(goColumnNameCollection, fmt.Sprintf("%v.%v.String()", table.variableName(), column.GoName))
				} else {
					goColumnNameCollection = append(goColumnNameCollection, fmt.Sprintf("%v.%v", table.variableName(), column.GoName))
				}
			}
		}
	}

	updateCode += strings.Join(columnNameArgPositionPairCollection, ",\n\t\t\t")
	updateCode += "\n\t\tWHERE \n\t\t\t"

	for k, column := range table.PkColumnList {
		comma := " AND "
		if k == len(table.PkColumnList)-1 {
			comma = ""
		}
		updateCode += fmt.Sprintf(`"%v" = $%v%v`, column.Name, i+1, comma)
		if column.DataType == "json" || column.DataType == "jsonb" {
			goColumnNameCollection = append(goColumnNameCollection, fmt.Sprintf("%v.%v.String()", table.variableName(), column.GoName))
		} else {
			goColumnNameCollection = append(goColumnNameCollection, fmt.Sprintf("%v.%v", table.variableName(), column.GoName))
		}
		i += 1
	}

	updateCode += "`\n\n"
	updateCode += fmt.Sprintf("_, err = %v.Exec(updateQuery,\n%v)\n", upperFirstChar(g.Config.DbModelPackageName), strings.Join(goColumnNameCollection, ",\n"))
	updateCode += "if err != nil {\n"
	updateCode += "return fmt.Errorf(\"E#" + newUniqueLmid() + " - Could not update " + table.fullyQualifiedStructName() + " in database: %v\", err)"
	updateCode += "}\n"

	// updateCode += `fmt.Printf("query: %v\n", updateQuery)`

	updateCode += "\nreturn nil\n"
	updateCode += "}\n\n"

	//importList = g.addToImports(baseGoModuleName+"/resources", importList)

	return updateCode, importList
}
