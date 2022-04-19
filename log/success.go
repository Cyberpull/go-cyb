package log

func Success(v ...interface{}) {
	Color(FgGreen, v...)
}

func Successf(format string, v ...interface{}) {
	Colorf(FgGreen, format, v...)
}

func Successfln(format string, v ...interface{}) {
	Colorfln(FgGreen, format, v...)
}

func Successln(v ...interface{}) {
	Colorln(FgGreen, v...)
}
