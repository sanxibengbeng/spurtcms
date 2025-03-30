package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

// Log levels
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = map[LogLevel]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

// LogFormat defines the output format of log messages
type LogFormat int

const (
	TextFormat LogFormat = iota
	JSONFormat
)

// Config holds the logger configuration
type Config struct {
	Level      LogLevel
	Format     LogFormat
	OutputPath string
	Stdout     bool
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      InfoLevel,
		Format:     TextFormat,
		OutputPath: "logs/spurtcms.log",
		Stdout:     true,
	}
}

// ImprovedLogger is a structured logger for SpurtCMS
type ImprovedLogger struct {
	config Config
	writer io.Writer
	file   *os.File
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	File      string         `json:"file,omitempty"`
	Line      int            `json:"line,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
}

// NewLogger creates a new ImprovedLogger with the given configuration
func NewLogger(config Config) (*ImprovedLogger, error) {
	// Create logs directory if it doesn't exist
	if config.OutputPath != "" {
		dir := filepath.Dir(config.OutputPath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	var writers []io.Writer
	var file *os.File

	// Add file writer if output path is specified
	if config.OutputPath != "" {
		var err error
		file, err = os.OpenFile(config.OutputPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		writers = append(writers, file)
	}

	// Add stdout writer if enabled
	if config.Stdout {
		writers = append(writers, os.Stdout)
	}

	// Create multi-writer if we have multiple outputs
	var writer io.Writer
	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		// Default to stdout if no writers specified
		writer = os.Stdout
	}

	return &ImprovedLogger{
		config: config,
		writer: writer,
		file:   file,
	}, nil
}

// Close closes the logger's file handle if it exists
func (l *ImprovedLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// log logs a message at the specified level with optional fields
func (l *ImprovedLogger) log(level LogLevel, msg string, fields map[string]any) {
	if level < l.config.Level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	// Shorten file path for readability
	file = filepath.Base(file)

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     levelNames[level],
		Message:   msg,
		File:      file,
		Line:      line,
		Fields:    fields,
	}

	switch l.config.Format {
	case JSONFormat:
		jsonData, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Error marshaling log entry: %v", err)
			return
		}
		fmt.Fprintln(l.writer, string(jsonData))
	case TextFormat:
		var fieldsStr string
		if len(fields) > 0 {
			parts := make([]string, 0, len(fields))
			for k, v := range fields {
				parts = append(parts, fmt.Sprintf("%s=%v", k, v))
			}
			fieldsStr = " " + strings.Join(parts, " ")
		}
		fmt.Fprintf(l.writer, "%s [%s] %s:%d %s%s\n", entry.Timestamp, entry.Level, entry.File, entry.Line, entry.Message, fieldsStr)
	}

	// Flush immediately for fatal logs
	if level == FatalLevel && l.file != nil {
		l.file.Sync()
	}
}

// Debug logs a message at debug level
func (l *ImprovedLogger) Debug(msg string, fields ...map[string]any) {
	var f map[string]any
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DebugLevel, msg, f)
}

// Info logs a message at info level
func (l *ImprovedLogger) Info(msg string, fields ...map[string]any) {
	var f map[string]any
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(InfoLevel, msg, f)
}

// Warn logs a message at warn level
func (l *ImprovedLogger) Warn(msg string, fields ...map[string]any) {
	var f map[string]any
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WarnLevel, msg, f)
}

// Error logs a message at error level
func (l *ImprovedLogger) Error(msg string, fields ...map[string]any) {
	var f map[string]any
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(ErrorLevel, msg, f)
}

// Fatal logs a message at fatal level and exits the application
func (l *ImprovedLogger) Fatal(msg string, fields ...map[string]any) {
	var f map[string]any
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FatalLevel, msg, f)
	os.Exit(1)
}

// WithField returns a new map with a single field added
func WithField(key string, value any) map[string]any {
	return map[string]any{key: value}
}

// WithFields returns a map with multiple fields
func WithFields(fields map[string]any) map[string]any {
	return fields
}

// WithError returns a map with an error field
func WithError(err error) map[string]any {
	return map[string]any{"error": err.Error()}
}

// Global logger instance
var globalLogger *ImprovedLogger

// init initializes the global logger with default configuration
func init() {
	// Load configuration from environment variables
	config := DefaultConfig()

	// Override with environment variables if set
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		switch strings.ToUpper(level) {
		case "DEBUG":
			config.Level = DebugLevel
		case "INFO":
			config.Level = InfoLevel
		case "WARN":
			config.Level = WarnLevel
		case "ERROR":
			config.Level = ErrorLevel
		case "FATAL":
			config.Level = FatalLevel
		}
	}

	if format := os.Getenv("LOG_FORMAT"); format != "" {
		switch strings.ToUpper(format) {
		case "JSON":
			config.Format = JSONFormat
		case "TEXT":
			config.Format = TextFormat
		}
	}

	if path := os.Getenv("LOG_PATH"); path != "" {
		config.OutputPath = path
	}

	if stdout := os.Getenv("LOG_STDOUT"); stdout != "" {
		config.Stdout = stdout == "true"
	}

	var err error
	globalLogger, err = NewLogger(config)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
}

// GetLogger returns the global logger instance
func GetLogger() *ImprovedLogger {
	return globalLogger
}

// Debug logs a message at debug level using the global logger
func Debug(msg string, fields ...map[string]any) {
	globalLogger.Debug(msg, fields...)
}

// Info logs a message at info level using the global logger
func Info(msg string, fields ...map[string]any) {
	globalLogger.Info(msg, fields...)
}

// Warn logs a message at warn level using the global logger
func Warn(msg string, fields ...map[string]any) {
	globalLogger.Warn(msg, fields...)
}

// Error logs a message at error level using the global logger
func Error(msg string, fields ...map[string]any) {
	globalLogger.Error(msg, fields...)
}

// Fatal logs a message at fatal level using the global logger and exits
func Fatal(msg string, fields ...map[string]any) {
	globalLogger.Fatal(msg, fields...)
}

// Compatibility functions for existing code

// ErrorLog returns a standard logger for backward compatibility
func ErrorLog() *log.Logger {
	return log.New(globalLogger.writer, "ERROR ", log.LstdFlags|log.Lshortfile)
}

// WarnLog returns a standard logger for backward compatibility
func WarnLog() *log.Logger {
	return log.New(globalLogger.writer, "WARN ", log.LstdFlags|log.Lshortfile)
}

// InfoLog returns a standard logger for backward compatibility
func InfoLog() *log.Logger {
	return log.New(globalLogger.writer, "INFO ", log.LstdFlags|log.Lshortfile)
}
