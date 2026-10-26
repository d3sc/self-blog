# Blog CMS backend

## Configuration

Required: `DATABASE_URL`, `SESSION_SECRET` (at least 32 random bytes, supplied as base64 or plain text). Optional: `PORT` (default `8080`), `COOKIE_SECURE` (default `true`; set `false` only for local HTTP), `TRUSTED_ORIGINS` (comma-separated origins, used for CORS), `MIGRATIONS_PATH` (default `file://migrations`).

## Run

```sh
go run ./cmd/server
```

The server applies SQL migrations on startup. It listens on `:8080` by default.

## First admin setup

There is no default or hardcoded credential. Before starting the server, set `DATABASE_URL`, then run:

```sh
go run ./cmd/create-admin --username admin
```

The command prompts for the password twice without printing it, validates the username and password, hashes with bcrypt, and inserts the first admin. It refuses duplicate usernames. Store credentials in a secret manager; never commit them or pass passwords in shell history.

## Tests

`go test ./...` runs unit tests (a PostgreSQL integration suite can be added by the deployment harness using `DATABASE_URL`).
