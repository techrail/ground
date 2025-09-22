package dbCodegen

import (
	"fmt"
)

func (g *Generator) buildEnumContentString(enum EnumDefinition, importList []string) (string, []string) {
	enumContentStr := ""

	enumTypeName := lowerFirstChar(enum.goNameSingular)
	if enum.Exported {
		enumTypeName = enum.goNameSingular
	}

	// Create the constants

	enumContentStr += fmt.Sprintf("type %v int16\n", enumTypeName)
	enumContentStr += "const(\n"
	enumContentStr += fmt.Sprintf("Undefined%v %v = -1\n", enumTypeName, enumTypeName)
	for name, value := range enum.Mappings {
		enumContentStr += fmt.Sprintf("%v %v = %v\n", name, enumTypeName, value)
	}
	enumContentStr += ")\n"

	// Type to string
	enumContentStr += fmt.Sprintf("func (t %v) String() string {\n", enumTypeName)
	enumContentStr += fmt.Sprintf("switch(t) {\n")
	for name := range enum.Mappings {
		enumContentStr += fmt.Sprintf("case %v: \n return \"%v\"\n", name, name)
	}
	enumContentStr += fmt.Sprintf("case Undefined%v: \n fallthrough \n", enumTypeName)
	enumContentStr += fmt.Sprintf("default: \n return \"Undefined%v\"\n", enumTypeName)
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
	for name := range enum.Mappings {
		enumContentStr += fmt.Sprintf("case \"%v\": \n return %v\n", name, name)
	}
	enumContentStr += fmt.Sprintf("case \"Undefined%v\": \n fallthrough \n", enumTypeName)
	enumContentStr += fmt.Sprintf("default: \n return Undefined%v\n", enumTypeName)
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
	enumContentStr += fmt.Sprintf("return Undefined%v\n", enumTypeName)
	enumContentStr += fmt.Sprintf("}\n")
	enumContentStr += "\n"

	// DB Methods
	if enum.IsDbType {
		// Value method
		importList = g.addToImports("database/sql/driver", importList)

		enumContentStr += fmt.Sprintf("func (t %v) Value() (driver.Value, error) { \n", enumTypeName)
		enumContentStr += fmt.Sprintf("if %vFromInt16(int16(t)) != Undefined%v {\n", enumTypeName, enumTypeName)
		enumContentStr += "return int16(t), nil \n"
		enumContentStr += "}\n"
		enumContentStr += "return -1, errors.New(\"E#" + newUniqueLmid() + " - Invalid value supplied for enumeration " + enumTypeName + "\")\n"

		importList = g.addToImports("errors", importList)

		enumContentStr += fmt.Sprintf("}\n")
		enumContentStr += "\n"
		enumContentStr += "\n // Scan method is not generated. Refer to NOTE 1NQKGL (search for that string) for details"

		// NOTE [1NQKGL] There is no Scan method that we would generate. The reason is that Scan requires a pointer receiver.
		// And a receiver function cannot update the receiver for null pointers to scalar values.
		// If an enumeration is being used as a database column and it is nullable then the corresponding
		// sql type (sql.NullInt16) would actually be a struct and structs can't be constants in go. As such
		// it is much easier to just not let enums be null.
		// If a null is required to be handled, then for an enumerated column, we can create a separate value which
		// can be considered as null on the developer/application level. In addition,
		// TODO: When building the code generator for the enumerated column in DB, make sure to check that the column
		//   is NOT NULL!
	}

	return enumContentStr, importList
}
