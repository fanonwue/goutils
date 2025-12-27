package goutils

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fanonwue/goutils/dsext"
)

func PanicHandler(handler func(err any)) {
	if err := recover(); err != nil {
		handler(err)
	}
}

type EnvVarHelper string

const DefaultEnvVarPrefix = "APP_"

func NewEnvVarHelper(prefix string) EnvVarHelper {
	if prefix == "" {
		prefix = DefaultEnvVarPrefix
	}

	if prefix[len(prefix)-1] != '_' {

	}

	return EnvVarHelper(prefix)
}

func (evh EnvVarHelper) PrefixVar(s string) string {
	return string(evh) + s
}

func (evh EnvVarHelper) Bool(key string, defaultValue bool) (bool, error) {
	raw := os.Getenv(evh.PrefixVar(key))
	if raw == "" {
		return defaultValue, nil
	}
	return strconv.ParseBool(raw)
}

func (evh EnvVarHelper) Int(key string, defaultValue int64) (int64, error) {
	raw := os.Getenv(evh.PrefixVar(key))
	if raw == "" {
		return defaultValue, nil
	}
	return strconv.ParseInt(raw, 10, 64)
}

var truthyValues = dsext.NewSetSlice([]string{"1", "true", "yes", "on", "enable"})

func TruthyValues() dsext.Set[string] {
	return truthyValues
}

func IsTruthy(s string) bool {
	return truthyValues.Contains(strings.ToLower(s))
}

func EpochStringToTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("empty epoch string")
	}

	timeAttr, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(timeAttr, 0), nil
}

func TruncateStringWholeWords(s string, maxLength uint) string {
	lastSpaceIx := -1
	length := uint(0)
	for i, r := range s {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		length++
		if length >= maxLength {
			if lastSpaceIx != -1 {
				return s[:lastSpaceIx] + "..."
			}
			// If here, s is longer than maxLength but has no spaces
		}
	}
	// If here, s is shorter than maxLength
	return s
}

// WithFile opens the file specified by the path using [os.Open] and calls f with it. The file is closed after f returns.
func WithFile[T any](path string, f func(file *os.File) (T, error)) (T, error) {
	return withFileOs(path, f, os.Open)
}

// WithFileRoot opens the file specified by the path within the given root and calls f with it. The file is closed after f returns.
func WithFileRoot[T any](path string, root *os.Root, f func(file *os.File) (T, error)) (T, error) {
	if root == nil {
		var result T
		return result, errors.New("root is nil")
	}
	return withFileOs(path, f, root.Open)
}

func withFileOs[T any](path string, f func(file *os.File) (T, error), openFunc func(string) (*os.File, error)) (T, error) {
	file, err := openFunc(path)
	if err != nil {
		var result T
		return result, err
	}
	defer func() { _ = file.Close() }()
	return f(file)
}

// WithFileFS opens the file in filesystem fs specified by the path and calls f with it. The file is closed after f returns.
func WithFileFS[T any](path string, fs fs.FS, f func(file fs.File) (T, error)) (T, error) {
	file, err := fs.Open(path)
	if err != nil {
		var result T
		return result, err
	}
	defer func() { _ = file.Close() }()
	return f(file)
}

// SplitAny splits the string s around each instance of one of the Unicode code points in seps.
func SplitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}
