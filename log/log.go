package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// This is only basic skeleton. We are working on it, it's ongoing proccess :D
// This logger will be in a separate repository and will be used across all our services.

// TODO: Throw panic if logger was used but goroutine wasn't started. It will hange !!!

var file *os.File
var config Config

type logMsg struct {
	level string
	data  string
}

type Config struct {
	Logs chan logMsg

	ErrorWriters []io.StringWriter
	InfoWriters  []io.StringWriter
	DebugWriters []io.StringWriter
	TraceWriters []io.StringWriter
}

var Colors =  map[string]string {
  "red":     "\033[31m",
  "green":   "\033[32m",
  "yellow":  "\033[33m",
  "black":   "\033[30m",
  "blue":    "\033[34m",
  "magenta": "\033[35m",
  "cyan":    "\033[36m",
  "white":   "\033[37m",
}

// Standard colors.
const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Black   = "\033[30m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Reset   = "\033[0m"
)

const (
	ErrorPrefix = Red     + "ERROR"
	DebugPrefix = Magenta + "DEBUG"
	InfoPrefix  = Green   + "INFO"
	TracePrefix = Blue    + "TRACE"
)

// Quick implementation of io.StringWriter for stdout.
type StdoutWriter struct {}

func (stdout StdoutWriter) WriteString(log string) (int, error) {
	return fmt.Println(log)
}

// @TODO: Make it atomic so we can only call it once.
func Configure(userConf *Config) error {
	// Default writers for stdout and file (file is disabled for now).
	defaultWriters := []io.StringWriter{StdoutWriter{}}

	config = Config{
		// Setting writers for each log level.
		ErrorWriters: append(defaultWriters, userConf.ErrorWriters...),
		InfoWriters:  append(defaultWriters, userConf.InfoWriters...),
		DebugWriters: append(defaultWriters, userConf.DebugWriters...),
		TraceWriters: append(defaultWriters, userConf.TraceWriters...),
	}

	config.Logs = make(chan logMsg, 100)

	// file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// Start main goroutine that will receive all logs and will pass them to writers.
// This way only one goroutine will be responsible for writing things down and it
// won't slow down the rest of the system.
func Start() {
	// Reading logs one by one.
	for log := range config.Logs {
		switch log.level {
    case "ERROR":
			for _, writer := range config.ErrorWriters {
				writer.WriteString(log.data)
			}

    case "INFO":
			for _, writer := range config.InfoWriters {
				writer.WriteString(log.data)
			}

    case "DEBUG":
			for _, writer := range config.DebugWriters {
				writer.WriteString(log.data)
			}

		case "TRACE":
			for _, writer := range config.TraceWriters {
				writer.WriteString(log.data)
			}
    }
	}
}

func Error(format string, args ...any) {
  send(logMsg{"ERROR", formatLog(ErrorPrefix, format, args...)})
}

func Info(format string, args ...any) {
	send(logMsg{"INFO", formatLog(InfoPrefix, format, args...)})
}

func Debug(format string, args ...any) {
	send(logMsg{"DEBUG", formatLog(DebugPrefix, format, args...)})
}

func Trace(format string, args ...any) {
	send(logMsg{"TRACE", formatLog(TracePrefix, format, args...)})
}

// Send log to channel - If no receiving goroutine is set, display it on screen.
func send(log logMsg) {
  select {
  case config.Logs <- log:
  default:
    fmt.Println(log.data)
  }
}

// Format log and add default prefix to it.
func formatLog(prefix, format string, args ...any) string {
	log := fmt.Sprintf(format, args...)
	return (defaultPrefix(prefix) + log + Reset)
}

// Return default prefix with caller file name and line number. 
// 
// @TODO: should we also add caller function name?
func defaultPrefix(prefix string) string {
	_, file, line, _ := runtime.Caller(4)
	file = filepath.Base(file)

	// Remove file extension.
	file = strings.Trim(file, ".go")

	fileline := fmt.Sprintf("|%s:%d|", file, line)
	prefix    = fmt.Sprintf("%-11s%-18s", prefix, fileline)

	return prefix
}
