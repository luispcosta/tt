package utils

import (
	"fmt"
	"testing"
)

func TestNotFoundError(t *testing.T) {
	notFoundError := NewNotFoundError("something")
	err := notFoundError.Error()

	if err != fmt.Sprintf("NotFound Error: something") {
		t.Error("NotFound error message does not match expected message")
	}
}

func TestReportNotImplementedError(t *testing.T) {
	notFoundError := NewReportNotImplementedError()
	err := notFoundError.Error()

	if err != fmt.Sprintf("This report format is not implemented yet") {
		t.Error("ReportNotImplementedError error message does not match expected message")
	}
}
