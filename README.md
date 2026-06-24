# task

A simple, priority-first CLI todo list manager. Tasks are persisted in a local SQLite database and automatically ordered by priority (block > high > medium > low) so you always focus on what matters most.

The core idea is based on [todo-manager](https://github.com/pykeras/BashUtils#todo-manager).

## Installation

```bash
go install github.com/arxenn/tasks@latest
```

Or build from source:

```bash
git clone https://github.com/arxenn/tasks.git
cd tasks
go build -o task .
```

## Shell Integration

Display your task list automatically on shell startup:

```bash
task shell enable   # add `task ls` to your shell config
task shell disable  # remove it
```

Supported shells: bash, zsh, fish, sh, ksh, and PowerShell (Windows).

## Usage

```bash
# Add a new task
task add Review pull requests -p high

# List pending tasks
task list

# Limit the number of shown tasks
task list -n 5

# Filter by priority
task list -p high

# Show completed tasks
task list --done

# Mark a task as done
task done 42

# Remove a task
task remove 15

# Clear completed tasks (or all with --all)
task clear
task clear --all
```

## Priorities

Tasks support four priority levels, listed highest to lowest:

| Priority | Flag          |
|----------|---------------|
| block    | `-p block`    |
| high     | `-p high`     |
| medium   | `-p medium`   |
| low      | `-p low`      |

If no priority is given, tasks default to `medium`.


## Data Storage

Tasks are stored in a SQLite database under your OS application data directory:

- **Linux:** `~/.local/share/tasks/`
- **macOS:** `~/Library/Application Support/tasks/`
- **Windows:** `%APPDATA%\tasks\`

## License

MIT
