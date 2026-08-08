package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogLevelNames(t *testing.T) {
	want := []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"}
	for level, name := range want {
		got := LogLevel(level).Name()
		if got != name || LogLevel(level).String() != name {
			t.Errorf("level %d name = %q, want %q", level, got, name)
		}
	}
	if got := LevelInfo.NameFormatted(); got != "[INFO]  " {
		t.Errorf("formatted info level = %q, want %q", got, "[INFO]  ")
	}
}

func TestSetLogLevel(t *testing.T) {
	previous := logLevel
	t.Cleanup(func() { logLevel = previous })

	if err := SetLogLevelByName("debug"); err != nil || logLevel != LevelDebug {
		t.Fatalf("SetLogLevelByName(debug) = level %v, error %v", logLevel, err)
	}
	if err := SetLogLevelByName("not-a-level"); err == nil {
		t.Fatal("SetLogLevelByName accepted an unknown level")
	}
	if err := SetLogLevel(LogLevel(99)); err == nil {
		t.Fatal("SetLogLevel accepted an unknown level")
	}

	envName := "GOUTILS_TEST_LOG_LEVEL"
	t.Setenv(envName, "trace")
	if err := SetLogLevelFromEnvironment(envName); err != nil || logLevel != LevelTrace {
		t.Fatalf("SetLogLevelFromEnvironment(trace) = level %v, error %v", logLevel, err)
	}
	t.Setenv(envName, "")
	if err := SetLogLevelFromEnvironment(envName); err != nil || logLevel != LevelTrace {
		t.Fatalf("empty environment value changed level: %v, error %v", logLevel, err)
	}
	if levels := LogLevels(); len(levels) != 7 {
		t.Fatalf("LogLevels() returned %d levels, want 7", len(levels))
	}
}

func TestLogfFilteringAndFormatting(t *testing.T) {
	previousLevel := logLevel
	previousDefaultOutput := defaultLogger.Writer()
	previousErrorOutput := errorLogger.Writer()
	previousDefaultFlags := defaultLogger.Flags()
	previousErrorFlags := errorLogger.Flags()
	t.Cleanup(func() {
		logLevel = previousLevel
		defaultLogger.SetOutput(previousDefaultOutput)
		errorLogger.SetOutput(previousErrorOutput)
		defaultLogger.SetFlags(previousDefaultFlags)
		errorLogger.SetFlags(previousErrorFlags)
	})

	var standard, errors bytes.Buffer
	defaultLogger.SetOutput(&standard)
	errorLogger.SetOutput(&errors)
	defaultLogger.SetFlags(0)
	errorLogger.SetFlags(0)
	logLevel = LevelInfo

	Debug("hidden")
	Infof("hello %s", "world")
	Error("broken")
	if strings.Contains(standard.String(), "hidden") || !strings.Contains(standard.String(), "\t[INFO]  hello world") {
		t.Errorf("unexpected standard log output: %q", standard.String())
	}
	if !strings.Contains(errors.String(), "[ERROR] broken") {
		t.Errorf("unexpected error log output: %q", errors.String())
	}
}

func TestSetLogLevelFromEnvironmentMissing(t *testing.T) {
	name := "GOUTILS_MISSING_LOG_LEVEL"
	t.Setenv(name, "")
	if err := SetLogLevelFromEnvironment(name); err != nil {
		t.Fatalf("missing environment variable returned error: %v", err)
	}
}
