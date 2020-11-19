package logger

import (
	"fmt"
	"log"
	"os"
)

func NewLogger(name string) *Logger {
	return &Logger{
		Info: log.New(os.Stdout, fmt.Sprintf("%s: %s ", "INFO", name), log.Ldate|log.Lmicroseconds|log.Lshortfile),

		Error: log.New(os.Stdout, fmt.Sprintf("%s: %s ", "ERROR", name), log.Ldate|log.Lmicroseconds|log.Lshortfile),
	}
}

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}
