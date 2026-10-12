#!/bin/bash
set -e

echo "Starting StreamLite multi-database initialization..."

databases=(
  "streamlite_auth"
  "streamlite_catalog"
  "streamlite_transcoder"
  "streamlite_watch"
  "streamlite_recommendation"
)

for db in "${databases[@]}"; do
  echo "Checking database: $db"
  if psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -lqt | cut -d \| -f 1 | grep -qw "$db"; then
    echo "Database $db already exists, skipping."
  else
    echo "Creating database: $db"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
      CREATE DATABASE $db;
      GRANT ALL PRIVILEGES ON DATABASE $db TO $POSTGRES_USER;
EOSQL
  fi
done

echo "StreamLite multi-database initialization complete."
