FROM migrate/migrate:v4.15.1
COPY ./migrations /migrations
ENTRYPOINT [ "sh", "-c", "migrate -path /migrations -database postgres://${POSTGRES_USER}:${PGPASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}?sslmode=disable up"]