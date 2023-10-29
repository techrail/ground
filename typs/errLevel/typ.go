package errLevel

// This type is compatible with the Bark error levels, except an addition of the "unknown" level

type Typ int

const (
	Panic   Typ = 10
	Alert   Typ = 9
	Error   Typ = 8
	Warning Typ = 4
	Notice  Typ = 3
	Info    Typ = 0
	Debug   Typ = -4
	Unknown Typ = -10
)

func (el Typ) String() string {
	switch el {
	case Panic:
		return "Panic"
	case Alert:
		return "Alert"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Notice:
		return "Notice"
	case Debug:
		return "Debug"
	case Unknown:
		return "Unknown"
	case Info:
		fallthrough
	default:
		return "Info"
	}
}

func (el Typ) ShortStr() string {
	return el.String()[0:1]
}

func (el Typ) Int8() int8 {
	return int8(el)
}
func (el Typ) Int() int {
	return int(el)
}

func FromShortString(shortString string) Typ {
	switch shortString {
	case "P":
		return Panic
	case "A":
		return Alert
	case "E":
		return Error
	case "W":
		return Warning
	case "N":
		return Notice
	case "D":
		return Debug
	case "I":
		fallthrough
	default:
		return Info
	}
}
