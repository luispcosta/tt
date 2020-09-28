package reporter

import "testing"

func TestIsAllowedFormat(t *testing.T) {
	if IsAllowedFormat("xpto") {
		t.Error("Should not allow format 'xpto")
	}

	if IsAllowedFormat("") {
		t.Error("Should not allow format ''")
	}

	if IsAllowedFormat("csvz") {
		t.Error("Should not allow format 'csvz'")
	}

	if !IsAllowedFormat("jSoN") {
		t.Error("Should allow format 'jSoN'")
	}

	if !IsAllowedFormat("CSV") {
		t.Error("Should allow format 'CSV'")
	}

	if !IsAllowedFormat("cli") {
		t.Error("Should allow format 'cli'")
	}
}

func TestCreate(t *testing.T) {
	if CreateReporter("xpto") != nil {
		t.Error("Should not create reporter for format format 'xpto")
	}

	if CreateReporter("") != nil {
		t.Error("Should not create reporter for format format ''")
	}

	if CreateReporter("csvz") != nil {
		t.Error("Should not create reporter for format format 'csvz'")
	}

	if CreateReporter("jSoN") == nil {
		t.Error("Should have had created reporter with format 'jSoN'")
	}

	if CreateReporter("CSV") == nil {
		t.Error("Should have had created reporter with format 'CSV'")
	}

	if CreateReporter("cli") == nil {
		t.Error("Should have had created reporter with format 'cli'")
	}
}
