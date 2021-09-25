package core

import (
	"strconv"

	"github.com/luispcosta/go-tt/utils"
)

// DurationFormat is a generic interface to represent activity durations
type DurationFormat interface {
	Format(int) string
}

const (
	Human   = "h"
	Seconds = "s"
	Minutes = "m"
	Hours   = "r"
)

func ParseDurationFormat(f string) DurationFormat {
	switch f {
	case Human:
		return HumanDurationFormat{}
	case Seconds:
		return SecondsDurationFormat{}
	case Minutes:
		return MinutesDurationFormat{}
	case Hours:
		return HoursDurationFormat{}
	default:
		return HumanDurationFormat{}
	}
}

type HumanDurationFormat struct{}

func (f HumanDurationFormat) Format(secondsDuration int) string {
	return utils.SecondsToHuman(secondsDuration)
}

type SecondsDurationFormat struct{}

func (f SecondsDurationFormat) Format(secondsDuration int) string {
	return strconv.Itoa(secondsDuration)
}

type MinutesDurationFormat struct{}

func (f MinutesDurationFormat) Format(secondsDuration int) string {
	return strconv.Itoa(secondsDuration / 60)
}

type HoursDurationFormat struct{}

func (f HoursDurationFormat) Format(secondsDuration int) string {
	return strconv.Itoa(secondsDuration / 60 / 60)
}
