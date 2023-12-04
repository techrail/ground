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

	tableUpdateByIndexes := ""
	singleIdxUpdateFuncStr := ""
	if g.Config.BuildUpdateByUniqueIndex == true {
		for _, idx := range table.IndexList {
			if idx.IsUnique {
				singleIdxUpdateFuncStr, importList = g.buildTableUpdateMethodBySingleIndex(table, idx, importList)
				tableUpdateByIndexes += singleIdxUpdateFuncStr
			}
		}
	}

	tableDeleteFuncStr := ""
	tableDeleteFuncStr, importList = g.buildTableDeleteMethod(table, importList)

	// Foreign key funcs
	tabForeignKeyMethods := ""

	for _, fkey := range table.FKeyMap {
		tabSingleFkeyMethod, iList := g.buildSingleTableFkeyFunc(table, fkey, importList)
		tabForeignKeyMethods += tabSingleFkeyMethod + "\n\n"
		importList = iList
	}

	// Now the Dao Ones
	tableMethodAndDaoSeparator := "\n// ============================================="
	tableMethodAndDaoSeparator += "\n// Table methods end here. Dao functions below"
	tableMethodAndDaoSeparator += "\n// =============================================\n\n"
	// TODO: Write the Dao functions
	tableDaoStructAndNew := ""
	tableDaoStructAndNew, importList = g.buildTableDaoStructAndNewFunc(table, importList)

	tableDaoFunctions := ""
	tableDaoFunctions, importList = g.buildTableDaoIdxFuncCreator(table, importList)

	tableBaseFuncStr := tableBaseValidationFuncStr + tableCommonValidationFuncStr + tableInsertionValidationFuncStr +
		tableUpdateValidationFuncStr + tableInsertionFuncStr + tableUpdateFuncStr +
		tableUpdateByIndexes + tableDeleteFuncStr + tabForeignKeyMethods +
		tableMethodAndDaoSeparator + tableDaoStructAndNew + tableDaoFunctions
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
	if len(pkeyColumnNameSlice) > 0 || len(returningColsSlice) > 0 {
		insertCode += "RETURNING "
		insertCode += strings.Join(pkeyColumnNameSlice, ", ")
		for _, col := range returningColsSlice {
			insertCode += `, "` + col.Name + `"`
		}
	}
	insertCode += "`;\n\n"

	//
	insertCode += fmt.Sprintf("resultRow := %v.QueryRowx(insertQuery,\n", upperFirstChar(g.Config.DbModelPackageName))
	insertCode += strings.Join(goColumnNameSlice, ", \n\t\t\t")
	insertCode += ",\n)\n\n"
	insertCode += "if resultRow.Err() != nil {\n"
	insertCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not insert into database: %v\", resultRow.Err())\n"
	importList = g.addToImports("fmt", importList)
	importList = g.addToImports("errors", importList)
	insertCode += "logger.Println(errMsg)\n"
	importList = g.addToImports("github.com/techrail/ground/logger", importList)
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

	if len(pkeyColumnNameSlice) > 0 || len(returningColsSlice) > 0 {
		insertCode += fmt.Sprintf("\nerr = resultRow.Scan(%v)\n", returningColList)

		insertCode += "if err != nil {\n"
		insertCode += "return fmt.Errorf(\"E#" + newUniqueLmid() + " - Scan failed. Error: %v\", err)\n"
		insertCode += "}\n\n"
	}

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

