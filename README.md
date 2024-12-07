# Homework Microservice

## Description

The Homework domain manages all information related to assignments for each course, including workflows for completing homework, submission tracking, and grading. Staff can post workflows to guide students through the homework process, which enhances the user experience (UX) for both students and staff.

## Setup and Usage

### Prerequisites

- **Go** (>= 1.19)
- **tmux** (optional, required only for running the bash script)

### Running the Microservice

#### 1. Run the Service Using the Makefile

- Start the server and client:
  ```bash
  make
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

## Contact

For inquiries or support, please contact [Your Name/Team] at [your-email@example.com](mailto\:your-email@example.com).

