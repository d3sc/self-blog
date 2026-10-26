# Personal Blog CMS MVP

## Architecture

Modular monolith with three Compose services:
- `db`: PostgreSQL with a named persistent volume; internal-only network exposure.
- `backend`: Go + Gin REST API, pgx, SQL migrations, cookie session authentication, bcrypt password hashing.
- `frontend`: SvelteKit + TypeScript + Tailwind, public blog and protected admin UI.

The browser talks to the backend through `/api/v1`; frontend server/runtime configuration supplies the API base URL. Backend owns authorization, validation, publication rules, and persistence.

## Database schema

`admins`: `id UUID PK`, `username TEXT UNIQUE NOT NULL`, `password_hash TEXT NOT NULL`, `created_at TIMESTAMPTZ NOT NULL`.

`posts`: `id UUID PK`, `title TEXT NOT NULL`, `slug TEXT UNIQUE NOT NULL`, `excerpt TEXT NOT NULL DEFAULT ''`, `content TEXT NOT NULL`, `status TEXT NOT NULL CHECK (status IN ('draft','published'))`, `created_at TIMESTAMPTZ NOT NULL`, `updated_at TIMESTAMPTZ NOT NULL`, `published_at TIMESTAMPTZ NULL`. Index `(status, published_at DESC)` and `updated_at DESC`.

Sessions may use a signed/encrypted HTTP-only cookie backed by an in-memory/session store for the single-instance MVP; production deployment should keep a stable session secret and single backend replica. CSRF is enforced for cookie-authenticated state-changing requests using a same-site cookie plus CSRF token/header.

## Ownership

- Backend agent: `backend/` and backend contract notes only.
- Frontend agent: `frontend/` only.
- DevOps/QA agent: root Compose, Dockerfiles, env/gitignore, README operations, smoke scripts; no business logic.
- Astra: shared contract, integration fixes, final security review, verification.

## Dependencies

The API contract below is stable before parallel implementation. Frontend and backend can work independently. DevOps uses the service names, ports, environment names, and migration command defined here.

## Acceptance criteria

Backend tests, frontend check/build, Compose config/build, migration/startup, auth, CRUD, draft privacy, persistence, and documented operations must be verified where the environment permits.
