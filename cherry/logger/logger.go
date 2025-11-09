package cherryLogger

import (
	"log"
)

func Warnf(format string, v ...any) {
	log.Printf(format, v...)
}

func Errorf(format string, v ...any) {
	log.Printf(format, v...)
}

func Debugf(format string, v ...any) {
	log.Printf(format, v...)
}

func Infof(format string, v ...any) {
	log.Printf(format, v...)
}

func Info(args ...any) {
	log.Print(args...)
}
