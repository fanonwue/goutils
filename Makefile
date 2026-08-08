test: test-default test-inject

test-default:
	go test ./...

# Testing the buildinfo package requires setting a timestamp at build time
test-inject:
	go test ./buildinfo -ldflags='-X github.com/fanonwue/goutils/buildinfo.timestamp=2026-08-09T12:34:56+0200'
