#!/bin/bash
# This script will watch only existing file but not new file.
if [ -z "$1" ]; then
  echo "Please provide a directory path."
  exit 1
fi
if [ -z "$2" ]; then
  echo "Please provide a command to run."
  exit 1
fi
echo $API_PORT
if [ -z "$3" ]; then
  echo "Please provide a port number."
  exit 1
fi

watch_path=$1
server_command=$2
PORT=$3

watch_files=$(find "$watch_path" -type f \( \
  -name "*.templ" \
  -o -name "*.js" \
  -o \( -name "*.css" -a ! -name "output.css" \) \
  -o \( -name "*.go" -a ! -name "*_templ.go" \) \
  \))

if [ -z "$watch_files" ]; then
  echo "No matching files found."
  exit 1
fi

echo "Watching files in $watch_path"

start_server() {
  # Run the server command in the background
  $server_command &
  # Capture the PID of the last background command
  server_pid=$!
}

wait_for_port() {
  while lsof -i:$PORT >/dev/null 2>&1; do
    echo "Waiting for port $PORT to be free..."
    sleep 1
  done
}

kill_port() {
  lsof -ti:$PORT | xargs -r kill -9
}

# Function to clean up on exit
cleanup() {
  echo "Cleaning up..."
  kill_port
  if [ -n "$server_pid" ]; then
    kill -SIGTERM "$server_pid"
    wait "$server_pid" 2>/dev/null
  fi
  exit
}
# Trap SIGINT (Ctrl-C)
trap cleanup SIGINT

start_server
# Watch for changes and restart the server command on changes
fswatch --monitor=poll_monitor $watch_files | while read file_path; do
  echo "File changed: $file_path"
  
  # Kill the previous instance of the server command
  if [ -n "$server_pid" ]; then
    echo "Stopping server with PID $server_pid..."
    kill -SIGTERM "$server_pid" # kill "$server_pid"
    wait "$server_pid" 2>/dev/null
  fi
  
  kill_port
  wait_for_port
  # Restart the server command
  start_server
done
