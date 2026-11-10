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

package engine

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/tforce-io/tf-attunic-examples/go/config"
	"github.com/tforce-io/tf-attunic-examples/go/diag"
)

// controller holds what every command needs: the settings and the logger.
type controller struct {
	config *config.RootConfig
	logger zerolog.Logger
}

// initController returns a controller with the settings of the config folder, logging to stderr.
func initController() (*controller, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}
	cfg, err := config.Load(dir)
	if err != nil {
		return nil, err
	}
	level, err := diag.ParseLevel(cfg.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("setting log.level: %w", err)
	}
	return &controller{
		config: cfg,
		logger: diag.NewLogger(os.Stderr, level),
	}, nil
}

// commandLogger returns a logger tagging each line with the module and the command it comes from.
func (c *controller) commandLogger(module, command string) zerolog.Logger {
	return c.logger.With().Str("module", module).Str("command", command).Logger()
}
