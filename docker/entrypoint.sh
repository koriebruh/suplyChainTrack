#!/bin/sh
set -e

echo ">>> Menunggu PostgreSQL siap di $DB_HOST:$DB_PORT..."
until nc -z "$DB_HOST" "$DB_PORT"; do
  echo "Masih nunggu DB..."
  sleep 5
done

echo ">>> DB ready, menjalankan migrasi..."
migrate -database "$DATABASE_URL" -path db/migrations up

echo ">>> Menjalankan aplikasi..."
mkdir -p logger
touch logger/application.log
exec ./app
