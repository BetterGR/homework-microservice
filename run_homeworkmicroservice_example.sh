#!/bin/bash

set -e  # Exit on any error

SERVER_FILE="server/main.go"
CLIENT_FILE="example/client.go"
PORT=1234

# Create a new tmux session for the server
tmux new-session -d -s homework_session -n server bash -c "echo -e '\n\n********** SERVER **********\n\n'; go run $SERVER_FILE"

# Split the window for the client
tmux split-window -h bash -c "echo -e '\n\n********** CLIENT **********\n\n'; go run $CLIENT_FILE --addr localhost:$PORT; read -p 'Press Enter to close this pane...'"

# Focus on the server pane initially
tmux select-pane -t 0

# Attach to the tmux session
tmux attach-session -t homework_session
