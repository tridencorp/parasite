package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// This is only basic skeleton. We are working on it, it's ongoing proccess :D
// This logger will be in a separate repository and will be used across all our services.

var file *os.File
var config Config

type Config struct {
	ErrorWriters []io.StringWriter
	InfoWriters  []io.StringWriter
	DebugWriters []io.StringWriter
	TraceWriters []io.StringWriter
}

// Standard colors.
const (
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	black   = "\033[30m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	reset   = "\033[0m"
)

// Quick implementation of io.StringWriter for stdout.
type StdoutWriter struct {}

func (stdout StdoutWriter) WriteString(log string) (int, error) {
	return fmt.Println(log)
}

// @TODO: Make it atomic so we can only call it once.
func Configure(conf *Config) error {
	// Default writers for stdout and file(disabled for now).
	defaultWriters := []io.StringWriter{StdoutWriter{}}

	// We are setting package config.
	config = Config{
		// Setting writers for each log level.
		ErrorWriters: append(defaultWriters, conf.ErrorWriters...),
		InfoWriters:  append(defaultWriters, conf.InfoWriters...),
		DebugWriters: append(defaultWriters, conf.DebugWriters...),
		TraceWriters: append(defaultWriters, conf.TraceWriters...),
	}

	// file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func Error(format string, args ...any) {
	log := formatLog("ERROR ", format, args...)

	for _, writer := range config.ErrorWriters {
		writer.WriteString(red + log + reset)
	}
}

func Info(format string, args ...any) {
	log := formatLog("INFO ", format, args...)

	for _, writer := range config.InfoWriters {
		writer.WriteString(magenta + log + reset)
	}
}

func Debug(format string, args ...any) {
	log := formatLog("DEBUG ", format, args...)

	for _, writer := range config.DebugWriters {
		writer.WriteString(green + log + reset)
	}
}

func Trace(format string, args ...any) {
	log := formatLog("TRACE ", format, args...)

	for _, writer := range config.TraceWriters {
		writer.WriteString(blue + log + reset)
	}
}

// Format log and add default prefix to it.
func formatLog(prefix, format string, args ...any) string {
	log := fmt.Sprintf(format, args...)
	return (prefix + defaultPrefix() + log)
}

// Return default prefix with caller file name and line number. 
// 
// @TODO: should we also add caller function name?
func defaultPrefix() string {
	_, file, line, _ := runtime.Caller(3)

	file = filepath.Base(file)
	prefix := fmt.Sprintf("(%s:%d) ", file, line)

	return prefix
}
