#!/bin/sh
set -e

echo "Starting API Service..."

# Wait for database
if [ -n "$DB_HOST" ]; then
    echo "Waiting for database at $DB_HOST:$DB_PORT..."
    until nc -z "$DB_HOST" "$DB_PORT"; do
        echo "Waiting for database..."
        sleep 2
    done
    echo "Database is ready!"
fi

# Run database migrations if enabled
if [ "$RUN_MIGRATIONS" = "true" ]; then
    echo "Running database migrations..."
    /app/migrate up
fi

# Start the API
exec "$@"