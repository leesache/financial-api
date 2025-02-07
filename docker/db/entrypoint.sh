#!/bin/bash
set -e

# Start PostgreSQL in the background using the default entrypoint
/usr/local/bin/docker-entrypoint.sh postgres &

# Wait for PostgreSQL to be ready
until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
  echo "Waiting for PostgreSQL to start..."
  sleep 2
done

# Run Goose migrations
echo "Running Goose migrations..."
goose -dir /migrations postgres "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:5432/$POSTGRES_DB?sslmode=disable" up

# Keep the container running
wait