func (g *Generator) buildTableUpdateMethodBySingleIndex(table DbTable, index DbIndex, importList []string) (string, []string) {
	if !index.IsUnique {
		return "", importList
	}

	updateCode := ""

	updateCode += fmt.Sprintf("func (%v *%v) updateBy%v() error {\n",
		table.variableName(), table.fullyQualifiedStructName(), index.GetFuncNamePart())

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

	for k, column := range index.ColumnList {
		comma := " AND "
		if k == len(index.ColumnList)-1 {
			comma = ""
		}
		updateCode += fmt.Sprintf(`"%v" = $%v%v`, column.Name, i+1, comma)
		if column.DataType == "json" {
			fmt.Printf("E#1OJ6NP - It is not possible to compare two JSON values while JSOB values are comparable, in fact. For column: %v\n", column.fullyQualifiedColumnName())
			return "", importList
		}

		if column.DataType == "jsonb" {
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

func (g *Generator) buildTableUpdateMethod(table DbTable, importList []string) (string, []string) {
	updateCode := ""

	updateCode += fmt.Sprintf("func (%v *%v) update() error {\n",
		table.variableName(), table.fullyQualifiedStructName())

	if len(table.PkColumnList) == 0 {
		updateCode += "return errors.New(\"E#" + newUniqueLmid() + " - Cannot update " + table.fullyQualifiedTableName() + " because of no primary key. Please write update query yourself\")\n"
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
		if column.DataType == "json" {
			fmt.Printf("E#1OJ6NP - It is not possible to compare two JSON values while JSOB values are comparable, in fact. For column: %v\n", column.fullyQualifiedColumnName())
			return "", importList
		}

		if column.DataType == "jsonb" {
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

func (g *Generator) buildTableDeleteMethod(table DbTable, importList []string) (string, []string) {
	deleteCode := ""
	deleteCode += fmt.Sprintf("func (%v *%v) delete() error {\n",
		table.variableName(), table.fullyQualifiedStructName())

	if len(table.PkColumnList) == 0 {
		deleteCode += "return errors.New(\"E#" + newUniqueLmid() + " - Cannot delete " + table.fullyQualifiedTableName() + " because of no primary key. Please write deletion query yourself\")\n"
		deleteCode += "}\n"
		importList = g.addToImports("errors", importList)
		return deleteCode, importList
	}

	deleteCode += fmt.Sprintf("_, err := %v.Exec(`DELETE FROM %v WHERE ",
		upperFirstChar(g.Config.DbModelPackageName), table.fullyQualifiedTableName())
	// id = $1`, user.Id)")
	pks := []string{}
	for k, column := range table.PkColumnList {
		comma := " AND "
		if k == len(table.PkColumnList)-1 {
			comma = ""
		}
		pks = append(pks, fmt.Sprintf("%v.%v", lowerFirstChar(table.GoNameSingular), column.GoName))
		deleteCode += fmt.Sprintf("%v = $%v%v", column.Name, k+1, comma)
	}
	deleteCode += "`," + strings.Join(pks, ", ") + ")\n"
	deleteCode += "if err != nil {\n"
	deleteCode += "return fmt.Errorf(\"E#" + newUniqueLmid() + " - Could not delete " + table.GoNameSingular + " from database: %v\", err)"
	deleteCode += "}\n"

	deleteCode += "\nreturn nil\n"
	deleteCode += "}\n\n"

	//importList = addToImports(baseGoModuleName+"/resources", importList)

	return deleteCode, importList
}

func (g *Generator) buildSingleTableFkeyFunc(table DbTable, fkey DbFkInfo, importList []string) (string, []string) {
	tabFKeyMethod := ""

	// The target table should be there
	_, ok := g.Schemas[fkey.ToSchema]
	if !ok {
		panic(fmt.Sprintf("E#1OJ0XX - Expected the toSchema %v to be there but was not", fkey.ToSchema))
	}
	targetTable, ok := g.Schemas[fkey.ToSchema].Tables[fkey.ToTable]
	if !ok {
		panic(fmt.Sprintf("E#1OJ13B - Expected the toTable %v in toSchema %v to be there but was not", fkey.ToTable, fkey.ToSchema))
	}

	funcNamePart := ""
	// funcArgs := make([]string, 0)
	queryValPairs := make([]string, 0)
	queryVars := make([]string, 0)
	i := 1
	for fromColName, toColName := range fkey.References {
		fromCol, err := g.getColumnFromListByName(fromColName, table.ColumnList)

		if err != nil {
			panic(fmt.Sprintf("E#1OJ1JG - Expected column %v is not prsent in table %v in schema %v",
				fromColName, table.Name, table.Schema))
		}

		funcNamePart += fromCol.GoName
		queryValPairs = append(queryValPairs, fmt.Sprintf("%v = $%v", toColName, i))
		i += 1
		if fromCol.Nullable && fromCol.GoDataType != "interface{}" {
			switch fromCol.GoDataType {
			case "sql.NullInt64":
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName+".Int64")
			case "sql.NullInt32":
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName+".Int32")
			case "sql.NullInt16":
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName+".Int16")
			case "sql.NullFloat64":
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName+".Float64")
			case "sql.NullBool":
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName+".Bool")
			case "sql.NullString":
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName+".String")
			case "types.JsonObject":
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName+".String()")
			default:
				queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName)
			}
		} else {
			queryVars = append(queryVars, lowerFirstChar(table.GoNameSingular)+"."+fromCol.GoName)
		}
	}
	fmt.Println("E#1C7C24 -", funcNamePart)

	tabFKeyMethod += fmt.Sprintf("func (%v *%v) Get%vFromDbBy%v(getFromMainDb ...bool) (%v, error) {\n",
		table.variableName(), table.fullyQualifiedStructName(), targetTable.GoNameSingular, funcNamePart, targetTable.fullyQualifiedStructName())
	tabFKeyMethod += "var err error\n"
	tabFKeyMethod += fmt.Sprintf("query := `SELECT * FROM %v WHERE %v;`\n", targetTable.fullyQualifiedTableName(), strings.Join(queryValPairs, " AND "))
	tabFKeyMethod += fmt.Sprintf("linked%v := %v{}\n\n", targetTable.GoNameSingular, targetTable.fullyQualifiedStructName())

	tabFKeyMethod += "if len(getFromMainDb) > 0 && getFromMainDb[0] == true {\n"
	tabFKeyMethod += fmt.Sprintf("err = %v.Get(&linked%v, query, %v)\n", upperFirstChar(g.Config.DbModelPackageName), targetTable.GoNameSingular, strings.Join(queryVars, ", "))
	tabFKeyMethod += "} else {\n"
	tabFKeyMethod += fmt.Sprintf("err = %vReader.Get(&linked%v, query, %v)\n", upperFirstChar(g.Config.DbModelPackageName), targetTable.GoNameSingular, strings.Join(queryVars, ", "))
	tabFKeyMethod += "}\n"
	tabFKeyMethod += "\nif err == sql.ErrNoRows {\n"
	importList = g.addToImports("database/sql", importList)
	tabFKeyMethod += fmt.Sprintf("return linked%v, err\n", targetTable.GoNameSingular)
	tabFKeyMethod += "}\n\n"

	tabFKeyMethod += "if err != nil {\n"
	tabFKeyMethod += `errMsg := fmt.Sprintf("E#` + newUniqueLmid() + ` - Could not load ` + ` by Id Error: %v", err)` + "\n"
	tabFKeyMethod += "logger.Println(errMsg)\n"
	importList = g.addToImports("github.com/techrail/ground/logger", importList)
	tabFKeyMethod += fmt.Sprintf("return linked%v, errors.New(errMsg)\n", targetTable.GoNameSingular)
	tabFKeyMethod += "}\n"

	tabFKeyMethod += fmt.Sprintf("return linked%v, nil\n", targetTable.GoNameSingular)
	tabFKeyMethod += "}\n"

	return tabFKeyMethod, importList
}

