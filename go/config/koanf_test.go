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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadWithoutConfigFileUsesDefaults(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()

	cfg, err := Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Logger.Level != "warn" {
		t.Errorf("expected log.level warn, got %q", cfg.Logger.Level)
	}
	if want := filepath.Join(dir, "tasks.json"); cfg.Store.Path != want {
		t.Errorf("expected store.path %q, got %q", want, cfg.Store.Path)
	}
	if want := filepath.Join(dir, "todo.yaml"); cfg.ConfigFile != want {
		t.Errorf("expected config file %q, got %q", want, cfg.ConfigFile)
	}
}

func TestLoadReadsConfigFile(t *testing.T) {
	clearEnv(t)
	dir := writeConfigFile(t, "log:\n  level: debug\nstore:\n  path: /data/tasks.json\n")

	cfg, err := Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Logger.Level != "debug" {
		t.Errorf("expected log.level debug, got %q", cfg.Logger.Level)
	}
	if cfg.Store.Path != "/data/tasks.json" {
		t.Errorf("expected store.path /data/tasks.json, got %q", cfg.Store.Path)
	}
}

func TestLoadPrefersEnvironmentVariablesOverConfigFile(t *testing.T) {
	clearEnv(t)
	dir := writeConfigFile(t, "log:\n  level: debug\n")
	t.Setenv("TFTODO_LOG_LEVEL", "error")

	cfg, err := Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Logger.Level != "error" {
		t.Errorf("expected log.level error, got %q", cfg.Logger.Level)
	}
}

func TestLoadIgnoresEnvironmentVariablesNamingNoSetting(t *testing.T) {
	clearEnv(t)
	t.Setenv("TFTODO_LOG_COLOR", "true")

	cfg, err := Load(t.TempDir())

	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Logger.Level != "warn" {
		t.Errorf("expected log.level warn, got %q", cfg.Logger.Level)
	}
}

func TestLoadRejectsInvalidConfigFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"unknown key", "log:\n  color: true\n", "log.color"},
		{"section given a value", "log: debug\n", "unknown keys log"},
		{"malformed YAML", "log: [debug\n", "todo.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			dir := writeConfigFile(t, tt.content)

			_, err := Load(dir)

			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected the error to name %q, got %v", tt.want, err)
			}
		})
	}
}

// clearEnv removes every variable naming a setting for the duration of the test, so the environment of the developer cannot change its outcome.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{"TFTODO_LOG_LEVEL", "TFTODO_STORE_PATH"} {
		t.Setenv(name, "") // restores the original value after the test.
		if err := os.Unsetenv(name); err != nil {
			t.Fatalf("unsetting %s: %v", name, err)
		}
	}
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return dir
}
