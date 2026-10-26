#!/usr/bin/env sh
set -eu

# Use the example only for structural validation; required secret interpolation is
# intentionally supplied as ephemeral process environment, never written to disk.
export POSTGRES_PASSWORD='compose-config-check-db-password'
export SESSION_SECRET='compose-config-check-session-secret-01234567890123456789'

config=$(docker compose --env-file .env.example config 2>&1)
printf '%s\n' "$config" | grep -q '^services:'
printf '%s\n' "$config" | grep -q '^  db:'
printf '%s\n' "$config" | grep -q '^  backend:'
printf '%s\n' "$config" | grep -q '^  frontend:'
printf '%s\n' "$config" | grep -q 'internal: true'
printf '%s\n' "$config" | grep -q 'blog_cms_postgres_data'
printf '%s\n' "$config" | grep -q 'MIGRATIONS_PATH: file:///app/migrations'
printf 'Compose structure and interpolation checks passed.\n'
