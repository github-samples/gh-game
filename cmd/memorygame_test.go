package cmd

import (
"bytes"
"strings"
"testing"

"github.com/chrisreddington/gh-game/internal/memorygame"
)

func TestDisplayUserSequence(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Test that we can display a sequence of colors on a single line
	colorStyles := map[memorygame.Color]string{
		memorygame.Red:    "Red",
		memorygame.Yellow: "Yellow",
		memorygame.Green:  "Green",
		memorygame.Blue:   "Blue",
	}

	// Test with an empty sequence
	userSequence := []memorygame.Color{}
	buf.Reset()
	for _, color := range userSequence {
		cstr := string(color)
		if len(cstr) > 0 {
			cstr = strings.ToUpper(cstr[:1]) + cstr[1:]
		}
		buf.WriteString(colorStyles[color] + " ")
	}
	if len(buf.String()) != 0 {
		t.Errorf("Expected empty string for empty sequence, got %q", buf.String())
	}

	// Test with a single color
	userSequence = []memorygame.Color{memorygame.Red}
	buf.Reset()
	for _, color := range userSequence {
		cstr := string(color)
		if len(cstr) > 0 {
			cstr = strings.ToUpper(cstr[:1]) + cstr[1:]
		}
		buf.WriteString(colorStyles[color] + " ")
	}
	expected := "Red "
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}

	// Test with multiple colors
	userSequence = []memorygame.Color{memorygame.Red, memorygame.Blue, memorygame.Green}
	buf.Reset()
	for _, color := range userSequence {
		cstr := string(color)
		if len(cstr) > 0 {
			cstr = strings.ToUpper(cstr[:1]) + cstr[1:]
		}
		buf.WriteString(colorStyles[color] + " ")
	}
	expected = "Red Blue Green "
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
cd /Users/chrisreddington/Documents/Code/gh-game && go test -v ./cmd/... ./internal/...}
