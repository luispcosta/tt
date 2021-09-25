package reporter

import (
	"strings"

	"github.com/luispcosta/go-tt/core"
)

const jsonFormat = "json"
const csvFormat = "csv"
const cliFormat = "cli"

// AllowedFormats creates a map with the allowed report formats and their implementations
func AllowedFormats() map[string]core.Reporter {
	allowedFormats := make(map[string]core.Reporter)
	allowedFormats[csvFormat] = NewEmptyReporter()
	allowedFormats[jsonFormat] = NewJsonReporter()
	allowedFormats[cliFormat] = NewCliReporter()
	return allowedFormats
}

// AllowedFormatsCollection returns the collection of allowed formats
func AllowedFormatsCollection() []string {
	return []string{jsonFormat, csvFormat, cliFormat}
}

// IsAllowedFormat returns true if the format is allowed
func IsAllowedFormat(format string) bool {
	return CreateReporter(format) != nil
}

// CreateReporter creates a new reporter
func CreateReporter(format string) core.Reporter {
	return AllowedFormats()[strings.ToLower(format)]
}
