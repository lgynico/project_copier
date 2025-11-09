package cherryLogger

import (
	"log"
)

func Warnf(format string, v ...any) {
	log.Printf(format, v...)
}

func Warn(args ...any) {
	log.Print(args...)
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

func Fatalf(format string, args ...any) {
	log.Fatalf(format, args...)
}

func Info(args ...any) {
	log.Print(args...)
}
