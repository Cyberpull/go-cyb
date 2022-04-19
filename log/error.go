package log

func Error(v ...interface{}) {
	Color(FgRed, v...)
}

func Errorf(format string, v ...interface{}) {
	Colorf(FgRed, format, v...)
}

func Errorfln(format string, v ...interface{}) {
	Colorfln(FgRed, format, v...)
}

func Errorln(v ...interface{}) {
	Colorln(FgRed, v...)
}
