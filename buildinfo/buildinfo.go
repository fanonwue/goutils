package buildinfo

import (
	"errors"
	"time"
)

const TimestampFormat = "2006-01-02T15:04:05Z0700"

var timestamp = ""
var timestampParsed = time.Time{}
var timestampParseError error

var ErrNoTimestamp = errors.New("no build timestamp set")

// RawTimestamp returns the build timestamp as a string. The build timestamp is set by the build script using a linker flag.
//
// To set the timestamp, run `go build -ldflags "-X github.com/fanonwue/goutils/buildinfo.timestamp=TIMESTAMP"`, where TIMESTAMP
// is the desired timestamp in RFC3339 or ISO8601 format.
func RawTimestamp() string {
	return timestamp
}

// Timestamp returns the build timestamp as a time.Time object.
// If an error occurs while parsing the timestamp, the error is returned.
// If no timestamp is set, ErrNoTimestamp is returned.
func Timestamp() (time.Time, error) {
	if timestampParseError != nil {
		return time.Time{}, timestampParseError
	}
	if timestampParsed.IsZero() {
		return time.Time{}, ErrNoTimestamp
	}
	return timestampParsed, nil
}

func parseTimestamp() (time.Time, error) {
	// Try the ISO8601 format first
	ts, err := time.Parse(TimestampFormat, timestamp)
	if err == nil {
		return ts, nil
	}

	// Fallback to RFC3339
	return time.Parse(time.RFC3339, timestamp)
}

func init() {
	ts, err := parseTimestamp()
	if err == nil {
		timestampParsed = ts
	}
	timestampParseError = err

}
