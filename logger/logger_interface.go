package logger

type Logger interface {
	Panic(string)
	Alert(string, bool)
	Error(string)
	Warn(string)
	Notice(string)
	Info(string)
	Debug(string)
	Println(string)
	Default(string)
}
