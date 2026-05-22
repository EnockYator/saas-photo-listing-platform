#!/bin/sh
set -e

echo "Starting Migration Service..."

# Wait for database
if [ -n "$DB_HOST" ]; then
    echo "Waiting for database at $DB_HOST:$DB_PORT..."
    until nc -z "$DB_HOST" "$DB_PORT"; do
        echo "Waiting for database..."
        sleep 2
    done
    echo "Database is ready!"
fi

# Run migration command
echo "Running migration command: $@"
exec /app/migrate "$@"