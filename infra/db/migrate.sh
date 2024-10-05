#!/bin/sh

echo "POSTGRES_USER = ${POSTGRES_USER} | PGPASSWORD = ${PGPASSWORD} | POSTGRES_HOST = ${POSTGRES_HOST} | POSTGRES_PORT = ${POSTGRES_PORT} | POSTGRES_DB = ${POSTGRES_DB}"
echo "postgres://${POSTGRES_USER}:${PGPASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

until migrate -path "${MIGRATIONS_PATH}" -database "postgres://${POSTGRES_USER}:${PGPASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done