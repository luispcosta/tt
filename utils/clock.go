package utils

import (
	"errors"
	"time"
)

// Clock is a general interface to control the time.
type Clock interface {
	Now() time.Time
	SetNow(time.Time) error
}

// LiveClock represents a "correct" clock (time has not been tempered with)
type LiveClock struct{}

// MockedClock represents a clock where the current time is changed
type MockedClock struct {
	MockedNow time.Time
}

// NewMockedClock returns a new mocked clock
func NewMockedClock(now time.Time) *MockedClock {
	return &MockedClock{MockedNow: now}
}

// NewClickClock returns a new, correct, clock
func NewLiveClock() *LiveClock {
	return &LiveClock{}
}

// Now returns the current time in the correct clock
func (clock LiveClock) Now() time.Time {
	return time.Now()
}

// SetNow will not work on a "correct" clock, since it's always moving forwards
func (clock LiveClock) SetNow(time time.Time) error {
	return errors.New("Cannot set 'Now' for a live clock")
}

// Now returns the current time in the mocked clock
func (clock *MockedClock) Now() time.Time {
	return clock.MockedNow
}

// SetNow changes the current time in the mocked clock
func (clock *MockedClock) SetNow(time time.Time) error {
	clock.MockedNow = time
	return nil
}
