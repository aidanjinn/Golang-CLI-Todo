
# Task Manager CLI

A simple command-line task management application written in Go that allows you to create, track, and manage tasks with dependencies.

## Features

- ✅ Create tasks with titles and notes
- ✅ Mark tasks as complete/incomplete
- ✅ Set task dependencies
- ✅ View tasks in a nicely formatted table
- ✅ Review task details including dependencies
- ✅ Update task titles and notes
- ✅ Delete tasks by ID, completed tasks, or all tasks
- ✅ Persistent storage using JSON

## Installation

1. Ensure you have Go installed (version 1.16 or higher recommended)
2. Clone this repository or download the source code
3. Install the required dependency

```bash
go get github.com/jedib0t/go-pretty/v6/table
```

4. Build the program:
```bash
go build
```

## Usage

Run the compiled binary:
```bash
./todo
```

### Available Commands

| Command  | Description                                      | Example                     |
|----------|--------------------------------------------------|-----------------------------|
| help     | Show available commands                          | `help`                      |
| post     | Create a new task                                | `post`                      |
| review   | View details of a specific task                  | `review`                    |
| update   | Modify an existing task                          | `update`                    |
| delete   | Remove tasks (by ID, completed, or all)          | `delete`                    |
| display  | Show all tasks in a table view                   | `display`                   |
| exit     | Quit the application                             | `exit`                      |

### Command Details

#### Post (Create a Task)
- Prompts for:
  - Task title
  - Notes (optional)
  - Dependencies (comma-separated task IDs, optional)

#### Review Task
- Shows complete details of a task including:
  - Title, completion status, creation date, notes
  - List of all dependencies with their details

#### Update Task
- Allows updating:
  - Completion status (done/not done)
  - Title
  - Notes

#### Delete Tasks
- Options:
  - Delete by ID (single task)
  - Delete all completed tasks
  - Delete all tasks

## Data Storage

Tasks are stored in a `tasks.json` file in the same directory as the executable. The file is automatically created when you add your first task.

## Dependencies

- [go-pretty](https://github.com/jedib0t/go-pretty) - For creating beautiful tables in the terminal

## License

This project is open source and available under the [MIT License](LICENSE).
