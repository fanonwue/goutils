package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type LogLevel uint8

func (ll LogLevel) Name() string {
	switch ll {
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	}
	panic("unreachable")
}

const (
	LevelPanic LogLevel = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

func (ll LogLevel) NameFormatted() string {
	return levelNamesFormatted[ll]
}

func (ll LogLevel) Logger() *log.Logger {
	if ll <= stdErrThresholdLevel {
		return errorLogger
	}
	return defaultLogger
}

const DefaultLevel = LevelInfo
const DefaultCalldepth = 3
const stdErrThresholdLevel = LevelError
const loggerFlags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds

var logLevel = DefaultLevel

var defaultLogger = log.New(os.Stdout, "", loggerFlags)
var errorLogger = log.New(os.Stderr, "", loggerFlags)

var logLevels []LogLevel
var levelNamesFormatted []string

func LogLevels() []LogLevel {
	return logLevels
}

func init() {
	logLevels = []LogLevel{LevelPanic, LevelFatal, LevelError, LevelWarn, LevelInfo, LevelDebug, LevelTrace}
	slices.Sort(logLevels)

	levelNamesFormatted = make([]string, 0, len(logLevels))
	for _, level := range logLevels {
		name := level.Name()
		paddingLength := 6 - len(name)
		levelNamesFormatted = append(levelNamesFormatted, fmt.Sprintf("[%s]%-*s", name, paddingLength, " "))
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(loggerFlags)
}

func callerInfo(calldepth int) string {
	_, file, no, _ := runtime.Caller(calldepth)

	fileParts := strings.Split(file, "/")
	caller := fileParts[len(fileParts)-2:]
	return strings.Join(caller, "/") + ":" + strconv.Itoa(no)
}

func SetLogLevelFromEnvironment(envVar string) error {
	value := os.Getenv(envVar)
	if value == "" {
		return nil
	}
	return SetLogLevelByName(value)
}

func logLevelByName(levelName string) (LogLevel, error) {
	// Setting the log level should not be performance-critical, so we don't bother with a map, set or caching
	normalizedName := strings.ToUpper(levelName)
	for _, level := range logLevels {
		if level.Name() == normalizedName {
			return level, nil
		}
	}
	return 0, fmt.Errorf("unknown log level: %s", levelName)
}

func SetLogLevelByName(levelName string) error {
	level, err := logLevelByName(levelName)
	if err != nil {
		return err
	}
	return SetLogLevel(level)
}

func SetLogLevel(level LogLevel) error {
	// Setting the log level should not be performance-critical, so we don't bother with a map, set or caching
	if !slices.Contains(logLevels, level) {
		return fmt.Errorf("unknown log level: %d", level)
	}
	logLevel = level
	return nil
}

func Logf(level LogLevel, calldepth int, msg string, args ...any) {
	if logLevel < level {
		return
	}
	if calldepth == 0 {
		calldepth = DefaultCalldepth
	}

	newArgs := []any{level.NameFormatted()}

	err := level.Logger().Output(calldepth, fmt.Sprintf("\t%s"+msg, append(newArgs, args...)...))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: Could not write log message: %v", err)
	}
}

func Info(msg string) {
	Logf(LevelInfo, DefaultCalldepth, msg)
}

func Infof(msg string, args ...any) {
	Logf(LevelInfo, DefaultCalldepth, msg, args...)
}

func Warn(msg string) {
	Logf(LevelWarn, DefaultCalldepth, msg)
}

func Warnf(msg string, args ...any) {
	Logf(LevelWarn, DefaultCalldepth, msg, args...)
}

func Error(msg string) {
	Logf(LevelError, DefaultCalldepth, msg)
}

func Errorf(msg string, args ...any) {
	Logf(LevelError, DefaultCalldepth, msg, args...)
}

func Debug(msg string) {
	Logf(LevelDebug, DefaultCalldepth, msg)
}

func Debugf(msg string, args ...any) {
	Logf(LevelDebug, DefaultCalldepth, msg, args...)
}

func Panic(msg string) {
	Logf(LevelPanic, DefaultCalldepth, msg)
}

func Panicf(msg string, args ...any) {
	Logf(LevelPanic, DefaultCalldepth, msg, args...)
}

func Fatal(msg string) {
	Logf(LevelFatal, DefaultCalldepth, msg)
}

func Fatalf(msg string, args ...any) {
	Logf(LevelFatal, DefaultCalldepth, msg, args...)
}
