# Personal Blog CMS MVP

See [ARCHITECTURE.md](ARCHITECTURE.md) and [API_CONTRACT.md](API_CONTRACT.md) for the service and API contract.

## Local and production-shaped Compose run

Requirements: Docker Engine with Compose v2. Copy the example environment and replace both secret placeholders:

```sh
cp .env.example .env
openssl rand -hex 32  # use output for POSTGRES_PASSWORD and a separate output for SESSION_SECRET
# edit .env; never commit it

docker compose config
docker compose up -d --build
docker compose ps
```

The frontend is bound to `127.0.0.1:3000` by default. PostgreSQL has no host port and is on an internal-only network. The backend is reachable by the frontend over the private Compose network, and the browser uses the frontend's `/api/v1` path. Put a TLS reverse proxy in front of the frontend for a public deployment; set `TRUSTED_ORIGINS` and `COOKIE_SECURE` to match the HTTPS origin.

The backend waits for a healthy database, runs `/app/server migrate` on startup, and then runs `/app/server`. The backend image also provides `/app/create-admin`; it prompts for a password without echoing it. The health checks expect `/healthz` and `/` respectively. Keep one backend replica for the MVP because sessions may be in memory.

### Start, inspect, and create the first admin

```sh
docker compose up -d --build
docker compose ps
docker compose logs --tail=100 backend frontend db
# Run migrations explicitly (safe when there are no changes):
docker compose exec backend /app/server migrate
# Interactive only; type the password at both hidden prompts:
docker compose exec backend /app/create-admin --username admin
```

For local HTTP, keep `COOKIE_SECURE=false` (the `.env.example` default); otherwise browsers will not send the secure cookies. For production, put a TLS reverse proxy in front of the frontend, set `COOKIE_SECURE=true`, and set `TRUSTED_ORIGINS` to the exact HTTPS origin(s). Never use the local HTTP setting on a public deployment.

## Configuration and secret handling

`.env.example` documents all Compose inputs. `.env` is ignored by Git. Use a secret manager or protected deployment environment rather than placing production secrets in source control or image layers. Use distinct, long random values for the database password and session secret; rotate both according to your operational policy. Do not use the example values in production.

## Backup and restore

Back up the database from the running `db` service; no database host port is required. Ensure the backup destination is protected and encrypted:

```sh
mkdir -p backups
STAMP=$(date -u +%Y%m%dT%H%M%SZ)
docker compose exec -T db pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" --format=custom > "backups/blog-$STAMP.dump"
```

Test restores into a separate, disposable database or isolated Compose project before relying on a backup. A restore overwrites its target data; the following drill creates a new temporary Postgres container and never connects to the live `db` service:

```sh
RESTORE_PROJECT="blog-restore-$$"
RESTORE_DB="restore"
RESTORE_USER="restore"
RESTORE_PASS="$(openssl rand -hex 24)"
docker network create "${RESTORE_PROJECT}_default"
docker run -d --rm --name "$RESTORE_PROJECT-db" --network "${RESTORE_PROJECT}_default" \
  -e POSTGRES_DB="$RESTORE_DB" -e POSTGRES_USER="$RESTORE_USER" -e POSTGRES_PASSWORD="$RESTORE_PASS" \
  postgres:16-alpine
until docker exec "$RESTORE_PROJECT-db" pg_isready -U "$RESTORE_USER" -d "$RESTORE_DB" >/dev/null 2>&1; do sleep 2; done
docker run --rm --network "${RESTORE_PROJECT}_default" \
  -v "$PWD/backups/blog-YYYYMMDDTHHMMSSZ.dump:/tmp/backup.dump:ro" \
  -e PGPASSWORD="$RESTORE_PASS" -e RESTORE_HOST="$RESTORE_PROJECT-db" \
  -e RESTORE_USER="$RESTORE_USER" -e RESTORE_DB="$RESTORE_DB" postgres:16-alpine sh -ec \
  'pg_restore --clean --if-exists --no-owner -h "$RESTORE_HOST" -U "$RESTORE_USER" -d "$RESTORE_DB" /tmp/backup.dump'
docker rm -f "$RESTORE_PROJECT-db" >/dev/null
docker network rm "${RESTORE_PROJECT}_default" >/dev/null
```

Replace the dump filename with the actual backup. If the restore command fails, remove the temporary container and network before retrying; do not point `pg_restore` at the live Compose database.

## Update and rollback

1. Review the image/Dockerfile and migration changes; back up the database.
2. Pull or build the reviewed revision, then run `docker compose config`.
3. Start with `docker compose up -d --build`; Compose recreates only changed services.
4. Run `scripts/smoke.sh http://127.0.0.1:3000` and inspect `docker compose ps`/logs.
5. Roll back application code to the prior reviewed revision and run `docker compose up -d --build`. Treat irreversible migrations as a separately reviewed release; restore the backup only through an approved maintenance procedure.

Do not use `docker compose down -v` in routine operations: it deletes the named database volume. Stop without deleting data with `docker compose stop`; remove containers while retaining data with `docker compose down`.

## Non-destructive verification checklist

- [ ] `docker compose --env-file .env config` succeeds and shows `db`, `backend`, and `frontend`.
- [ ] `docker compose up -d --build` completes; `docker compose ps` shows all services healthy.
- [ ] `scripts/smoke.sh` returns success for the frontend URL.
- [ ] `GET /api/v1/posts?page=1&limit=10` returns the contract envelope (usually an empty published list on a fresh DB).
- [ ] Login, CSRF, admin CRUD, publish/unpublish, and logout are exercised using test credentials only.
- [ ] A draft is absent from both public list and slug endpoints.
- [ ] Restarting services preserves data: `docker compose restart`, then re-run checks.
- [ ] A `pg_dump` backup is created and a restore drill succeeds in an isolated target.
- [ ] `scripts/integration-qa.sh` passes with random credentials and removes its project/volume on exit.
- [ ] No credentials, `.env`, or database ports are exposed in Git or the host port list.

## QA scripts

`scripts/verify-compose.sh` validates Compose interpolation and service structure without starting or deleting anything. `scripts/smoke.sh [base-url]` performs read-only frontend and API checks with `curl`. `scripts/integration-qa.sh` starts an isolated random Compose project, creates a random admin through `/app/create-admin`, exercises authentication, CSRF, draft privacy, publish/unpublish, logout, and restart persistence, then removes only its own project and volume. Set `QA_HOST_PORT` if port 3000 is occupied. It never prints credentials.
