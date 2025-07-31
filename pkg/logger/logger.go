package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = [...]string{"DEBUG", "INFO ", "WARN ", "ERROR", "FATAL"}

type Logger struct {
	mu    sync.Mutex
	level Level
	out   *log.Logger
}

var Log = NewLogger(INFO)

func NewLogger(lvl Level) *Logger {
	return &Logger{
		level: lvl,
		out:   log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) log(lvl Level, msg string, v ...interface{}) {
	if lvl < l.level {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.mu.Lock()
	defer l.mu.Unlock()
	line := fmt.Sprintf("[%s] [%s] %s", timestamp, levelNames[lvl], fmt.Sprintf(msg, v...))
	l.out.Println(line)
	if lvl == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) SetLevel(lvl Level) {
	l.level = lvl
}

func (l *Logger) Debug(msg string, v ...interface{}) { l.log(DEBUG, msg, v...) }
func (l *Logger) Info(msg string, v ...interface{})  { l.log(INFO, msg, v...) }
func (l *Logger) Warn(msg string, v ...interface{})  { l.log(WARN, msg, v...) }
func (l *Logger) Error(msg string, v ...interface{}) { l.log(ERROR, msg, v...) }
func (l *Logger) Fatal(msg string, v ...interface{}) { l.log(FATAL, msg, v...) }