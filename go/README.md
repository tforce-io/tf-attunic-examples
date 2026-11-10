# TFtodo

Command line to-do list, and the Go example of [TFattunic](https://github.com/tforce-io/tf-attunic).

## Usage

```
todo <command> [arguments] [flags]
```

Run `todo --help` to list the commands, and add `--help` to any command to see its flags.

```sh
todo add Buy milk          # Added task 1: Buy milk
todo add "Call the bank"   # Added task 2: Call the bank
todo done 1                # Completed task 1: Buy milk
todo list                  #   1 [x] Buy milk
                           #   2 [ ] Call the bank
todo list --pending        #   2 [ ] Call the bank
todo remove 2              # Removed task 2: Call the bank
```

A task title has at most 200 characters and no line break. A new task takes the ID after the highest one in the list, so removing the last task frees its ID for the next task added.

### Configuration

Settings come from the config file `todo.yaml`, and environment variables override them; a setting found in neither keeps its default. The config folder holds the config file and, by default, the tasks file `tasks.json`:

- Windows: `%APPDATA%\TFtodo`.
- Linux and macOS: `~/.config/tftodo`.

| Key | Default | Meaning |
|-----|---------|---------|
| `log.level` | `warn` | Lowest level of the log lines written to stderr: `trace`, `debug`, `info`, `warn`, `error` or `disabled`. |
| `store.path` | `tasks.json` in the config folder | Tasks file; a relative path is relative to the current folder. |

```yaml
log:
  level: info
store:
  path: /home/me/Documents/tasks.json
```

An environment variable is named `TFTODO_` followed by the key in upper case with `_` in place of `.`, such as `TFTODO_LOG_LEVEL=debug`.

## Development

```sh
go build -o ./build/todo .   # append .exe on Windows
go test ./...
```

[.attunic/context.md](.attunic/context.md) holds the layout of the project, its architecture layers, and the commands to build, test and lint it. AI coding agents start from [AGENTS.md](AGENTS.md).

## License

TFtodo is licensed under MIT license. See [LICENSE](LICENSE) file for details.
