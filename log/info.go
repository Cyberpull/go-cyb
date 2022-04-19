package log

func Info(v ...interface{}) {
	Color(FgCyan, v...)
}

func Infof(format string, v ...interface{}) {
	Colorf(FgCyan, format, v...)
}

func Infofln(format string, v ...interface{}) {
	Colorfln(FgCyan, format, v...)
}

func Infoln(v ...interface{}) {
	Colorln(FgCyan, v...)
}
