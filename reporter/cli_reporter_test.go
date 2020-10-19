package reporter

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/luispcosta/go-tt/configuration"
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

func TestProduceReportWithNoActivityOnPeriod(t *testing.T) {
	period, errPeriod := core.PeriodFromDateStrings("2020-10-10", "2020-10-11")
	if errPeriod != nil {
		t.Error("Raised an error initializing correct period")
	}

	config := configuration.NewConfig()
	repo := persistence.NewCustomJSONActivityRepository(utils.TestDataFolder, utils.LogTestFolder, *config, utils.NewLiveClock())
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	reporter := NewCliReporter()
	errInitialize := reporter.Initialize(repo, period)
	if errInitialize != nil {
		t.Error("Should not have failed to initialize CLI reporter")
	}

	errProduce := reporter.ProduceReport()

	if errProduce == nil {
		t.Error("Should have failed to produce report when no reports exist in period")
	}
}

func TestProduceReportWithSomeActivitiesInPeriod(t *testing.T) {
	defer utils.ClearTestFolder()
	defer utils.ClearLogTestFolder()

	printer := ReportLines{}
	period, errPeriod := core.PeriodFromDateStrings("2020-10-10", "2020-10-15")
	if errPeriod != nil {
		t.Error("Raised an error initializing correct period")
	}

	mockedNow, errParseDate := time.Parse("2006-01-02", "2020-10-12")
	clock := utils.NewMockedClock(mockedNow)
	if errParseDate != nil {
		t.Error("Should not have failed parsing a correct date for clock")
	}

	config := configuration.NewConfig()
	repo := persistence.NewCustomJSONActivityRepository(utils.TestDataFolder, utils.LogTestFolder, *config, clock)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity := core.Activity{}
	activity.Name = "Some_Name"
	errUpdating := repo.Update(activity)

	if errUpdating != nil {
		t.Error("Should not have failed creating a valid activity")
	}

	repo.Start(activity)
	advancedTime := mockedNow.Add(time.Hour * 2)
	clock.SetNow(advancedTime)
	repo.Stop(activity)

	timeBeforePeriod := mockedNow.Add(-(time.Hour * 24) * 10)
	clock.SetNow(timeBeforePeriod)
	repo.Start(activity)
	repo.Stop(activity)

	reporter := NewCustomCLIReporter(printer.Printer)
	errInitialize := reporter.Initialize(repo, period)
	if errInitialize != nil {
		t.Error("Should not have failed to initialize CLI reporter")
	}

	errProduce := reporter.ProduceReport()

	if errProduce != nil {
		t.Error("Should not have failed to produce report when no reports exist in period")
	}

	linesPrinted := printer.Lines

	expectedFirstLine := "Number of activities found: 1\n"
	assertReportLineIsCorrect(t, expectedFirstLine, linesPrinted, 0)

	expectedThirdLine := "  some_name () - 7200s/2.00h - (100.00%)\n"
	assertReportLineIsCorrect(t, expectedThirdLine, linesPrinted, 2)

	expectedFifthLine := "Date: 2020-10-12\n"
	assertReportLineIsCorrect(t, expectedFifthLine, linesPrinted, 4)

	expectedSeventhLine := "  Total Duration: 7200s/2.00h\n"
	assertReportLineIsCorrect(t, expectedSeventhLine, linesPrinted, 6)
}

