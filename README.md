# TODO
A simple command-line TODO list manager written with GO.

## Features
- Add, edit, view, and delete tasks
- Mark tasks as complete (closed) or incomplete (open)
- Archive tasks
- SQLite3 storage in `%Home%/.todo/data.sqlite3`

## Installation
Download the application using:
```bash
go install github.com/BennettB123/todo@latest
```

## Usage

### Create a Task
To create a new task, use:
```bash
todo new "Task description"
```

### List Tasks
To list tasks, use:
```bash
todo list
```
or
```bash
todo ls
```

Tasks are assigned a unique ID upon creation. Tasks will be displayed with their ID, a symbol representing open or closed, and the task's name. Example:
```
1   [ ] Task description
2   [X] This is another task
```

### Edit a Task
To edit the description of an existing task, use:
```bash
todo edit <task-id> "Updated task description"
```

### Close (complete) a Task
To mark a task as complete, use:
```bash
todo done <task-id>
```

### Open (incomplete) a Task
To mark a task as incomplete, use:
```bash
todo open <task-id>
```

### Delete a Task
To delete a task, use:
```bash
todo delete <task-id>
```
or
```bash
todo rm <task-id>
```

### Archive Tasks
To archive completed tasks, use:
```bash
todo archive <task-id>
```

Archiving tasks are not listed by default. Use `todo list -a` to view them.

### Support for Multiple Task IDs
The `done`, `open`, `delete`, and `archive` commands support multiple space-separated IDs. Example:
```bash
todo done 2 7 12
```