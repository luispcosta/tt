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
