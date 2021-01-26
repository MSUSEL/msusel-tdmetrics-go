package parser

import (
	"fmt"
	"strings"
)

// RecordProcess is a process handler which writes to the given buffer.
func RecordProcess(buf *strings.Builder) ProcessHandler {
	return func(filename string, data interface{}) {
		if buf.Len() > 0 {
			fmt.Fprintln(buf)
		}
		fmt.Fprint(buf, data)
	}
}

// RecordError is an error handler which writes to the given buffer.
func RecordError(buf *strings.Builder) OnErrorHander {
	return func(err error) {
		if buf.Len() > 0 {
			fmt.Fprintln(buf)
		}
		fmt.Fprintf(buf, `Error: %s`, err.Error())
	}
}

// DefaultOnError handles an error which occurred by
// printing the error to the console.
// This will be used if the OnError is not set or nil.
func DefaultOnError(filename string, err error) {
	fmt.Print("Error in ", filename, ": ", err, "\n")
}

// DefaultUpdateProgress handles the progress being updated by
// printing the progress to the console.
// This will be used if the UpdateProgress is not set or nil.
func DefaultUpdateProgress(finished, total int) {
	progress := 1.0
	if total > 0 {
		progress = float64(finished) / float64(total)
	}
	fmt.Printf("%6.2f%% (%d/%d)\n", progress*100.0, finished, total)
}

// MuteProgress will mute the progress output.
func MuteProgress(finished, total int) {}

// IsTestFile returns true if the given filename is a test file.
func IsTestFile(filename string) bool {
	return strings.HasSuffix(filename, `_test.go`)
}

// IsGitFile returns true if the given filename is in a git path.
func IsGitFile(filename string) bool {
	return strings.Contains(filename, `/.git/`)
}

// IsVendorFile returns true if the given filename is in a vendor path.
func IsVendorFile(filename string) bool {
	return strings.Contains(filename, `/vendor/`)
}

// IsNotGoFile returns true if the given filename is NOT a *.go file.
func IsNotGoFile(filename string) bool {
	return !strings.HasSuffix(filename, `.go`)
}
