# Golang Task API Test Assignment

## Technical Specification

Develop an HTTP service in Go that allows creating, tracking, and deleting long-running asynchronous tasks (I/O bound) that simulate prolonged work (3-5 minutes). All data should be stored in memory without using external systems.

## API Description

`POST /api/tasks` - register a task in the system

`GET /api/tasks/{task_id}` - get information about a task by task_id

`DELETE /api/tasks/{task_id}` - delete a task from the system

## Configuration

- host and port - server address <host:port> - default "localhost:8080"
- file - Path to the data save/load file - default value "/output/task-db.json"
- interval - Data save interval to file (specified as whole number of seconds) - default value "3 seconds"

## Requirements

Go - v1.24.3

## Running

### Running from source code

1. Clone the repository
   `git clone https://github.com/NotNikita/golang-test-io-server.git` and enter the directory

2. Run the command `task check`

3. Start the server `go run ./cmd/main.go`

### Running using binary file

1. Download the binary file corresponding to your OS from Releases

2. Run the binary file
