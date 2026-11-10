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
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	majorVersion = 0
	minorVersion = 1
	patchVersion = 0
)

// gitCommit is the commit the program is built from, set at build time with
// -ldflags "-X github.com/tforce-io/tf-attunic-examples/go/engine.gitCommit=<hash>".
var gitCommit string

// controllerKey is the context key of a controller that commands use instead of loading one from the config folder.
type controllerKey struct{}

// BuildRootCmd returns the root command of the program, holding every command of the engine.
func BuildRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "todo",
		Long: fmt.Sprintf(`TFtodo v%s - keep a list of tasks to do.

Copyright (c) 2025 T-Force I/O.
Licensed under MIT license. See LICENSE file along with this program for more details.`,
			version()),
		Short:   "Keep a list of tasks to do.",
		Version: version(),

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	rootCmd.AddCommand(addCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(doneCmd())
	rootCmd.AddCommand(removeCmd())
	return rootCmd
}

// initApp returns the controller carried by ctx, or else a new one with the settings of the config folder.
func initApp(ctx context.Context) (*controller, error) {
	if c, ok := ctx.Value(controllerKey{}).(*controller); ok {
		return c, nil
	}
	return initController()
}

// version returns the version of the program, followed by the short commit hash when built with one.
func version() string {
	v := fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion)
	if len(gitCommit) >= 8 {
		v += "-" + gitCommit[:8]
	}
	return v
}
