package core

import (
	"testing"
)

func TestValidNameWhenEmpty(t *testing.T) {
	activity := Activity{Name: ""}

	err := activity.ValidateName()

	if err == nil {
		t.Fatal("Should have failed when name is blank")
	}
}

func TestValidNameWhenNull(t *testing.T) {
	activity := Activity{}

	err := activity.ValidateName()

	if err == nil {
		t.Fatal("Should have failed when name is missing")
	}
}

func TestValidNameHasWrongSpaces(t *testing.T) {
	activity := Activity{Name: "Some     name"}

	err := activity.ValidateName()

	if err == nil {
		t.Fatal("Should have failed when name contains consecutive spaces")
	}
}

func TestValidNameHasWithSpaces(t *testing.T) {
	activity := Activity{Name: "Some name"}

	err := activity.ValidateName()

	if err == nil {
		t.Fatal("Should have failed when name contains spaces")
	}
}

func TestValidateNameWhenReservedNameIsUsed(t *testing.T) {
	activity := Activity{Name: "inDeX"}

	err := activity.ValidateName()

	if err == nil {
		t.Fatal("Should have failed when using reserved names for activity name")
	}
}

func TestValidateNameWhenNameHasDots(t *testing.T) {
	activity := Activity{Name: "someActivity.json"}

	err := activity.ValidateName()

	if err == nil {
		t.Fatal("Should have failed when using dots on activity name")
	}
}

func TestValidateNameWhenNameHasBackSlashes(t *testing.T) {
	activity := Activity{Name: "someActivity/other/"}

	err := activity.ValidateName()

	if err == nil {
		t.Fatal("Should have failed when using slashes on activity name")
	}
}

func TestValidNameHasWithCorrectChars(t *testing.T) {
	activity := Activity{Name: "Name123_-"}

	err := activity.ValidateName()

	if err != nil {
		t.Fatal("Should not have failed with correct name chars")
	}
}