func TestProduceReportWithSomeActivitiesInPeriodAndActivitiesWithZeroDuration(t *testing.T) {
	defer utils.ClearTestFolder()
	defer utils.ClearLogTestFolder()

	printer := ReportLines{}
	period, errPeriod := core.PeriodFromDateStrings("2020-10-10", "2020-10-15")
	if errPeriod != nil {
		t.Error("Raised an error initializing correct period")
	}

	mockedNow, errParseDate := time.Parse("2006-01-02", "2020-10-12")
	clock := utils.NewMockedClock(mockedNow)
	if errParseDate != nil {
		t.Error("Should not have failed parsing a correct date for clock")
	}

	config := configuration.NewConfig()
	repo := persistence.NewCustomJSONActivityRepository(utils.TestDataFolder, utils.LogTestFolder, *config, clock)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity := core.Activity{}
	activity.Name = "Some_Name"
	errUpdating := repo.Update(activity)

	if errUpdating != nil {
		t.Error("Should not have failed creating a valid activity")
	}

	otherActivity := core.Activity{}
	otherActivity.Name = "Activity2"
	otherActivity.Description = "Some big and long description"
	otherActivity.Alias = "a2"
	errUpdating2 := repo.Update(otherActivity)

	if errUpdating2 != nil {
		t.Error("Should not have failed creating a valid activity")
	}

	yetAnotherActivity := core.Activity{}
	yetAnotherActivity.Name = "Activity3"
	yetAnotherActivity.Description = "Some big and long description for yet another activity"
	yetAnotherActivity.Alias = "a3"
	errUpdating3 := repo.Update(yetAnotherActivity)

	if errUpdating3 != nil {
		t.Error("Should not have failed creating a valid activity")
	}

	repo.Start(activity)
	advancedTime := mockedNow.Add(time.Hour * 2)
	clock.SetNow(advancedTime)
	repo.Stop(activity)

	timeBeforePeriod := mockedNow.Add(-(time.Hour * 24) * 10) // this should not appear because its before period
	clock.SetNow(timeBeforePeriod)
	repo.Start(activity)
	repo.Stop(activity)

	clock.SetNow(mockedNow)
	repo.Start(otherActivity)
	repo.Stop(otherActivity) // 0 duration

	timeAfterPeriod := mockedNow.Add(time.Hour * 24 * 10) // this should not appear because its after period
	clock.SetNow(timeAfterPeriod)
	repo.Start(yetAnotherActivity)
	repo.Stop(yetAnotherActivity)

	advancedTime2 := mockedNow.Add(time.Hour * 48)
	clock.SetNow(advancedTime2)
	repo.Start(yetAnotherActivity)

	advancedTime3 := advancedTime2.Add(time.Minute * 50)
	clock.SetNow(advancedTime3)
	repo.Stop(yetAnotherActivity)

	reporter := NewCustomCLIReporter(printer.Printer)
	errInitialize := reporter.Initialize(repo, period)
	if errInitialize != nil {
		t.Error("Should not have failed to initialize CLI reporter")
	}

	errProduce := reporter.ProduceReport()

	if errProduce != nil {
		t.Error("Should not have failed to produce report when no reports exist in period")
	}

	linesPrinted := printer.Lines

	expectedFirstLine := "Number of activities found: 3\n"
	assertReportLineIsCorrect(t, expectedFirstLine, linesPrinted, 0)

	expectedThirdLine := "  some_name () - 7200s/2.00h - (70.59%)\n"
	assertReportLineIsCorrect(t, expectedThirdLine, linesPrinted, 2)

	expectedFourth := "  activity3 (Some big and long description for yet another activity) - 3000s/0.83h - (29.41%)\n"
	assertReportLineIsCorrect(t, expectedFourth, linesPrinted, 3)

	expectedFifthLine := "  activity2 (Some big and long description) - 0s/0.00h - (0.00%)\n"
	assertReportLineIsCorrect(t, expectedFifthLine, linesPrinted, 4)

	expectedSeventhLine := "Date: 2020-10-12\n"
	assertReportLineIsCorrect(t, expectedSeventhLine, linesPrinted, 6)
	expectedLinePos7 := "  Activity: some_name\n"
	assertReportLineIsCorrect(t, expectedLinePos7, linesPrinted, 7)
	expectedLinePos8 := "  Total Duration: 7200s/2.00h\n"
	assertReportLineIsCorrect(t, expectedLinePos8, linesPrinted, 8)
	expectedLinePos9 := "  Activity: activity2\n"
	assertReportLineIsCorrect(t, expectedLinePos9, linesPrinted, 9)
	expectedLinePos10 := "  Total Duration: 0s/0.00h\n"
	assertReportLineIsCorrect(t, expectedLinePos10, linesPrinted, 10)

	expectedLinePos11 := "Date: 2020-10-14\n"
	assertReportLineIsCorrect(t, expectedLinePos11, linesPrinted, 11)
	expectedLinePos12 := "  Activity: activity3\n"
	assertReportLineIsCorrect(t, expectedLinePos12, linesPrinted, 12)
	expectedLinePos13 := "  Total Duration: 3000s/0.83h\n"
	assertReportLineIsCorrect(t, expectedLinePos13, linesPrinted, 13)
}

