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
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// commandError is an error in a command's work, as opposed to an error in its command line.
type commandError struct {
	err error
}

// Error returns the inner error's message.
func (e *commandError) Error() string {
	return e.err.Error()
}

// Unwrap returns the inner error.
func (e *commandError) Unwrap() error {
	return e.err
}

// newCommandError wraps err into a commandError; a nil err stays nil.
func newCommandError(err error) error {
	if err == nil {
		return nil
	}
	return &commandError{err: err}
}

// PrintError prints err, returned by executing cmd, to w.
// An error in the command line is followed by a pointer to the command's help.
func PrintError(w io.Writer, cmd *cobra.Command, err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(w, cmd.ErrPrefix(), err)
	var ce *commandError
	if errors.As(err, &ce) {
		return
	}
	fmt.Fprintf(w, "Run '%s --help' for usage.\n", cmd.CommandPath())
}
