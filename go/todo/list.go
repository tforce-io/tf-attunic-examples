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
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// MaxTitleLength is the maximum number of characters in a task title.
const MaxTitleLength = 200

var (
	// ErrNotFound reports that no task in the list has the requested ID.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID reports a task ID that is zero or negative.
	ErrInvalidID = errors.New("ID is not a positive number")

	// ErrDuplicateID reports a task ID held by another task of the same list.
	ErrDuplicateID = errors.New("ID is held by another task")

	// ErrEmptyTitle reports a task title that is empty or only white space.
	ErrEmptyTitle = errors.New("task title is empty")

	// ErrTitleTooLong reports a task title longer than MaxTitleLength characters.
	ErrTitleTooLong = fmt.Errorf("task title is longer than %d characters", MaxTitleLength)

	// ErrControlCharacter reports a task title holding a control character, such as a line break,
	// which would break the one-line-per-task listing or send escape sequences to the terminal.
	ErrControlCharacter = errors.New("task title has a control character")
)

// Task is one thing to do.
type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

// List is an ordered list of tasks, each with a unique, positive ID.
// The zero value is an empty list ready to use.
type List struct {
	tasks []Task
}

// NewList returns a list of tasks in the given order.
// It returns an error if a task has an invalid title, or an ID that is not positive or not unique.
func NewList(tasks []Task) (*List, error) {
	seen := make(map[int]bool, len(tasks))
	for _, t := range tasks {
		if t.ID <= 0 {
			return nil, fmt.Errorf("task %d: %w", t.ID, ErrInvalidID)
		}
		if seen[t.ID] {
			return nil, fmt.Errorf("task %d: %w", t.ID, ErrDuplicateID)
		}
		seen[t.ID] = true
		if _, err := normalizeTitle(t.Title); err != nil {
			return nil, fmt.Errorf("task %d: %w", t.ID, err)
		}
	}
	return &List{tasks: slices.Clone(tasks)}, nil
}

// Tasks returns a copy of every task, in order.
func (l *List) Tasks() []Task {
	return slices.Clone(l.tasks)
}

// Pending returns a copy of the tasks not done yet, in order.
func (l *List) Pending() []Task {
	var pending []Task
	for _, t := range l.tasks {
		if !t.Done {
			pending = append(pending, t)
		}
	}
	return pending
}

// Add appends a task with the given title, stripped of leading and trailing white space, and returns it.
// The new task's ID is one more than the highest ID in the list, so the ID of a removed last task is given again.
func (l *List) Add(title string, createdAt time.Time) (Task, error) {
	normalized, err := normalizeTitle(title)
	if err != nil {
		return Task{}, err
	}
	t := Task{ID: l.nextID(), Title: normalized, CreatedAt: createdAt}
	l.tasks = append(l.tasks, t)
	return t, nil
}

// Complete marks the task with the given ID as done and returns it; completing a done task changes nothing.
// It returns an error wrapping ErrNotFound if no task has that ID.
func (l *List) Complete(id int) (Task, error) {
	i, ok := l.indexOf(id)
	if !ok {
		return Task{}, fmt.Errorf("task %d: %w", id, ErrNotFound)
	}
	l.tasks[i].Done = true
	return l.tasks[i], nil
}

// Remove deletes the task with the given ID and returns it.
// It returns an error wrapping ErrNotFound if no task has that ID.
func (l *List) Remove(id int) (Task, error) {
	i, ok := l.indexOf(id)
	if !ok {
		return Task{}, fmt.Errorf("task %d: %w", id, ErrNotFound)
	}
	removed := l.tasks[i]
	l.tasks = slices.Delete(l.tasks, i, i+1)
	return removed, nil
}

func (l *List) indexOf(id int) (int, bool) {
	i := slices.IndexFunc(l.tasks, func(t Task) bool { return t.ID == id })
	return i, i >= 0
}

func (l *List) nextID() int {
	highest := 0
	for _, t := range l.tasks {
		highest = max(highest, t.ID)
	}
	return highest + 1
}

func normalizeTitle(title string) (string, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return "", ErrEmptyTitle
	}
	if utf8.RuneCountInString(trimmed) > MaxTitleLength {
		return "", ErrTitleTooLong
	}
	if strings.ContainsFunc(trimmed, unicode.IsControl) {
		return "", ErrControlCharacter
	}
	return trimmed, nil
}
