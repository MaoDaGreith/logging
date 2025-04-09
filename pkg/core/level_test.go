package core

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{Debug, "DEBUG"},
		{Info, "INFO"},
		{Warning, "WARNING"},
		{Error, "ERROR"},
		{Level(99), "UNKNOWN(99)"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.level.String(); got != test.expected {
				t.Errorf("Level.String() = %q, want %q", got, test.expected)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		levelStr string
		expected Level
		wantErr  bool
	}{
		{"DEBUG", Debug, false},
		{"debug", Debug, false},
		{"INFO", Info, false},
		{"info", Info, false},
		{"WARNING", Warning, false},
		{"warning", Warning, false},
		{"WARN", Warning, false},
		{"warn", Warning, false},
		{"ERROR", Error, false},
		{"error", Error, false},
		{"ERR", Error, false},
		{"err", Error, false},
		{"invalid", Info, true},
		{"", Info, true},
	}

	for _, test := range tests {
		t.Run(test.levelStr, func(t *testing.T) {
			got, err := ParseLevel(test.levelStr)

			// Check error condition
			if (err != nil) != test.wantErr {
				t.Errorf("ParseLevel(%q) error = %v, wantErr %v", test.levelStr, err, test.wantErr)
				return
			}

			// Check result if no error expected
			if !test.wantErr && got != test.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", test.levelStr, got, test.expected)
			}
		})
	}
}
