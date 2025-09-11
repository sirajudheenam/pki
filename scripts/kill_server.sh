#!/bin/bash

PORT=8443

# Find the PID(s) of the process using the port
PIDS=$(lsof -ti tcp:$PORT)

if [ -z "$PIDS" ]; then
  echo "No process found using port $PORT"
else
  echo "Killing process(es) on port $PORT: $PIDS"
  # Kill the process(es)
  kill -9 $PIDS
  echo "Done."
fi
