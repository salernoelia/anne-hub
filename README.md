# anne Hub

Routing and processing server for anne wear and anne companion.

# Setup the PostgreSQL DB and Migrations

Make sure you have postgres 14 installed and running.

```sh
brew install postgres@14
```

```sh
brew services start postgresql@14
```

Get the CLI tool for migrations (with brew on Linux or macOS, otherwise you can get it from the releases of the [official repository](https://github.com/golang-migrate/migrate))

```sh
brew install golang-migrate
```

Create the Database

```sh
psql -U <username> -tc "SELECT 1 FROM pg_database WHERE datname = 'anne_hub';" | grep -q 1 || psql -U <username> -c "CREATE DATABASE anne_hub;"
```

Create a migration

```sh
migrate create -ext sql -dir db/migrations -seq migration_name
```

Apply a migration

```sh
migrate -database $ANNE_HUB_DB -path db/migrations up
```

# env example

```env
GROQ_API_KEY=xxxxx
DB_HOST=host.docker.internal
DB_PORT=5432
DB_USERNAME=xxxxx
DB_NAME=anne_hub
DB_PASSWORD=xxxxx
DB_SSLMODE=disable
```

# Quickstart with Docker

For building:

```sh
sudo docker build -t anne-hub .
```

To run:

```sh
sudo docker run --env-file .env -p 1323:1323 anne-hub
```

# Quickstart with go cli

```sh
go mod tidy
```

```sh
go build -o bin/anne-hub-v*
```
