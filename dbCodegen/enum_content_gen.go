package dbCodegen

import "fmt"

func (g *Generator) buildEnumContentString(enum EnumDefinition, importList []string) (string, []string) {
	enumContentStr := ""
	enumContentStr += fmt.Sprintf("type %v int16\n", enum.goNameSingular)
	enumContentStr += "const(\n"
	enumContentStr += fmt.Sprintf("Undefined %v = -1\n", enum.goNameSingular)
	for name, value := range enum.Mappings {
		enumContentStr += fmt.Sprintf("%v %v = %v\n", name, enum.goNameSingular, value)
	}
	enumContentStr += ")\n"

	enumContentStr += fmt.Sprintf("func (t %v) String() string {\n", enum.goNameSingular)
	enumContentStr += fmt.Sprintf("switch(t) {\n")
	for name, _ := range enum.Mappings {
		enumContentStr += fmt.Sprintf("case %v: \n return \"%v\"\n", name, name)
	}
	enumContentStr += fmt.Sprintf("case Undefined: \n fallthrough \n")
	enumContentStr += fmt.Sprintf("default: \n return \"Undefined\"\n")

	enumContentStr += fmt.Sprintf("}\n")
	enumContentStr += fmt.Sprintf("}\n")

	enumContentStr += "\n"

	return enumContentStr, importList
}