func TestProduceReportWithSomeActivitiesInPeriodAndSomeActivitiesCannotBeFound(t *testing.T) {
	defer utils.ClearTestFolder()
	defer utils.ClearLogTestFolder()

	printer := ReportLines{}
	period, errPeriod := core.PeriodFromDateStrings("2020-10-10", "2020-10-15")
	if errPeriod != nil {
		t.Error("Raised an error initializing correct period")
	}

	mockedNow, errParseDate := time.Parse("2006-01-02", "2020-10-12")
	clock := utils.NewMockedClock(mockedNow)
	if errParseDate != nil {
		t.Error("Should not have failed parsing a correct date for clock")
	}

	config := configuration.NewConfig()
	repo := persistence.NewCustomJSONActivityRepository(utils.TestDataFolder, utils.LogTestFolder, *config, clock)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity := core.Activity{}
	activity.Name = "Some_Name"
	errUpdating := repo.Update(activity)

	if errUpdating != nil {
		t.Error("Should not have failed creating a valid activity")
	}

	otherActivity := core.Activity{}
	otherActivity.Name = "Activity2"
	otherActivity.Description = "Some big and long description"
	otherActivity.Alias = "a2"
	errUpdating2 := repo.Update(otherActivity)

	if errUpdating2 != nil {
		t.Error("Should not have failed creating a valid activity")
	}

	yetAnotherActivity := core.Activity{}
	yetAnotherActivity.Name = "Activity3"
	yetAnotherActivity.Description = "Some big and long description for yet another activity"
	yetAnotherActivity.Alias = "a3"
	errUpdating3 := repo.Update(yetAnotherActivity)

	if errUpdating3 != nil {
		t.Error("Should not have failed creating a valid activity")
	}

	repo.Start(activity)
	advancedTime := mockedNow.Add(time.Hour * 2)
	clock.SetNow(advancedTime)
	repo.Stop(activity)

	timeBeforePeriod := mockedNow.Add(-(time.Hour * 24) * 10) // this should not appear because its before period
	clock.SetNow(timeBeforePeriod)
	repo.Start(activity)
	repo.Stop(activity)

	clock.SetNow(mockedNow)
	repo.Start(otherActivity)
	repo.Stop(otherActivity) // 0 duration

	timeAfterPeriod := mockedNow.Add(time.Hour * 24 * 10) // this should not appear because its after period
	clock.SetNow(timeAfterPeriod)
	repo.Start(yetAnotherActivity)
	repo.Stop(yetAnotherActivity)

	advancedTime2 := mockedNow.Add(time.Hour * 48)
	clock.SetNow(advancedTime2)
	repo.Start(yetAnotherActivity)

	advancedTime3 := advancedTime2.Add(time.Minute * 50)
	clock.SetNow(advancedTime3)
	repo.Stop(yetAnotherActivity)

	repo.Delete(yetAnotherActivity.Name)

	reporter := NewCustomCLIReporter(printer.Printer)
	errInitialize := reporter.Initialize(repo, period)
	if errInitialize != nil {
		t.Error("Should not have failed to initialize CLI reporter")
	}

	errProduce := reporter.ProduceReport()

	if errProduce != nil {
		t.Error("Should not have failed to produce report when no reports exist in period")
	}

	linesPrinted := printer.Lines

	expectedFirstLine := "Number of activities found: 2\n"
	assertReportLineIsCorrect(t, expectedFirstLine, linesPrinted, 0)

	expectedThirdLine := "  some_name () - 7200s/2.00h - (70.59%)\n"
	assertReportLineIsCorrect(t, expectedThirdLine, linesPrinted, 2)

	expectedFifthLine := "  activity2 (Some big and long description) - 0s/0.00h - (0.00%)\n"
	assertReportLineIsCorrect(t, expectedFifthLine, linesPrinted, 3)

	expectedSeventhLine := "Date: 2020-10-12\n"
	assertReportLineIsCorrect(t, expectedSeventhLine, linesPrinted, 5)
	expectedLinePos7 := "  Activity: some_name\n"
	assertReportLineIsCorrect(t, expectedLinePos7, linesPrinted, 6)
	expectedLinePos8 := "  Total Duration: 7200s/2.00h\n"
	assertReportLineIsCorrect(t, expectedLinePos8, linesPrinted, 7)
	expectedLinePos9 := "  Activity: activity2\n"
	assertReportLineIsCorrect(t, expectedLinePos9, linesPrinted, 8)
	expectedLinePos10 := "  Total Duration: 0s/0.00h\n"
	assertReportLineIsCorrect(t, expectedLinePos10, linesPrinted, 9)
}
