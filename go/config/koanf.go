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

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	configFileName = "todo.yaml"
	tasksFileName  = "tasks.json"

	// envPrefix starts the name of every environment variable that overrides a setting,
	// such as TFTODO_LOG_LEVEL for log.level.
	envPrefix = "TFTODO_"
)

// Dir returns the folder holding the config file and, by default, the tasks file:
// %APPDATA%\TFtodo on Windows, ~/.config/tftodo on Linux and macOS.
func Dir() (string, error) {
	if runtime.GOOS == "windows" {
		appData, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("locating config folder: %w", err)
		}
		return filepath.Join(appData, "TFtodo"), nil
	}
	// Not os.UserConfigDir, which follows XDG_CONFIG_HOME and is ~/Library/Application Support on macOS.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locating config folder: %w", err)
	}
	return filepath.Join(home, ".config", "tftodo"), nil
}

// Load returns the settings for the config folder dir. Each setting comes from the first of these that gives it:
// an environment variable named envPrefix followed by its key in upper case with "_" in place of ".",
// the config file todo.yaml in dir, and the default value.
//
// A missing config file is not an error, and environment variables naming no setting are ignored;
// a config file that is not valid YAML or has a key naming no setting is an error.
func Load(dir string) (*RootConfig, error) {
	k, err := defaults()
	if err != nil {
		return nil, err
	}
	known := k.Keys()

	configFile := filepath.Join(dir, configFileName)
	if err := mergeFile(k, configFile, known); err != nil {
		return nil, err
	}
	if err := mergeEnv(k, known); err != nil {
		return nil, err
	}

	var cfg RootConfig
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("decoding settings: %w", err)
	}
	cfg.ConfigFile = configFile
	if cfg.Store.Path == "" {
		cfg.Store.Path = filepath.Join(dir, tasksFileName)
	}
	return &cfg, nil
}

// defaults returns the default value of every setting; its keys are the only keys a setting can have.
func defaults() (*koanf.Koanf, error) {
	k := koanf.New(".")
	err := k.Load(structs.Provider(RootConfig{
		Logger: &LoggerConfig{
			Level: "warn",
		},
		Store: &StoreConfig{
			Path: "",
		},
	}, "koanf"), nil)
	if err != nil {
		return nil, fmt.Errorf("loading default settings: %w", err)
	}
	return k, nil
}

// mergeFile merges the settings of the YAML file at path into k, unless the file does not exist.
func mergeFile(k *koanf.Koanf, path string, known []string) error {
	loaded := koanf.New(".")
	err := loaded.Load(file.Provider(path), yaml.Parser())
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("config file %s: %w", path, err)
	}

	var unknown []string
	for _, key := range loaded.Keys() {
		if !slices.Contains(known, key) {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("config file %s: unknown keys %s", path, strings.Join(unknown, ", "))
	}
	return k.Merge(loaded)
}

// mergeEnv merges the settings of the environment variables named for a known key into k.
func mergeEnv(k *koanf.Koanf, known []string) error {
	vars := koanf.New(".")
	err := vars.Load(env.Provider(envPrefix, ".", func(name string) string {
		key := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(name, envPrefix)), "_", ".")
		if !slices.Contains(known, key) {
			return "" // the provider skips a variable mapped to no key.
		}
		return key
	}), nil)
	if err != nil {
		return fmt.Errorf("reading environment variables: %w", err)
	}
	return k.Merge(vars)
}
