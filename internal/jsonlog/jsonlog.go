package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Creer une niveau de log
// une type ? parce'que better ?

type Level int8

const (
	LevelInfo Level = iota // Has a value of 0
	LevelError
	LevelFatal
	LevelOff // Has a value of 3
)

// Map each enum level, a function

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Creer nouveau de type de logger.
// Consists of
// Level
// The writer <--
// Mutex <- to prevent race condition

type Logger struct {
	out      io.Writer
	minLevel Level
	mutex    sync.Mutex
}

func New(writer io.Writer, level Level) *Logger {
	return &Logger{
		minLevel: level,
		out:      writer,
	}
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.out.Write(append(line, '\n'))
}

func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}

//
