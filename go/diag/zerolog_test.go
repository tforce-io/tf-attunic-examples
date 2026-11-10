// Copyright (c) 2025 T-Force I/O
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package diag

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestParseLevelAcceptsNames(t *testing.T) {
	tests := []struct {
		name string
		want zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"WARN", zerolog.WarnLevel},
		{"disabled", zerolog.Disabled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.name)

			if err != nil {
				t.Fatalf("ParseLevel(%q): %v", tt.name, err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestParseLevelRejectsUnknownNames(t *testing.T) {
	for _, name := range []string{"", "verbose", "1"} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseLevel(name); err == nil {
				t.Errorf("ParseLevel(%q): expected an error, got nil", name)
			}
		})
	}
}

func TestNewLoggerWritesOnlyLevelAndAbove(t *testing.T) {
	var out bytes.Buffer
	logger := NewLogger(&out, zerolog.WarnLevel)

	logger.Info().Msg("info-sentinel")
	logger.Warn().Msg("warn-sentinel")

	if strings.Contains(out.String(), "info-sentinel") {
		t.Errorf("expected the info line skipped, got:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "WRN warn-sentinel") {
		t.Errorf("expected the warn line written, got:\n%s", out.String())
	}
}
