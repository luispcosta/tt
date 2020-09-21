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
	allowedFormats[jsonFormat] = NewEmptyReporter()
	allowedFormats[cliFormat] = NewCliReporter()
	return allowedFormats
}

// IsAllowedFormat returns true if the format is allowed
func IsAllowedFormat(format string) bool {
	return Create(format) != nil
}

// Create creates a new reporter
func Create(format string) core.Reporter {
	return AllowedFormats()[strings.ToLower(format)]
}