func (g *Generator) buildTableDaoStructAndNewFunc(table DbTable, importList []string) (string, []string) {
	daoCode := ""

	daoCode += fmt.Sprintf("\ntype %vDao struct{}\n", table.fullyQualifiedStructName())
	daoCode += fmt.Sprintf("func New%vDao() *%vDao {\n", table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	daoCode += fmt.Sprintf("return &%vDao{}\n", table.fullyQualifiedStructName())
	daoCode += "}\n"

	return daoCode, importList
}

func (g *Generator) buildTableDaoIdxFuncCreator(table DbTable, importList []string) (string, []string) {
	daoIdxCode := ""
	for _, idx := range table.IndexList {
		daoSingleIdxCode, iList := g.buildSingleTableDaoIdxFunc(table, idx, importList)
		daoIdxCode += daoSingleIdxCode + "\n\n"
		importList = iList
	}

	return daoIdxCode, importList
}

func (g *Generator) buildSingleTableDaoIdxFunc(table DbTable, idx DbIndex, importList []string) (string, []string) {
	daoSingleIdxCode := ""
	// Columns in IndexList
	funcNamePart := ""
	argList := ""
	argListWithoutTypes := []string{}
	for i, col := range idx.ColumnList {
		funcNamePart += col.GoName
		argList += lowerFirstChar(col.GoNameSingular) + " " + col.GoDataType
		argListWithoutTypes = append(argListWithoutTypes, lowerFirstChar(col.GoNameSingular))
		if i+1 != len(idx.ColumnList) {
			argList += ","
		}
		i += 1
	}

	if idx.IsUnique {
		// Create a function to get a single item
		daoSingleIdxCode += fmt.Sprintf("func (%vDao *%vDao)GetFromDbBy%v(%v,getFromMainDb ...bool) (%v, error) {\n",
			lowerFirstChar(table.GoNameSingular), table.GoNameSingular, funcNamePart, argList, table.GoNameSingular)
		daoSingleIdxCode += "var err error\n"

		// Create the query now
		daoSingleIdxCode += fmt.Sprintf("query:=`SELECT * FROM %v WHERE ", table.Name)
		for k, column := range idx.ColumnList {
			comma := " AND "
			if k == len(idx.ColumnList)-1 {
				comma = ""
			}
			daoSingleIdxCode += fmt.Sprintf(`"%v" = $%v%v`, column.Name, k+1, comma)
			k += 1
		}
		daoSingleIdxCode += "`\n"

		daoSingleIdxCode += fmt.Sprintf("%v := %v{}\n", lowerFirstChar(table.GoNameSingular), table.GoNameSingular)

		daoSingleIdxCode += "\nif len(getFromMainDb) > 0 && getFromMainDb[0] == true {\n"
		daoSingleIdxCode += fmt.Sprintf("err = resources.MainDb.Get(&%v, query, %v)\n", lowerFirstChar(table.GoNameSingular), strings.Join(argListWithoutTypes, ", "))
		daoSingleIdxCode += "} else {\n"
		daoSingleIdxCode += fmt.Sprintf("err = resources.ReaderDb.Get(&%v, query, %v)\n", lowerFirstChar(table.GoNameSingular), strings.Join(argListWithoutTypes, ", "))
		daoSingleIdxCode += "}\n\n"

		daoSingleIdxCode += "if err==sql.ErrNoRows {\n"
		daoSingleIdxCode += fmt.Sprintf("return %v, err\n", lowerFirstChar(table.GoNameSingular))
		daoSingleIdxCode += "}\n"
		importList = g.addToImports("database/sql", importList)

		daoSingleIdxCode += "if err!=nil {\n"
		daoSingleIdxCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not load " + table.GoName + " by " + funcNamePart + " Error: %v\", err)\n"
		daoSingleIdxCode += "resources.Logger.LogString(errMsg)\n"
		daoSingleIdxCode += fmt.Sprintf("return %v, errors.New(errMsg)\n", lowerFirstChar(table.GoNameSingular))
		daoSingleIdxCode += "}\n"

		daoSingleIdxCode += fmt.Sprintf("return %v, nil\n", lowerFirstChar(table.GoNameSingular))
		daoSingleIdxCode += "}\n "
	} else {
		// Create a function to get a list of items
		daoSingleIdxCode += fmt.Sprintf("// GetListFromDbBy%v fetches a list of %v items from DB using given parameters\n",
			funcNamePart, table.GoNameSingular)
		daoSingleIdxCode += "// NOTE: This function does not implement pagination.\n"
		daoSingleIdxCode += fmt.Sprintf("func (%vDao *%vDao)GetListFromDbBy%v(%v,getFromMainDb ...bool) ([]*%v, error) {\n",
			lowerFirstChar(table.GoNameSingular), table.GoNameSingular, funcNamePart, argList, table.GoNameSingular)

		daoSingleIdxCode += "var err error\n"

		// Create the query now
		daoSingleIdxCode += fmt.Sprintf("query:=`SELECT * FROM %v WHERE ", table.Name)
		for k, column := range idx.ColumnList {
			comma := " AND "
			if k == len(idx.ColumnList)-1 {
				comma = ""
			}
			daoSingleIdxCode += fmt.Sprintf(`"%v" = $%v%v`, column.Name, k+1, comma)
			k += 1
		}
		daoSingleIdxCode += "`\n"

		daoSingleIdxCode += fmt.Sprintf("%v := make([]*%v, 0)\n", lowerFirstChar(table.GoNamePlural), table.GoNameSingular)
		daoSingleIdxCode += "\nif len(getFromMainDb) > 0 && getFromMainDb[0] == true {\n"
		daoSingleIdxCode += fmt.Sprintf("err = resources.MainDb.Select(&%v, query, %v)\n", lowerFirstChar(table.GoNamePlural), strings.Join(argListWithoutTypes, ", "))
		daoSingleIdxCode += "} else {\n"
		daoSingleIdxCode += fmt.Sprintf("err = resources.ReaderDb.Select(&%v, query, %v)\n", lowerFirstChar(table.GoNamePlural), strings.Join(argListWithoutTypes, ", "))
		daoSingleIdxCode += "}\n\n"

		daoSingleIdxCode += "if err!=nil {\n"
		daoSingleIdxCode += "errMsg := fmt.Sprintf(\"E#" + newUniqueLmid() + " - Could not load " + table.GoName + " by " + funcNamePart + " Error: %v\", err)\n"
		daoSingleIdxCode += "resources.Logger.LogString(errMsg)\n"
		daoSingleIdxCode += "return nil, errors.New(errMsg)\n"
		daoSingleIdxCode += "}\n"

		daoSingleIdxCode += fmt.Sprintf("return %v, nil\n", lowerFirstChar(table.GoNamePlural))
		daoSingleIdxCode += "}\n "
	}

	return daoSingleIdxCode, importList
}
