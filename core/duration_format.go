package core

import (
	"fmt"

	"github.com/luispcosta/go-tt/utils"
)

// DurationFormat is a generic interface to represent activity durations
type DurationFormat interface {
	Format(int) string
}

const (
	Auto    = "a"
	Seconds = "s"
	Minutes = "m"
	Hours   = "h"
)

func ParseDurationFormat(f string) DurationFormat {
	switch f {
	case Auto:
		return AutoDurationFormat{}
	case Seconds:
		return SecondsDurationFormat{}
	case Minutes:
		return MinutesDurationFormat{}
	case Hours:
		return HoursDurationFormat{}
	default:
		return AutoDurationFormat{}
	}
}

type AutoDurationFormat struct{}

func (f AutoDurationFormat) Format(secondsDuration int) string {
	return utils.SecondsToHuman(secondsDuration)
}

type SecondsDurationFormat struct{}

func (f SecondsDurationFormat) Format(secondsDuration int) string {
	return fmt.Sprintf("%v seconds", secondsDuration)
}

type MinutesDurationFormat struct{}

func (f MinutesDurationFormat) Format(secondsDuration int) string {
	minutes := secondsDuration / 60
	return fmt.Sprintf("%v minutes", minutes)
}

type HoursDurationFormat struct{}

func (f HoursDurationFormat) Format(secondsDuration int) string {
	hours := secondsDuration / 60 / 60
	return fmt.Sprintf("%v hours", hours)
}
