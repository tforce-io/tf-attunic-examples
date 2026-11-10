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
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// levels maps each accepted log level name to its level; zerolog.ParseLevel also accepts any number, which is not a documented setting.
var levels = map[string]zerolog.Level{
	"trace":    zerolog.TraceLevel,
	"debug":    zerolog.DebugLevel,
	"info":     zerolog.InfoLevel,
	"warn":     zerolog.WarnLevel,
	"error":    zerolog.ErrorLevel,
	"disabled": zerolog.Disabled,
}

// ParseLevel returns the log level named name, in any case: trace, debug, info, warn, error or disabled.
func ParseLevel(name string) (zerolog.Level, error) {
	level, ok := levels[strings.ToLower(name)]
	if !ok {
		return zerolog.NoLevel, fmt.Errorf("unknown log level %q, use trace, debug, info, warn, error or disabled", name)
	}
	return level, nil
}

// NewLogger returns a logger writing human-readable, timestamped lines of level and above to w.
func NewLogger(w io.Writer, level zerolog.Level) zerolog.Logger {
	console := zerolog.ConsoleWriter{Out: w, NoColor: true, TimeFormat: time.DateTime}
	return zerolog.New(console).Level(level).With().Timestamp().Logger()
}
