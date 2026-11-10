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

package config

// RootConfig contains all the settings of the program.
type RootConfig struct {
	// ConfigFile is the path of the config file the settings were read from, whether it exists or not.
	ConfigFile string `koanf:"-"`

	Logger *LoggerConfig `koanf:"log"`
	Store  *StoreConfig  `koanf:"store"`
}

// LoggerConfig contains the settings for logging.
type LoggerConfig struct {
	// Level is the lowest level of the log lines written: trace, debug, info, warn, error or disabled.
	Level string `koanf:"level"`
}

// StoreConfig contains the settings for keeping tasks.
type StoreConfig struct {
	// Path is the tasks file. Load sets it to tasks.json in the config folder when no setting gives it.
	Path string `koanf:"path"`
}
