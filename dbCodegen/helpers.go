package dbCodegen

import (
	"github.com/techrail/ground/typs/integer"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var baseLmidSeconds int64

func init() {
	baseLmidSeconds = 1700000000
}

func (g *Generator) getGoImportForDataType(datatype string, nullable bool) string {
	switch datatype {
	case "bigint", "integer", "smallint", "boolean", "character varying", "text", "numeric":
		if nullable {
			return "database/sql"
		}
		return ""
	case "uuid":
		if nullable {
			return "database/sql"
		}
		return ""
	case "jsonb", "json":
		return "github.com/techrail/ground/typs/jsonObject"
	case "timestamp without time zone", "timestamp", "timestampz", "timestamp with time zone":
		if nullable {
			return "database/sql"
		}
		return "time"
	default:
		return ""
	}
}

func (g *Generator) removeTrailingNewlines(input string) string {
	// Split by new lines
	inputParts := strings.Split(input, "\n")
	multipleNewLinesInEnd := true
	for multipleNewLinesInEnd {
		if strings.TrimSpace(inputParts[len(inputParts)-1]) == "" &&
			strings.TrimSpace(inputParts[len(inputParts)-2 : len(inputParts)-1][0]) == "" {
			// If last two lines are empty strings (effectively), then remove the last one and check again
			multipleNewLinesInEnd = true
			inputParts = inputParts[:len(inputParts)-2]
		} else {
			multipleNewLinesInEnd = false
		}
	}
	return strings.Join(inputParts, "\n")
}

// Function to get the Go name for a given PostgreSQL table or column name
func getGoName(name string) string {
	nameParts := strings.Split(name, ".")
	if len(nameParts) > 1 {
		name = nameParts[1]
	}
	caser := cases.Title(language.English)
	retVal := strings.ReplaceAll(caser.String(strings.ReplaceAll(name, "_", " ")), " ", "")
	return retVal
}

func isColumnInList(columnName string, list []DbColumn) bool {
	for _, col := range list {
		if col.Name == columnName {
			return true
		}
	}
	return false
}

func newUniqueLmid() string {
	repeatZeros := func(times int) string {
		r := ""
		for i := 0; i < times; i++ {
			r += "0"
		}
		return r
	}
	baseLmidSeconds += 1
	lmid := integer.Base10ToBase36(baseLmidSeconds - 1600000000)
	prefix := ""
	if len(lmid) < 6 {
		prefix = repeatZeros(6 - len(lmid))
	}
	return prefix + lmid
}

func lowerFirstChar(input string) string {
	return strings.ToLower(input[:1]) + input[1:]
}

func upperFirstChar(input string) string {
	return strings.ToUpper(input[:1]) + input[1:]
}

func breakInPiecesOfStringSlicesIn3()
