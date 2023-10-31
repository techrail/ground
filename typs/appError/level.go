package appError

// This type is compatible with the Bark error levels, except an addition of the "unknown" level

type Level int

const (
	Panic   Level = 10
	Alert   Level = 9
	Error   Level = 8
	Warning Level = 4
	Notice  Level = 3
	Info    Level = 0
	Debug   Level = -4
	Unknown Level = -10
)

func (el Level) String() string {
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

func (el Level) ShortStr() string {
	return el.String()[0:1]
}

func (el Level) Int8() int8 {
	return int8(el)
}
func (el Level) Int() int {
	return int(el)
}

func FromShortString(shortString string) Level {
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
