package logger

import "log"

func Info(message string, args ...interface{}) {
	log.Printf(message, args...)
}
