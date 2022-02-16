package src

import (
	"testing"
)

func TestFormatBreakLength(t *testing.T) {
	length, unit := formatBreakLength(10)
	if length != 10 || unit != "seconds" {
		t.Error("TestFormatBreakLength: Unexpected result 10 seconds")
	}

	length, unit = formatBreakLength(100)
	if length != 1 || unit != "minute" {
		t.Error("TestFormatBreakLength: Unexpected result 100 seconds")
	}

	length, unit = formatBreakLength(120)
	if length != 2 || unit != "minutes" {
		t.Error("TestFormatBreakLength: Unexpected result 120 seconds")
	}

	length, unit = formatBreakLength(7200)
	if length != 2 || unit != "hours" {
		t.Error("TestFormatBreakLength: Unexpected result 7200 seconds")
	}

}
