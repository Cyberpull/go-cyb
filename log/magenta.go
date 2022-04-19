package log

func Magenta(v ...interface{}) {
	Color(FgMagenta, v...)
}

func Magentaf(format string, v ...interface{}) {
	Colorf(FgMagenta, format, v...)
}

func Magentafln(format string, v ...interface{}) {
	Colorfln(FgMagenta, format, v...)
}

func Magentaln(v ...interface{}) {
	Colorln(FgMagenta, v...)
}
