package log

import (
	"github.com/fatih/color"
)

type ColorAttribute color.Attribute

const (
	FgGreen   ColorAttribute = ColorAttribute(color.FgGreen)
	FgRed     ColorAttribute = ColorAttribute(color.FgRed)
	FgCyan    ColorAttribute = ColorAttribute(color.FgCyan)
	FgMagenta ColorAttribute = ColorAttribute(color.FgMagenta)
	FgBlack   ColorAttribute = ColorAttribute(color.FgBlack)
	FgBlue    ColorAttribute = ColorAttribute(color.FgBlue)
)

func Color(a ColorAttribute, v ...interface{}) {
	color.Set(color.Attribute(a))
	Print(v...)
	color.Unset()
}

func Colorf(a ColorAttribute, format string, v ...interface{}) {
	color.Set(color.Attribute(a))
	Printf(format, v...)
	color.Unset()
}

func Colorfln(a ColorAttribute, format string, v ...interface{}) {
	color.Set(color.Attribute(a))
	Printfln(format, v...)
	color.Unset()
}

func Colorln(a ColorAttribute, v ...interface{}) {
	color.Set(color.Attribute(a))
	Println(v...)
	color.Unset()
}
