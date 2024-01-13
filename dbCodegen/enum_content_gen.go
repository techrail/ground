package dbCodegen

import "fmt"

func (g *Generator) buildEnumContentString(enum EnumDefinition, importList []string) (string, []string) {
	enumContentStr := ""

	enumTypeName := lowerFirstChar(enum.goNameSingular)
	if enum.Exported {
		enumTypeName = enum.goNameSingular
	}

	// Create the constants

	enumContentStr += fmt.Sprintf("type %v int16\n", enumTypeName)
	enumContentStr += "const(\n"
	enumContentStr += fmt.Sprintf("Undefined %v = -1\n", enumTypeName)
	for name, value := range enum.Mappings {
		enumContentStr += fmt.Sprintf("%v %v = %v\n", name, enumTypeName, value)
	}
	enumContentStr += ")\n"

	// Type to string
	enumContentStr += fmt.Sprintf("func (t %v) String() string {\n", enumTypeName)
	enumContentStr += fmt.Sprintf("switch(t) {\n")
	for name, _ := range enum.Mappings {
		enumContentStr += fmt.Sprintf("case %v: \n return \"%v\"\n", name, name)
	}
	enumContentStr += fmt.Sprintf("case Undefined: \n fallthrough \n")
	enumContentStr += fmt.Sprintf("default: \n return \"Undefined\"\n")
	enumContentStr += fmt.Sprintf("}\n")
	enumContentStr += fmt.Sprintf("}\n")

	enumContentStr += "\n"

	// String to Type
	if enum.Exported {
		enumContentStr += fmt.Sprintf("func StringTo%v(input string) %v {\n", enumTypeName, enumTypeName)
	} else {
		enumContentStr += fmt.Sprintf("func stringTo%v(input string) %v {\n", enum.goNameSingular, enumTypeName)
	}

	enumContentStr += fmt.Sprintf("switch(input) {\n")
	for name, _ := range enum.Mappings {
		enumContentStr += fmt.Sprintf("case \"%v\": \n return %v\n", name, name)
	}
	enumContentStr += fmt.Sprintf("case \"Undefined\": \n fallthrough \n")
	enumContentStr += fmt.Sprintf("default: \n return Undefined\n")
	enumContentStr += fmt.Sprintf("}\n")
	enumContentStr += "}\n"
	enumContentStr += "\n"

	// Int16 to Type
	enumContentStr += fmt.Sprintf("func %vFromInt16(input int16) %v {\n", enumTypeName, enumTypeName)
	inputComparisons := []string{}
	for _, intVal := range enum.Mappings {
		inputComparisons = append(inputComparisons, fmt.Sprintf("input == %v", intVal))
	}
	enumContentStr += fmt.Sprintf("if %v {", groupBy3(inputComparisons, " || ", "\n\t"))
	enumContentStr += fmt.Sprintf("return %v(input)\n", enumTypeName)
	enumContentStr += fmt.Sprintf("}\n")
	enumContentStr += fmt.Sprintf("return Undefined\n")
	enumContentStr += fmt.Sprintf("}\n")
	enumContentStr += "\n"

	return enumContentStr, importList
}
