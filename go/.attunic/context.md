# Project Context

Part of the project memory: human-readable memory of what this project is.

The reading and writing rules for this file and the rest of `.attunic/` are defined in [attunic-protocol.md](../.agents/attunic-protocol.md). The user fills this file in, by hand or by asking an agent to help; change it only on such a request and after the user explicitly confirms that change (see [attunic-protocol.md's Write Rules](../.agents/attunic-protocol.md#write-rules)).

## Abstract

A newcomer should be able to tell what the software does without knowing what it is built with. One paragraph that describes the domain, not the technology, in simple terms:
<!-- Project-specific / Abstract -->

TFtodo is a command line to-do list. Users add tasks, list them, mark them done, and remove them; the list is kept in a file in the user's config folder between runs. It is the Go example of TFattunic, so its layout and conventions matter as much as its features.

## Technical Specifications

A brief overview of the ingredients developers of this project should know.

Write one bullet per ingredient, as `- <Ingredient>: <value>.`, such as primary language, frameworks, architecture style, build system, test runner, linter, databases, environment prerequisites.

Write `None` as the value of an ingredient the project lacks but an agent would look for, such as a build system, test runner, or linter.

<!-- Project-specific / Technical Specifications -->

- Primary language: Go.
- Command line: Cobra.
- Configuration: Koanf, from a YAML file and environment variables.
- Logging: Zerolog.
- Databases: None.
- Build system: go.
- Test runner: go.
- Linter: gofmt, goimports, go vet.

## Layout

Layout of this project; [project-structure.md](../.agents/project-structure.md) holds the standard layouts of its languages and frameworks, kept for reference.

List of important folders, files and their roles in the project:

Write the list as a tree in a fenced block, one entry per line with its role as a `#` comment, mirroring the actual file system.

When the list is not empty and a top-level file/folder is missing from it, prompt the user to update it.

<!-- Project-specific / Layout -->

```
├── .agents/            # project AI agent instructions (attunic-protocol, coding-conventions, development-guideline, license-header, memory-template, project-structure, pull-request-guideline)
├── .attunic/           # project AI agent memory (context, memory, question)
├── app/                # wiring root of the program: assembles the root command and runs it
├── config/             # settings from default values, the config file, environment variables
├── diag/               # logging
├── engine/             # core of the program: Cobra commands, the controller, the modules carrying out the commands
├── store/              # JSON file persistence of the to-do list
├── todo/               # rules of the to-do list: tasks, their IDs and titles
├── .copier-answers.yml # project Copier answers
├── .editorconfig       # EditorConfig configuration
├── .gitignore          # git ignore patterns
├── AGENTS.md           # project AI agent entry point
├── CLAUDE.md           # project Claude Code entry point
├── LICENSE             # MIT license text
├── README.md           # project overview
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── main.go             # program entry point
```

## Layer Dependencies

List of architecture layers and their relationships:

- Declare the innermost architecture layer first.
- An architecture layer may depend only on itself and on the architecture layers in its `Depends on` column, all declared before it.
- If an architecture layer appears in the table below, the dependency graph must be enforced.

| Architecture layer | Path | Depends on | Must never contain |
|--------------------|------|------------|--------------------|
| <!-- Project-specific / Layer Dependencies --> | | | |
| domain | todo/ | itself | file, terminal or network access; settings; logging; Cobra types |
| support | config/, diag/ | itself | business rules; Cobra types |
| storage | store/ | domain | business rules; Cobra types |
| commands | engine/ | domain, support, storage | business rules; tasks file layout |
| wiring | app/, main.go | commands | business rules |

### Deliberate Exceptions

List of documented exceptions that cannot follow the architecture layers above:

- Any place the rules above are knowingly broken, with the reason and the condition for fixing it.
- Only exceptions to the architecture layers above belong here.
- An undocumented exception is a defect; a documented one is allowed, and revisited once its `Revisit-if` condition holds.

| What | Why | Revisit-if |
|------|-----|------------|
| <!-- Project-specific / Deliberate Exceptions --> | | |

## Supporting Tools

All the supporting tools declared here are opinionated for the sake of consistent results across development environments.

List each required tool/library under its own heading, followed by the commands for each use case.

Common use cases include: Build, Test, Lint, Format, Install.

<!-- Project-specific / Supporting Tools -->

### Build

- `go build -o ./build/todo .` - build the program from the root package, output in `./build` (append `.exe` on Windows).

### Test

- `go test ./...` - run the whole test suite with the go test runner.
- `go test -race ./...` - run the suite with the race detector; required for changes that touch concurrency, per [development-guideline.md](../.agents/development-guideline.md).
- Narrower: `go test ./engine/...`, `go test ./todo/...`, etc.

### Lint

- `gofmt -l .` - check formatting; must report no files.
- `gofmt -w .` - apply canonical formatting.
- `goimports -l -local github.com/tforce-io/tf-attunic-examples/go .` / `goimports -w -local github.com/tforce-io/tf-attunic-examples/go .` - check/fix import grouping (stdlib, then third-party, then local); the check must report no files.
- `go vet ./...` - go static analysis.
- `go mod tidy -diff` / `go mod tidy` - check/fix that `go.mod` and `go.sum` match the imports; the check must report no diff.

## Project-specific

<!-- Project-specific / Project-specific -->
