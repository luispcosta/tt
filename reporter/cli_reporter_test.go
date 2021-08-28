package reporter

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/persistence"
	"github.com/luispcosta/go-tt/utils"
)

type ReportLines struct {
	Lines []string
}

func (reportLines *ReportLines) Printer(a ...interface{}) (int, error) {
	for _, arg := range a {
		reportLines.Lines = append(reportLines.Lines, reflect.ValueOf(arg).String())
	}
	return 0, nil
}

func assertReportLineIsCorrect(t *testing.T, expected string, lines []string, position int) {
	got := lines[position]
	if expected != got {
		t.Errorf("Report line at position %v is wrong. Expected: %s, but instead got: %s", position, strings.TrimSpace(expected), strings.TrimSpace(got))
	}
}
