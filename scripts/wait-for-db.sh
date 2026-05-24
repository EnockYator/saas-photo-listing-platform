#!/bin/bash
# Wait for database to be ready

set -e

TIMEOUT=30
HOST=$1
PORT=$2
shift 2
CMD="$@"

echo "Waiting for $HOST:$PORT to be ready..."

for i in $(seq 1 $TIMEOUT); do
    if nc -z "$HOST" "$PORT" 2>/dev/null; then
        echo "$HOST:$PORT is ready!"
        exec $CMD
    fi
    echo "Retrying... ($i/$TIMEOUT)"
    sleep 1
done

echo "Timeout waiting for $HOST:$PORT"
exit 1