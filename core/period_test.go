package core

import (
	"testing"
	"time"
)

func TestNumberOfDaysWhenSdAndEdAreSameDate(t *testing.T) {
	period := Period{Sd: time.Now(), Ed: time.Now()}

	if period.NumberOfDays() != 1 {
		t.Error("Number of perid days should be 1 when Sd and Ed are equal")
	}
}

func TestNumberOfDaysWhenSdIsBeforeEd(t *testing.T) {
	period := Period{Sd: time.Now().AddDate(0, 0, -1), Ed: time.Now()}

	if period.NumberOfDays() != 1 {
		t.Error("Number of perid days should be 1 when Sd is before Ed")
	}
}

func TestNumberOfDays(t *testing.T) {
	period1 := Period{Sd: time.Now(), Ed: time.Now().AddDate(0, 0, 10)}

	if period1.NumberOfDays() != 10 {
		t.Error("Number of perid days should be 10 when ed is 10 days after sd")
	}

	period2, err := PeriodFromDateStrings("2020-10-10", "2020-10-20")

	if err != nil {
		t.Error("Should not have failed creating period from '2020-10-10', '2020-10-20'")
	}

	if period2.NumberOfDays() != 10 {
		t.Error("Number of period days should be 10 when period created from: '2020-10-10', '2020-10-20'")
	}
}

func TestPeriodFromDateStringsWhenSdIsEmpty(t *testing.T) {
	_, err := PeriodFromDateStrings("", "2020-10-10")

	if err == nil {
		t.Error("Should have failed to create a period with invalid start date string (empty)")
	}
}

func TestPeriodFromDateStringsWhenEdIsEmpty(t *testing.T) {
	_, err := PeriodFromDateStrings("2020-10-10", "")

	if err == nil {
		t.Error("Should have failed to create a period with invalid end date string (empty)")
	}
}

func TestPeriodFromDateStringsWhenEdAndEdIsEmpty(t *testing.T) {
	_, err := PeriodFromDateStrings("", "")

	if err == nil {
		t.Error("Should have failed to create a period with invalid end date and start date string (empty)")
	}
}

func TestPeriodFromDateStringsWhenSdIsInvalidDate(t *testing.T) {
	_, err := PeriodFromDateStrings("2020-06-32", "2020-10-10")

	if err == nil {
		t.Error("Should have failed to create a period with invalid start date")
	}
}

func TestPeriodFromDateStringsWhenEdIsInvalidDate(t *testing.T) {
	_, err := PeriodFromDateStrings("2020-06-30", "xxxx")

	if err == nil {
		t.Error("Should have failed to create a period with invalid end date")
	}
}

func TestPeriodFromDateStringsWhenSdIsAfterEd(t *testing.T) {
	period, err := PeriodFromDateStrings("2020-06-30", "2020-06-29")

	if err != nil {
		t.Error("Should not have failed to create a period from valid Sd and Ed strings")
	}

	sd := period.Sd
	ed := period.Ed

	y1, m1, d1 := sd.Date()
	y2, m2, d2 := ed.Date()

	if y1 != 2020 && y2 != 2020 {
		t.Error("Invalid period created when sd is after ed")
	}

	if m1 != 6 && m2 != 6 {
		t.Error("Invalid period created when sd is after ed")
	}

	if d1 != 29 {
		t.Error("Invalid period created when sd is after ed")
	}

	if d2 != 30 {
		t.Error("Invalid period created when sd is after ed")
	}
}

func TestPeriodFromDateStringsWhenSdIsBeforeEd(t *testing.T) {
	period, err := PeriodFromDateStrings("2020-06-29", "2020-06-30")

	if err != nil {
		t.Error("Should not have failed to create a period from valid Sd and Ed strings")
	}

	sd := period.Sd
	ed := period.Ed

	y1, m1, d1 := sd.Date()
	y2, m2, d2 := ed.Date()

	if y1 != 2020 && y2 != 2020 {
		t.Error("Invalid period created when sd is before ed")
	}

	if m1 != 6 && m2 != 6 {
		t.Error("Invalid period created when sd is before ed")
	}

	if d1 != 29 {
		t.Error("Invalid period created when sd is before ed")
	}

	if d2 != 30 {
		t.Error("Invalid period created when sd is before ed")
	}
}

func TestPeriodFromKeyWordWithInvalidKeyword(t *testing.T) {
	period1 := PeriodFromKeyWord("")
	if period1 != PeriodFromKeyWord("day") {
		t.Error("Calling PeriodFromKeyWord with empty string didnt produce default period")
	}

	period2 := PeriodFromKeyWord("invalid")
	if period2 != PeriodFromKeyWord("day") {
		t.Error("Calling PeriodFromKeyWord with invalid keyword didnt produce default period")
	}
}

func TestPeriodFromKeyWordDay(t *testing.T) {
	period := PeriodFromKeyWord("day")
	sd := period.Sd
	ed := period.Ed
	now := time.Now()
	sdYear, sdMonth, sdDay := sd.Date()
	edYear, edMonth, edDay := ed.Date()
	currentYear, currentMonth, currentDay := now.Date()

	if sdYear != currentYear && edYear != currentYear {
		t.Error("Invalid period created for keyword 'day' (wrong year)")
	}

	if sdMonth != currentMonth && edMonth != currentMonth {
		t.Error("Invalid period created for keyword 'day' (wrong month)")
	}

	if sdDay != currentDay-1 {
		t.Error("Invalid period created for keyword 'day' (wrong sd day)")
	}

	if edDay != currentDay {
		t.Error("Invalid period created for keyword 'day' (wrong sd day)")
	}
}

func TestPeriodFromKeyWordWeek(t *testing.T) {
	period := PeriodFromKeyWord("week")
	sd := period.Sd
	ed := period.Ed
	now := time.Now()
	sdYear, sdMonth, sdDay := sd.Date()
	edYear, edMonth, edDay := ed.Date()
	currentYear, currentMonth, currentDay := now.Date()

	if sdYear != currentYear && edYear != currentYear {
		t.Error("Invalid period created for keyword 'week' (wrong year)")
	}

	if sdMonth != currentMonth && edMonth != currentMonth {
		t.Error("Invalid period created for keyword 'week' (wrong month)")
	}

	if sdDay != currentDay-7 {
		t.Error("Invalid period created for keyword 'week' (wrong sd day)")
	}

	if edDay != currentDay {
		t.Error("Invalid period created for keyword 'week' (wrong ed day)")
	}
}

func TestPeriodFromKeyWordYear(t *testing.T) {
	period := PeriodFromKeyWord("year")
	sd := period.Sd
	ed := period.Ed
	now := time.Now()
	sdYear, sdMonth, sdDay := sd.Date()
	edYear, edMonth, edDay := ed.Date()
	currentYear, currentMonth, currentDay := now.Date()

	if sdYear != currentYear-1 && edYear != currentYear {
		t.Error("Invalid period created for keyword 'year' (wrong year)")
	}

	if sdMonth != currentMonth && edMonth != currentMonth {
		t.Error("Invalid period created for keyword 'year' (wrong month)")
	}

	if sdDay != currentDay {
		t.Error("Invalid period created for keyword 'year' (wrong sd day)")
	}

	if edDay != currentDay {
		t.Error("Invalid period created for keyword 'year' (wrong ed day)")
	}
}
