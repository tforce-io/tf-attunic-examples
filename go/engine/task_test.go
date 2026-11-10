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
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/tforce-io/tf-attunic-examples/go/config"
	"github.com/tforce-io/tf-attunic-examples/go/todo"
)

const usageHint = "--help' for usage."

func TestAddPrintsNewTask(t *testing.T) {
	c := newTestController(t)

	out := mustRun(t, c, "add", " Buy", "milk ")

	if out != "Added task 1: Buy milk\n" {
		t.Errorf("unexpected output %q", out)
	}
}

func TestListShowsEachTaskWithItsStatus(t *testing.T) {
	c := newTestController(t)
	mustRun(t, c, "add", "Buy milk")
	mustRun(t, c, "add", "Call the bank")
	mustRun(t, c, "done", "1")

	out := mustRun(t, c, "list")

	if want := "  1 [x] Buy milk\n  2 [ ] Call the bank\n"; out != want {
		t.Errorf("expected:\n%s\ngot:\n%s", want, out)
	}
}

func TestListPendingHidesDoneTasks(t *testing.T) {
	c := newTestController(t)
	mustRun(t, c, "add", "Buy milk")
	mustRun(t, c, "add", "Call the bank")
	mustRun(t, c, "done", "1")

	out := mustRun(t, c, "list", "--pending")

	if want := "  2 [ ] Call the bank\n"; out != want {
		t.Errorf("expected:\n%s\ngot:\n%s", want, out)
	}
}

func TestListWithoutTasksSaysSo(t *testing.T) {
	c := newTestController(t)

	out := mustRun(t, c, "list")

	if out != "No tasks.\n" {
		t.Errorf("unexpected output %q", out)
	}
}

func TestRemoveDeletesTask(t *testing.T) {
	c := newTestController(t)
	mustRun(t, c, "add", "Buy milk")

	removed := mustRun(t, c, "remove", "1")
	listed := mustRun(t, c, "list")

	if removed != "Removed task 1: Buy milk\n" {
		t.Errorf("unexpected output of remove %q", removed)
	}
	if listed != "No tasks.\n" {
		t.Errorf("expected no task left, got %q", listed)
	}
}

func TestAddRejectsEmptyTitle(t *testing.T) {
	c := newTestController(t)

	out, err := runCommand(t, c, "add", " ")

	if !errors.Is(err, todo.ErrEmptyTitle) {
		t.Errorf("expected %v, got %v", todo.ErrEmptyTitle, err)
	}
	if !strings.Contains(out, "Error: task title is empty") {
		t.Errorf("expected the error printed, got %q", out)
	}
	if _, err := os.Stat(c.config.Store.Path); !os.IsNotExist(err) {
		t.Errorf("expected no tasks file written, got %v", err)
	}
}

func TestCommandsRejectUnknownTaskID(t *testing.T) {
	for _, command := range []string{"done", "remove"} {
		t.Run(command, func(t *testing.T) {
			c := newTestController(t)
			mustRun(t, c, "add", "Buy milk")

			out, err := runCommand(t, c, command, "9")

			if !errors.Is(err, todo.ErrNotFound) {
				t.Errorf("expected %v, got %v", todo.ErrNotFound, err)
			}
			if !strings.Contains(out, "Error: task 9: not found") {
				t.Errorf("expected the error printed, got %q", out)
			}
			if strings.Contains(out, usageHint) {
				t.Errorf("expected no usage hint for a valid command line, got %q", out)
			}
		})
	}
}

func TestCommandsPointInvalidTaskIDToHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"word", []string{"done", "first"}},
		{"zero", []string{"remove", "0"}},
		{"fraction", []string{"done", "1.5"}},
		{"missing", []string{"done"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestController(t)

			out, err := runCommand(t, c, tt.args...)

			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if want := "Run 'todo " + tt.args[0] + " " + usageHint; !strings.Contains(out, want) {
				t.Errorf("expected %q in the output, got %q", want, out)
			}
		})
	}
}

func TestUnknownCommandPointsToRootHelp(t *testing.T) {
	c := newTestController(t)

	out, err := runCommand(t, c, "frobnicate")

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if want := "Run 'todo " + usageHint; !strings.Contains(out, want) {
		t.Errorf("expected %q in the output, got %q", want, out)
	}
}

func TestCommandsReportInvalidTasksFile(t *testing.T) {
	c := newTestController(t)
	if err := os.WriteFile(c.config.Store.Path, []byte("not JSON"), 0o600); err != nil {
		t.Fatalf("writing %s: %v", c.config.Store.Path, err)
	}

	out, err := runCommand(t, c, "add", "Buy milk")

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if data, rerr := os.ReadFile(c.config.Store.Path); rerr != nil || string(data) != "not JSON" {
		t.Errorf("expected the tasks file left as it was, got %q (%v)", data, rerr)
	}
	if strings.Contains(out, usageHint) {
		t.Errorf("expected no usage hint for a valid command line, got %q", out)
	}
}

// newTestController returns a controller keeping its tasks in a fresh temporary folder, logging nothing.
func newTestController(t *testing.T) *controller {
	t.Helper()
	return &controller{
		config: &config.RootConfig{Store: &config.StoreConfig{Path: filepath.Join(t.TempDir(), "tasks.json")}},
		logger: zerolog.Nop(),
	}
}

// runCommand executes the command line args with c as the controller, and returns everything printed and the error.
func runCommand(t *testing.T, c *controller, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	root := BuildRootCmd()
	root.SetArgs(args)
	root.SetOut(&out)
	root.SetErr(&out)

	cmd, err := root.ExecuteContextC(context.WithValue(t.Context(), controllerKey{}, c))
	PrintError(&out, cmd, err)
	return out.String(), err
}

// mustRun executes the command line args with c as the controller, and returns everything printed; it fails the test on error.
func mustRun(t *testing.T, c *controller, args ...string) string {
	t.Helper()
	out, err := runCommand(t, c, args...)
	if err != nil {
		t.Fatalf("todo %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return out
}
