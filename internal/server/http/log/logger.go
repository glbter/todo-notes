package log

import "log"

type Logger interface {
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type MyLogger struct {}

func (MyLogger) Info(v ...interface{}) {
	log.Printf("info: %v", v)
}

func (MyLogger) Warn(v ...interface{}) {
	log.Printf("warn: %v", v)
}

func (MyLogger) Error(v ...interface{}) {
	log.Printf("error: %v", v)
}