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
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/tforce-io/tf-attunic-examples/go/store"
	"github.com/tforce-io/tf-attunic-examples/go/todo"
)

// taskModule carries out the commands that read and change the to-do list.
type taskModule struct {
	logger zerolog.Logger
	out    io.Writer
	store  *store.JSONFile
}

// newTaskModule returns a taskModule for the command named cmdName, printing its results to out.
func newTaskModule(c *controller, cmdName string, out io.Writer) *taskModule {
	return &taskModule{
		logger: c.commandLogger("task", cmdName),
		out:    out,
		store:  store.NewJSONFile(c.config.Store.Path),
	}
}

// add adds a task with the given title.
func (m *taskModule) add(title string) error {
	return m.change("Added", func(l *todo.List) (todo.Task, error) {
		return l.Add(title, time.Now().UTC())
	})
}

// complete marks the task with the given ID as done.
func (m *taskModule) complete(id int) error {
	return m.change("Completed", func(l *todo.List) (todo.Task, error) {
		return l.Complete(id)
	})
}

// remove deletes the task with the given ID.
func (m *taskModule) remove(id int) error {
	return m.change("Removed", func(l *todo.List) (todo.Task, error) {
		return l.Remove(id)
	})
}

// listAll prints every task.
func (m *taskModule) listAll() error {
	l, err := m.load()
	if err != nil {
		return err
	}
	return m.printTasks(l.Tasks())
}

// listPending prints the tasks not done yet.
func (m *taskModule) listPending() error {
	l, err := m.load()
	if err != nil {
		return err
	}
	return m.printTasks(l.Pending())
}

// change applies apply to the stored list, saves the list, then prints "<verb> task <ID>: <title>" for the task apply returns.
func (m *taskModule) change(verb string, apply func(*todo.List) (todo.Task, error)) error {
	l, err := m.load()
	if err != nil {
		return err
	}
	t, err := apply(l)
	if err != nil {
		return err
	}
	if err := m.store.Save(l); err != nil {
		return err
	}
	m.logger.Debug().Str("path", m.store.Path()).Int("id", t.ID).Msg("Saved tasks.")

	_, err = fmt.Fprintf(m.out, "%s task %d: %s\n", verb, t.ID, t.Title)
	return err
}

// load reads the stored list.
func (m *taskModule) load() (*todo.List, error) {
	l, err := m.store.Load()
	if err != nil {
		return nil, err
	}
	m.logger.Debug().Str("path", m.store.Path()).Int("count", len(l.Tasks())).Msg("Loaded tasks.")
	return l, nil
}

// printTasks prints one line per task: its ID, [x] when it is done or [ ] when it is not, and its title.
func (m *taskModule) printTasks(tasks []todo.Task) error {
	if len(tasks) == 0 {
		_, err := fmt.Fprintln(m.out, "No tasks.")
		return err
	}
	var lines strings.Builder
	for _, t := range tasks {
		mark := " "
		if t.Done {
			mark = "x"
		}
		fmt.Fprintf(&lines, "%3d [%s] %s\n", t.ID, mark, t.Title)
	}
	_, err := io.WriteString(m.out, lines.String())
	return err
}

// runTaskCommand runs action with a taskModule for cmd, marking the error it returns as a commandError.
func runTaskCommand(cmd *cobra.Command, action func(*taskModule) error) error {
	c, err := initApp(cmd.Context())
	if err != nil {
		return newCommandError(err)
	}
	return newCommandError(action(newTaskModule(c, cmd.Name(), cmd.OutOrStdout())))
}

// parseID returns the task ID written in arg, a positive whole number.
func parseID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("task ID %q is not a positive whole number", arg)
	}
	return id, nil
}

// addCmd returns the command adding a task.
func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <title>...",
		Short: "Add a task.",
		Long: fmt.Sprintf(`Add a task titled with the arguments joined by spaces, and print its ID.
The title has at most %d characters and no line break.`, todo.MaxTitleLength),
		Example: `  todo add Buy milk
  todo add "Call the bank"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := strings.Join(args, " ")
			return runTaskCommand(cmd, func(m *taskModule) error {
				return m.add(title)
			})
		},
	}
}

// listCmd returns the command listing the tasks.
func listCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the tasks.",
		Long: `List the tasks in the order they were added, one per line: its ID, [x] when
it is done or [ ] when it is not, and its title.`,
		Example: `  todo list
  todo list --pending`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			pendingOnly, err := cmd.Flags().GetBool("pending")
			if err != nil {
				return err
			}
			return runTaskCommand(cmd, func(m *taskModule) error {
				if pendingOnly {
					return m.listPending()
				}
				return m.listAll()
			})
		},
	}
	command.Flags().BoolP("pending", "p", false, "List only the tasks not done yet.")
	return command
}

// doneCmd returns the command marking a task as done.
func doneCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "done <id>",
		Short:   "Mark a task as done.",
		Long:    `Mark the task with the given ID as done; "todo list" shows the IDs.`,
		Example: `  todo done 3`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			return runTaskCommand(cmd, func(m *taskModule) error {
				return m.complete(id)
			})
		},
	}
}

// removeCmd returns the command removing a task.
func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <id>",
		Aliases: []string{"rm"},
		Short:   "Remove a task.",
		Long: `Remove the task with the given ID; "todo list" shows the IDs.
A new task takes the ID after the highest one in the list, so removing the last
task frees its ID for the next task added.`,
		Example: `  todo remove 3`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			return runTaskCommand(cmd, func(m *taskModule) error {
				return m.remove(id)
			})
		},
	}
}
