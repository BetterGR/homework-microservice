# Homework Microservice

## Description

The Homework domain manages all information related to assignments for each course, including workflows for completing homework, submission tracking, and grading. Staff can post workflows to guide students through the homework process, which enhances the user experience (UX) for both students and staff.

## Setup and Usage

### Prerequisites

- **Go** (>= 1.23.4)
- **tmux** (optional, required only for running the bash script)
- **Docker** (optional, required for container operations)

### Available Make Commands

The following commands are available through the Makefile:

- `make all` - Build and check everything (proto, gomod, fmt, vet, lint)
- `make proto` - Generate Go code from proto files
- `make gomod` - Manage Go modules (tidy and verify)
- `make fmt` - Format Go code using gofumpt and gci
- `make vet` - Run Go vet checks on code
- `make lint` - Run golangci-lint checks
- `make build` - Build the server binary
- `make run` - Run the server
- `make docker-build` - Build Docker image
- `make docker-push` - Push Docker image to registry
- `make clean` - Clean up generated files
- `make help` - Display available make commands

### Running the Microservice

#### 1. Using Make Commands

```bash
# Build and run all checks
make

# Run the server
make run

# Build Docker image
make docker-build
```

#### 2. Manual Setup (Without tmux or Script)

- Start the server manually:
  ```bash
  go run server/main.go
  ```
- In another terminal, run the client:
  ```bash
  go run example/client.go --addr localhost:1234
  ```

#### 3. Using the Bash Script

If you prefer to use the bash script, ensure `tmux` is installed:

```bash
sudo apt install tmux
```

Run the script:

```bash
./run_homeworkmicroservice_example.sh
```

### Exiting the Microservice

- For `tmux` sessions:
  ```bash
  tmux attach -t homework_session
  Ctrl+B, then press D (to detach)
  tmux kill-session -t homework_session
  ```
- Alternatively, you can stop the server directly by pressing Ctrl+C in the terminal where it is running


