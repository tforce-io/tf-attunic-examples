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

package todo

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var createdAt = time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)

func TestAddAssignsNextID(t *testing.T) {
	tests := []struct {
		name   string
		tasks  []Task
		wantID int
	}{
		{"first task of an empty list", nil, 1},
		{"task after the highest ID", []Task{{ID: 5, Title: "a"}, {ID: 2, Title: "b"}}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := mustNewList(t, tt.tasks)

			got, err := l.Add("new", createdAt)

			if err != nil {
				t.Fatalf("Add: %v", err)
			}
			if got.ID != tt.wantID {
				t.Errorf("expected ID %d, got %d", tt.wantID, got.ID)
			}
		})
	}
}

func TestAddReusesIDOfRemovedLastTask(t *testing.T) {
	l := mustNewList(t, []Task{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}})
	if _, err := l.Remove(2); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	got, err := l.Add("c", createdAt)

	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if got.ID != 2 {
		t.Errorf("expected ID 2, got %d", got.ID)
	}
}

func TestAddStoresNormalizedTask(t *testing.T) {
	var l List

	got, err := l.Add("  Buy milk \n", createdAt)

	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	want := Task{ID: 1, Title: "Buy milk", Done: false, CreatedAt: createdAt}
	if got != want {
		t.Errorf("expected %+v, got %+v", want, got)
	}
	if tasks := l.Tasks(); len(tasks) != 1 || tasks[0] != want {
		t.Errorf("expected the list to hold only %+v, got %+v", want, tasks)
	}
}

func TestAddAcceptsTitleOfMaxLength(t *testing.T) {
	var l List
	title := strings.Repeat("é", MaxTitleLength) // multi-byte, so the limit is counted in characters, not bytes.

	got, err := l.Add(title, createdAt)

	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if got.Title != title {
		t.Errorf("expected the title kept whole, got %d characters", len([]rune(got.Title)))
	}
}

func TestAddRejectsInvalidTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  error
	}{
		{"empty", "", ErrEmptyTitle},
		{"only white space", " \t\n", ErrEmptyTitle},
		{"one character too long", strings.Repeat("a", MaxTitleLength+1), ErrTitleTooLong},
		{"line break inside", "Buy\nmilk", ErrControlCharacter},
		{"terminal escape sequence", "\x1b[31mBuy milk", ErrControlCharacter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l List

			_, err := l.Add(tt.title, createdAt)

			if !errors.Is(err, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, err)
			}
			if n := len(l.Tasks()); n != 0 {
				t.Errorf("expected no task added, got %d", n)
			}
		})
	}
}

func TestCompleteMarksTaskDone(t *testing.T) {
	l := mustNewList(t, []Task{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}})

	got, err := l.Complete(2)

	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if !got.Done {
		t.Error("expected the returned task done")
	}
	if pending := l.Pending(); len(pending) != 1 || pending[0].ID != 1 {
		t.Errorf("expected only task 1 pending, got %+v", pending)
	}
}

func TestCompleteKeepsDoneTaskDone(t *testing.T) {
	l := mustNewList(t, []Task{{ID: 1, Title: "a", Done: true}})

	got, err := l.Complete(1)

	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if !got.Done {
		t.Error("expected the task still done")
	}
}

func TestRemoveDeletesOnlyThatTask(t *testing.T) {
	l := mustNewList(t, []Task{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}, {ID: 3, Title: "c"}})

	got, err := l.Remove(2)

	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if got.Title != "b" {
		t.Errorf("expected the removed task returned, got %+v", got)
	}
	tasks := l.Tasks()
	if len(tasks) != 2 || tasks[0].ID != 1 || tasks[1].ID != 3 {
		t.Errorf("expected tasks 1 and 3 in order, got %+v", tasks)
	}
}

func TestChangesRejectUnknownID(t *testing.T) {
	tests := []struct {
		name   string
		change func(*List) (Task, error)
	}{
		{"complete", func(l *List) (Task, error) { return l.Complete(9) }},
		{"remove", func(l *List) (Task, error) { return l.Remove(9) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := mustNewList(t, []Task{{ID: 1, Title: "a"}})

			_, err := tt.change(l)

			if !errors.Is(err, ErrNotFound) {
				t.Errorf("expected %v, got %v", ErrNotFound, err)
			}
			if tasks := l.Tasks(); len(tasks) != 1 || tasks[0].Done {
				t.Errorf("expected the list unchanged, got %+v", tasks)
			}
		})
	}
}

func TestNewListRejectsInvalidTasks(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  error
	}{
		{"zero ID", []Task{{ID: 0, Title: "a"}}, ErrInvalidID},
		{"negative ID", []Task{{ID: -1, Title: "a"}}, ErrInvalidID},
		{"duplicate ID", []Task{{ID: 1, Title: "a"}, {ID: 1, Title: "b"}}, ErrDuplicateID},
		{"empty title", []Task{{ID: 1, Title: " "}}, ErrEmptyTitle},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewList(tt.tasks)

			if !errors.Is(err, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestTasksReturnsCopy(t *testing.T) {
	l := mustNewList(t, []Task{{ID: 1, Title: "a"}})

	l.Tasks()[0].Title = "changed"

	if got := l.Tasks()[0].Title; got != "a" {
		t.Errorf("expected the list unchanged by its caller, got title %q", got)
	}
}

func mustNewList(t *testing.T, tasks []Task) *List {
	t.Helper()
	l, err := NewList(tasks)
	if err != nil {
		t.Fatalf("NewList: %v", err)
	}
	return l
}
