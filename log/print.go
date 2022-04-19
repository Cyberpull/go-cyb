package log

import (
	"fmt"
	"log"
)

func Print(v ...interface{}) {
	log.Print(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Printfln(format string, v ...interface{}) {
	format = fmt.Sprintf(format, v...)
	log.Println(format)
}

func Println(v ...interface{}) {
	log.Println(v...)
}
