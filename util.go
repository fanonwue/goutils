package goutils

import (
	"errors"
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
