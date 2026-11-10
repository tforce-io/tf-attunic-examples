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
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/tforce-io/tf-attunic-examples/go/todo"
)

// formatVersion is the version of the tasks file layout written by Save, and the only one Load reads.
const formatVersion = 1

// ErrUnsupportedVersion reports a tasks file whose layout version this program cannot read.
var ErrUnsupportedVersion = errors.New("unsupported tasks file version")

// fileContent is the persisted layout of a tasks file; change it only together with formatVersion.
type fileContent struct {
	Version int          `json:"version"`
	Tasks   []taskRecord `json:"tasks"`
}

type taskRecord struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

// JSONFile keeps a to-do list in a JSON file.
type JSONFile struct {
	path string
}

// NewJSONFile returns a JSONFile for the file at path; the file and its folder are created on the first Save.
func NewJSONFile(path string) *JSONFile {
	return &JSONFile{path: path}
}

// Path returns the path of the file.
func (f *JSONFile) Path() string {
	return f.path
}

// Load reads the list from the file; a missing file is an empty list.
// It returns an error if the file is not a valid tasks file of a supported version.
func (f *JSONFile) Load() (*todo.List, error) {
	data, err := os.ReadFile(f.path)
	if errors.Is(err, fs.ErrNotExist) {
		return &todo.List{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading tasks file: %w", err)
	}

	var content fileContent
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("tasks file %s: %w", f.path, err)
	}
	if content.Version != formatVersion {
		return nil, fmt.Errorf("tasks file %s: %w %d", f.path, ErrUnsupportedVersion, content.Version)
	}
	l, err := todo.NewList(toTasks(content.Tasks))
	if err != nil {
		return nil, fmt.Errorf("tasks file %s: %w", f.path, err)
	}
	return l, nil
}

// Save writes the list to the file, readable and writable by the current user only.
// It replaces the file in one step, so a failed Save leaves the previous content whole.
func (f *JSONFile) Save(l *todo.List) error {
	content := fileContent{Version: formatVersion, Tasks: toRecords(l.Tasks())}
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding tasks: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(f.path), 0o700); err != nil {
		return fmt.Errorf("creating tasks folder: %w", err)
	}
	if err := replaceFile(f.path, data); err != nil {
		return fmt.Errorf("writing tasks file: %w", err)
	}
	return nil
}

func toTasks(records []taskRecord) []todo.Task {
	tasks := make([]todo.Task, len(records))
	for i, r := range records {
		tasks[i] = todo.Task{ID: r.ID, Title: r.Title, Done: r.Done, CreatedAt: r.CreatedAt}
	}
	return tasks
}

func toRecords(tasks []todo.Task) []taskRecord {
	records := make([]taskRecord, len(tasks))
	for i, t := range tasks {
		records[i] = taskRecord{ID: t.ID, Title: t.Title, Done: t.Done, CreatedAt: t.CreatedAt}
	}
	return records
}

// replaceFile writes data to a temporary file beside path, then renames it over path.
func replaceFile(path string, data []byte) error {
	// CreateTemp creates the file with mode 0600.
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	if err := writeAndClose(tmp, data); err != nil {
		return errors.Join(err, os.Remove(tmp.Name()))
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return errors.Join(err, os.Remove(tmp.Name()))
	}
	return nil
}

func writeAndClose(f *os.File, data []byte) error {
	if _, err := f.Write(data); err != nil {
		return errors.Join(err, f.Close())
	}
	if err := f.Sync(); err != nil {
		return errors.Join(err, f.Close())
	}
	return f.Close()
}
