package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	maxLogFiles = 5
	logsDir     = "logs"
)

var (
	instance *Logger
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
	file   *os.File
}

func init() {
	instance = newLogger()
}

func newLogger() *Logger {
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatalf("Falha ao criar diretório de logs: %v", err)
	}

	timestamp := time.Now().UnixMilli()
	logFilename := filepath.Join(logsDir, fmt.Sprintf("debug_%d.log", timestamp))

	file, err := os.Create(logFilename)
	if err != nil {
		log.Fatalf("Falha ao criar arquivo de log: %v", err)
	}

	logger := log.New(file, "", log.LstdFlags)
	log.SetOutput(file)

	cleanOldLogs()

	return &Logger{
		level:  DEBUG,
		logger: logger,
		file:   file,
	}
}

func cleanOldLogs() {
	files, err := filepath.Glob(filepath.Join(logsDir, "debug_*.log"))
	if err != nil {
		return
	}

	if len(files) <= maxLogFiles {
		return
	}

	sort.Strings(files)

	for i := 0; i < len(files)-maxLogFiles; i++ {
		os.Remove(files[i])
	}
}

func Debugf(format string, args ...any) {
	if instance.level <= DEBUG {
		instance.logf("DEBUG", format, args...)
	}
}

func Infof(format string, args ...any) {
	if instance.level <= INFO {
		instance.logf("INFO", format, args...)
	}
}

func Warnf(format string, args ...any) {
	if instance.level <= WARN {
		instance.logf("WARN", format, args...)
	}
}

func Errorf(format string, args ...any) {
	if instance.level <= ERROR {
		instance.logf("ERROR", format, args...)
	}
}

func Fatalf(format string, args ...any) {
	instance.logf("FATAL", format, args...)
	os.Exit(1)
}

func (l *Logger) logf(level, format string, args ...any) {
	var msg string
	if len(args) == 0 {
		msg = format
	} else {
		msg = fmt.Sprintf(format, args...)
	}
	l.logger.Printf("[%s] %s", level, msg)
}
