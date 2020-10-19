package reporter

import (
	"testing"

	"github.com/luispcosta/go-tt/configuration"
	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/persistence"
	"github.com/luispcosta/go-tt/utils"
)

func TestEmptyReporter(t *testing.T) {
	reporter := NewEmptyReporter()
	period, errPeriod := core.PeriodFromDateStrings("2020-10-10", "2020-10-11")
	if errPeriod != nil {
		t.Error("Raised an error initializing correct period")
	}

	clock := utils.NewLiveClock()
	config := configuration.NewConfig()
	repo := persistence.NewCustomJSONActivityRepository(utils.TestDataFolder, utils.LogTestFolder, *config, clock)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	errInitialize := reporter.Initialize(repo, period)

	if errInitialize != nil {
		if _, ok := errInitialize.(*utils.ReportNotImplementedError); !ok {
			t.Error("Error type is not ReportNotImplementdError")
		}
	} else {
		t.Error("Should have raised an error trying to intialize not implemented reporter")
	}

	errProduce := reporter.ProduceReport()

	if errProduce != nil {
		if _, ok := errProduce.(*utils.ReportNotImplementedError); !ok {
			t.Error("Error type is not ReportNotImplementdError")
		}
	} else {
		t.Error("Should have raised an error trying to produce a not implemented reporter")
	}
}
