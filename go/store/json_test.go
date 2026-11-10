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

package store

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/tforce-io/tf-attunic-examples/go/todo"
)

func TestLoadMissingFileReturnsEmptyList(t *testing.T) {
	f := NewJSONFile(filepath.Join(t.TempDir(), "tasks.json"))

	l, err := f.Load()

	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if n := len(l.Tasks()); n != 0 {
		t.Errorf("expected an empty list, got %d tasks", n)
	}
}

func TestSaveThenLoadKeepsTasks(t *testing.T) {
	want := []todo.Task{
		{ID: 1, Title: "Buy milk", Done: true, CreatedAt: time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)},
		{ID: 3, Title: "Call the bank", Done: false, CreatedAt: time.Date(2025, time.February, 3, 4, 5, 6, 0, time.UTC)},
	}
	f := NewJSONFile(filepath.Join(t.TempDir(), "tasks.json"))

	if err := f.Save(mustNewList(t, want)); err != nil {
		t.Fatalf("Save: %v", err)
	}
	l, err := f.Load()

	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	got := l.Tasks()
	if len(got) != len(want) {
		t.Fatalf("expected %d tasks, got %+v", len(want), got)
	}
	for i := range want {
		if !got[i].CreatedAt.Equal(want[i].CreatedAt) || got[i].ID != want[i].ID || got[i].Title != want[i].Title || got[i].Done != want[i].Done {
			t.Errorf("task %d: expected %+v, got %+v", i, want[i], got[i])
		}
	}
}

func TestSaveCreatesMissingFolder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "tasks.json")

	err := NewJSONFile(path).Save(&todo.List{})

	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected %s written: %v", path, err)
	}
}

func TestSaveLeavesNoTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	f := NewJSONFile(filepath.Join(dir, "tasks.json"))

	for range 2 { // the second Save replaces an existing file.
		if err := f.Save(&todo.List{}); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	if len(entries) != 1 || entries[0].Name() != "tasks.json" {
		t.Errorf("expected only tasks.json in the folder, got %v", entries)
	}
}

func TestSaveRestrictsFileToOwner(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows ignores Unix permission bits except the owner write bit")
	}
	path := filepath.Join(t.TempDir(), "tasks.json")

	if err := NewJSONFile(path).Save(&todo.List{}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("expected mode 0600, got %#o", got)
	}
}

func TestLoadRejectsInvalidFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    error // nil when any error will do.
	}{
		{"malformed JSON", `{"version": 1, "tasks": [`, nil},
		{"missing version", `{"tasks": []}`, ErrUnsupportedVersion},
		{"newer version", `{"version": 2, "tasks": []}`, ErrUnsupportedVersion},
		{"duplicate task ID", `{"version": 1, "tasks": [{"id": 1, "title": "a"}, {"id": 1, "title": "b"}]}`, todo.ErrDuplicateID},
		{"empty task title", `{"version": 1, "tasks": [{"id": 1, "title": ""}]}`, todo.ErrEmptyTitle},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "tasks.json")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("writing %s: %v", path, err)
			}

			_, err := NewJSONFile(path).Load()

			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func mustNewList(t *testing.T, tasks []todo.Task) *todo.List {
	t.Helper()
	l, err := todo.NewList(tasks)
	if err != nil {
		t.Fatalf("NewList: %v", err)
	}
	return l
}
