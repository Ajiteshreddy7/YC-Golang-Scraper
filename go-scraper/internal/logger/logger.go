package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel = INFO
	logger       = log.New(os.Stdout, "", 0)
)

// SetLevel configures the minimum log level to display
func SetLevel(level Level) {
	currentLevel = level
}

// SetOutput configures where logs are written
func SetOutput(file *os.File) {
	logger.SetOutput(file)
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// ParseLevel converts a string like "DEBUG"/"INFO"/"WARN"/"ERROR" to a Level
func ParseLevel(s string) (Level, bool) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return DEBUG, true
	case "INFO", "":
		return INFO, true
	case "WARN", "WARNING":
		return WARN, true
	case "ERROR":
		return ERROR, true
	default:
		return INFO, false
	}
}

// InitFromEnv sets log level from LOG_LEVEL environment variable (default INFO)
func InitFromEnv() {
	lvlStr := os.Getenv("LOG_LEVEL")
	if lvl, ok := ParseLevel(lvlStr); ok {
		SetLevel(lvl)
	} else {
		SetLevel(INFO)
	}
}

func Debug(format string, v ...interface{}) {
	if currentLevel <= DEBUG {
		logger.Printf("[DEBUG] %s - %s", timestamp(), fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	if currentLevel <= INFO {
		logger.Printf("[INFO]  %s - %s", timestamp(), fmt.Sprintf(format, v...))
	}
}

func Warn(format string, v ...interface{}) {
	if currentLevel <= WARN {
		logger.Printf("[WARN]  %s - %s", timestamp(), fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	if currentLevel <= ERROR {
		logger.Printf("[ERROR] %s - %s", timestamp(), fmt.Sprintf(format, v...))
	}
}

func Fatal(format string, v ...interface{}) {
	logger.Printf("[FATAL] %s - %s", timestamp(), fmt.Sprintf(format, v...))
	os.Exit(1)
}
