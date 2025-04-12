#!/bin/bash

# Script to run SQL migrations for the TODO List application

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
fi

# Check if PGPASSWORD is set
if [ -z "$DB_PASSWORD" ]; then
  echo "Error: DB_PASSWORD environment variable is not set."
  exit 1
fi

# Set PostgreSQL password
export PGPASSWORD=$DB_PASSWORD

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-todo_app}

# Function to run a migration file
run_migration() {
  local file=$1
  echo "Running migration: $file"
  
  psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $file
  
  if [ $? -eq 0 ]; then
    echo "Migration successful: $file"
  else
    echo "Migration failed: $file"
    exit 1
  fi
}

# Create database if it doesn't exist
echo "Checking if database exists..."
DB_EXISTS=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -w $DB_NAME | wc -l)

if [ $DB_EXISTS -eq 0 ]; then
  echo "Creating database: $DB_NAME"
  psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
else
  echo "Database $DB_NAME already exists."
fi

# Run all migration files in order
echo "Running migrations..."
for file in migrations/*.sql; do
  run_migration $file
done

echo "All migrations completed successfully."
