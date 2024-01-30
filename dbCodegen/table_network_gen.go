package dbCodegen

import (
	"fmt"
)

func (g *Generator) buildNetworkStructString(table DbTable, importList []string) (string, []string) {
	var networkStruct string
	networkStruct = ""

	// tableName := getGoName(table.Name)
	networkStruct += fmt.Sprintf("// %vForResponse struct represents the %v table in %v schema of the DB for network responses\n",
		table.fullyQualifiedStructName(), table.Name, table.Schema)
	networkStruct += fmt.Sprintf("// Table Comment: %v\n", table.commentForStruct())
	networkStruct += fmt.Sprintf("type %sForResponse struct {\n", table.fullyQualifiedStructName())
	for _, columnName := range table.ColumnListA2z {
		column := table.ColumnMap[columnName]
		if !column.CommentProperties.HideFromNetwork {
			if column.CommentProperties.StrConversionViaEnum != "" {
				if column.Nullable {
					networkStruct += fmt.Sprintf("\t%s *string `json:\"%s\"`\n",
						column.GoName, lowerFirstChar(column.GoName))
				} else {
					networkStruct += fmt.Sprintf("\t%s string `json:\"%s\"`\n",
						column.GoName, lowerFirstChar(column.GoName))
				}
			} else {
				networkStruct += fmt.Sprintf("\t%s %s `json:\"%s\"`\n",
					column.GoName, column.NetworkDataType, lowerFirstChar(column.GoName))
				if g.getGoImportForDataType(column.DataType, column.Nullable) != "database/sql" {
					importList = g.addToImports(g.getGoImportForDataType(column.DataType, column.Nullable), importList)
				}
			}
		}
	}
	networkStruct += "}\n\n"

	networkStruct += fmt.Sprintf(
		"// New%vForResponse will return pointer to a new %vForResponse type\n",
		table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("func New%vForResponse() *%vForResponse {\n",
		table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("return &%vForResponse{}\n", table.fullyQualifiedStructName())
	networkStruct += "}\n"

	networkStruct += fmt.Sprintf(
		"// FillFromDb%v will copy data from a %v.%v type to a %vForResponse type\n",
		table.fullyQualifiedStructName(), g.Config.DbModelPackageName, table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("func (%vForResponse *%vForResponse) FillFromDb%v(db%v %v.%v) %vForResponse {\n",
		table.variableName(), table.fullyQualifiedStructName(), table.fullyQualifiedStructName(),
		table.GoNameSingular, g.Config.DbModelPackageName, table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	for _, colName := range table.ColumnListA2z {
		col := table.ColumnMap[colName]
		if !col.CommentProperties.HideFromNetwork {
			if col.GoDataType == col.NetworkDataType {
				if col.CommentProperties.StrConversionViaEnum != "" {
					enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
					if err != nil {
						panic(fmt.Sprintf("P#1RDB3A - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
					}

					networkStruct += fmt.Sprintf("%vForResponse.%v = %v.%vFromInt16(db%v.%v).String()\n",
						table.variableName(), col.GoName,
						col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
						table.GoNameSingular, col.GoName)
					importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
				} else {
					networkStruct += fmt.Sprintf("%vForResponse.%v = db%v.%v\n",
						table.variableName(), col.GoName,
						table.GoNameSingular, col.GoName)
				}
			} else {
				networkStruct += fmt.Sprintf("%vForResponse.%v = nil\n",
					table.variableName(), col.GoName)
				if col.GoDataType == "sql.NullString" {
					networkStruct += fmt.Sprintf("if db%v.%v.Valid {\n", table.GoNameSingular, col.GoName)
					if col.CommentProperties.StrConversionViaEnum != "" {
						enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
						if err != nil {
							panic(fmt.Sprintf("P#1RDC5X - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
						}

						networkStruct += fmt.Sprintf("val := %v.%vFromInt16(db%v.%v.String).String()\n",
							col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
							table.GoNameSingular, col.GoName)
						networkStruct += fmt.Sprintf("%vForResponse.%v = &val\n",
							table.variableName(), col.GoName)
						importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
					} else {
						networkStruct += fmt.Sprintf("%vForResponse.%v = &db%v.%v.String\n",
							table.variableName(), col.GoName,
							table.GoNameSingular, col.GoName)
					}
					networkStruct += "}\n"
				}
				if col.GoDataType == "sql.NullInt16" {
					networkStruct += fmt.Sprintf("if db%v.%v.Valid {\n", table.GoNameSingular, col.GoName)
					if col.CommentProperties.StrConversionViaEnum != "" {
						enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
						if err != nil {
							panic(fmt.Sprintf("P#1RDC64 - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
						}

						networkStruct += fmt.Sprintf("val := %v.%vFromInt16(db%v.%v.Int16).String()\n",
							col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
							table.GoNameSingular, col.GoName)
						networkStruct += fmt.Sprintf("%vForResponse.%v = &val\n",
							table.variableName(), col.GoName)
						importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
					} else {
						networkStruct += fmt.Sprintf("%vForResponse.%v = &db%v.%v.Int16\n",
							table.variableName(), col.GoName,
							table.GoNameSingular, col.GoName)
					}
					networkStruct += "}\n"
				}
				if col.GoDataType == "sql.NullInt32" {
					networkStruct += fmt.Sprintf("if db%v.%v.Valid {\n", table.GoNameSingular, col.GoName)
					if col.CommentProperties.StrConversionViaEnum != "" {
						enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
						if err != nil {
							panic(fmt.Sprintf("P#1RDC69 - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
						}

						networkStruct += fmt.Sprintf("val := %v.%vFromInt16(db%v.%v.Int32).String()\n",
							col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
							table.GoNameSingular, col.GoName)
						networkStruct += fmt.Sprintf("%vForResponse.%v = &val\n",
							table.variableName(), col.GoName)
						importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
					} else {
						networkStruct += fmt.Sprintf("%vForResponse.%v = &db%v.%v.Int32\n",
							table.variableName(), col.GoName,
							table.GoNameSingular, col.GoName)
					}
					networkStruct += "}\n"
				}
				if col.GoDataType == "sql.NullInt64" {
					networkStruct += fmt.Sprintf("if db%v.%v.Valid {\n", table.GoNameSingular, col.GoName)
					if col.CommentProperties.StrConversionViaEnum != "" {
						enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
						if err != nil {
							panic(fmt.Sprintf("P#1RDC6G - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
						}

						networkStruct += fmt.Sprintf("val := %v.%vFromInt16(db%v.%v.Int64).String()\n",
							col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
							table.GoNameSingular, col.GoName)
						networkStruct += fmt.Sprintf("%vForResponse.%v = &val\n",
							table.variableName(), col.GoName)
						importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
					} else {
						networkStruct += fmt.Sprintf("%vForResponse.%v = &db%v.%v.Int64\n",
							table.variableName(), col.GoName,
							table.GoNameSingular, col.GoName)
					}
					networkStruct += "}\n"
				}
				if col.GoDataType == "sql.NullFloat64" {
					networkStruct += fmt.Sprintf("if db%v.%v.Valid {\n", table.GoNameSingular, col.GoName)
					if col.CommentProperties.StrConversionViaEnum != "" {
						enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
						if err != nil {
							panic(fmt.Sprintf("P#1RDC6R - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
						}

						networkStruct += fmt.Sprintf("val := %v.%vFromInt16(db%v.%v.Float64).String()\n",
							col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
							table.GoNameSingular, col.GoName)
						networkStruct += fmt.Sprintf("%vForResponse.%v = &val\n",
							table.variableName(), col.GoName)
						importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
					} else {
						networkStruct += fmt.Sprintf("%vForResponse.%v = &db%v.%v.Float64\n",
							table.variableName(), col.GoName,
							table.GoNameSingular, col.GoName)
					}
					networkStruct += "}\n"
				}
				if col.GoDataType == "sql.NullBool" {
					networkStruct += fmt.Sprintf("if db%v.%v.Valid {\n", table.GoNameSingular, col.GoName)
					if col.CommentProperties.StrConversionViaEnum != "" {
						enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
						if err != nil {
							panic(fmt.Sprintf("P#1RDC6X - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
						}

						networkStruct += fmt.Sprintf("val := %v.%vFromInt16(db%v.%v.Bool).String()\n",
							col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
							table.GoNameSingular, col.GoName)
						networkStruct += fmt.Sprintf("%vForResponse.%v = &val\n",
							table.variableName(), col.GoName)
						importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
					} else {
						networkStruct += fmt.Sprintf("%vForResponse.%v = &db%v.%v.Bool\n",
							table.variableName(), col.GoName,
							table.GoNameSingular, col.GoName)
					}
					networkStruct += "}\n"
				}
				if col.GoDataType == "sql.NullTime" {
					networkStruct += fmt.Sprintf("if db%v.%v.Valid {\n", table.GoNameSingular, col.GoName)
					if col.CommentProperties.StrConversionViaEnum != "" {
						enum, err := g.getEnumByName(col.CommentProperties.StrConversionViaEnum)
						if err != nil {
							panic(fmt.Sprintf("P#1RDC72 - Expected enum %v to be found but did not.", col.CommentProperties.StrConversionViaEnum))
						}

						networkStruct += fmt.Sprintf("val := %v.%vFromInt16(db%v.%v.Time).String()\n",
							col.CommentProperties.StrConversionViaEnum, enum.goTypeName,
							table.GoNameSingular, col.GoName)
						networkStruct += fmt.Sprintf("%vForResponse.%v = &val\n",
							table.variableName(), col.GoName)
						importList = g.addToImports(g.Config.DbModelPackageName+"/types/"+col.CommentProperties.StrConversionViaEnum, importList)
					} else {
						networkStruct += fmt.Sprintf("%vForResponse.%v = &db%v.%v.Time\n",
							table.variableName(), col.GoName,
							table.GoNameSingular, col.GoName)
					}
					networkStruct += "}\n"
				}
			}
		}
	}

	networkStruct += fmt.Sprintf("\nreturn *%vForResponse\n", table.variableName())
	networkStruct += "}\n"

	networkStruct += fmt.Sprintf(
		"// FillSliceFromDb%v will copy data from a slice of %v.%v type to a slice of %vForResponse type\n",
		table.fullyQualifiedStructName(), g.Config.DbModelPackageName, table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("func FillSliceFromDb%v(db%v []%v.%v) []%vForResponse {\n",
		table.fullyQualifiedStructName(), table.GoNameSingular, g.Config.DbModelPackageName, table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("networkSlice := make([]%vForResponse, 0)\n", table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("for _, loop%v := range db%v {\n", table.GoNameSingular, table.GoNameSingular)
	networkStruct += fmt.Sprintf("networkSlice = append(networkSlice, New%vForResponse().FillFromDb%v(loop%v))\n",
		table.fullyQualifiedStructName(), table.fullyQualifiedStructName(), table.GoNameSingular)
	networkStruct += "}\n"
	networkStruct += "return networkSlice\n"
	networkStruct += "}\n"

	networkStruct += fmt.Sprintf(
		"// FillSliceFromDb%vPointers will copy data from a slice of pointers of %v.%v type to a slice of %vForResponse type\n",
		table.fullyQualifiedStructName(), g.Config.DbModelPackageName, table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("func FillSliceFromDb%vPointers(db%v []*%v.%v) []%vForResponse {\n",
		table.fullyQualifiedStructName(), table.GoNameSingular, g.Config.DbModelPackageName, table.fullyQualifiedStructName(), table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("networkSlice := make([]%vForResponse, 0)\n", table.fullyQualifiedStructName())
	networkStruct += fmt.Sprintf("for _, loop%v := range db%v {\n", table.GoNameSingular, table.GoNameSingular)
	networkStruct += fmt.Sprintf("networkSlice = append(networkSlice, New%vForResponse().FillFromDb%v(*loop%v))\n",
		table.fullyQualifiedStructName(), table.fullyQualifiedStructName(), table.GoNameSingular)
	networkStruct += "}\n"
	networkStruct += "return networkSlice\n"
	networkStruct += "}\n"

	importList = g.addToImports(fmt.Sprintf("%s/%s", g.Config.ModelsContainerPackage, g.Config.DbModelPackageName), importList)

	return networkStruct, importList
}
