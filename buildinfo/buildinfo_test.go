package buildinfo

import (
	"errors"
	"testing"
	"time"
)

func TestTimestampWithoutBuildValue(t *testing.T) {
	if timestamp != "" {
		t.Skip("a linker-provided timestamp is set")
	}
	if got := RawTimestamp(); got != "" {
		t.Fatalf("RawTimestamp() = %q, want an empty value", got)
	}
	got, err := Timestamp()
	if !got.IsZero() || !errors.Is(err, ErrNoTimestamp) {
		t.Fatalf("Timestamp() = %v, %v; want zero time and ErrNoTimestamp", got, err)
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		want  string
		valid bool
	}{
		{name: "ISO8601", raw: "2026-08-09T12:34:56+0200", want: "2026-08-09T10:34:56Z", valid: true},
		{name: "RFC3339", raw: "2026-08-09T12:34:56+02:00", want: "2026-08-09T10:34:56Z", valid: true},
		{name: "invalid", raw: "not-a-timestamp"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			previous := timestamp
			t.Cleanup(func() { timestamp = previous })
			timestamp = test.raw

			got, err := parseTimestamp()
			if !test.valid {
				if err == nil {
					t.Fatal("parseTimestamp() accepted an invalid timestamp")
				}
				return
			}
			if err != nil || !got.Equal(mustParseTime(test.want)) {
				t.Fatalf("parseTimestamp() = %v, %v; want %v", got, err, test.want)
			}
		})
	}
}

func TestTimestampReturnsParseError(t *testing.T) {
	previousTimestamp := timestamp
	previousParsed := timestampParsed
	previousError := timestampParseError
	t.Cleanup(func() {
		timestamp = previousTimestamp
		timestampParsed = previousParsed
		timestampParseError = previousError
	})

	timestamp = "invalid"
	timestampParsed = time.Time{}
	_, timestampParseError = parseTimestamp()
	got, err := Timestamp()
	if !got.IsZero() || err == nil || errors.Is(err, ErrNoTimestamp) {
		t.Fatalf("Timestamp() = %v, %v; want a parse error", got, err)
	}
}

func TestTimestampFromLinkerValue(t *testing.T) {
	if timestamp == "" {
		t.Skip("run with -ldflags '-X github.com/fanonwue/goutils/buildinfo.timestamp=...' to test linker injection")
	}
	want, err := parseTimestamp()
	if err != nil {
		t.Fatalf("linker-provided timestamp %q is not parseable: %v", timestamp, err)
	}
	if got := RawTimestamp(); got != timestamp {
		t.Fatalf("RawTimestamp() = %q, want %q", got, timestamp)
	}
	got, err := Timestamp()
	if err != nil || !got.Equal(want) {
		t.Fatalf("Timestamp() = %v, %v; want %v", got, err, want)
	}
}

func mustParseTime(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return parsed
}
