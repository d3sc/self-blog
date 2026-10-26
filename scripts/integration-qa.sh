#!/usr/bin/env sh
# Disposable end-to-end Compose QA. It owns only a random Compose project/volume.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PROJECT="blog-qa-$$-${RANDOM:-0}"
TMP=$(mktemp -d "${TMPDIR:-/tmp}/blog-cms-qa.XXXXXX")
ENV_FILE="$TMP/.env"
DOCKER_ENV_FILE=$(cygpath -w "$ENV_FILE" 2>/dev/null || printf '%s' "$ENV_FILE")
COOKIE_JAR="$TMP/cookies"
CURL_COOKIE_JAR=$(cygpath -w "$COOKIE_JAR" 2>/dev/null || printf '%s' "$COOKIE_JAR")
CURL_RESPONSE=$(cygpath -w "$TMP/response" 2>/dev/null || printf '%s' "$TMP/response")
CURL_HEADERS=$(cygpath -w "$TMP/headers" 2>/dev/null || printf '%s' "$TMP/headers")
BASE_URL="${QA_BASE_URL:-http://127.0.0.1:${QA_HOST_PORT:-3000}}"
BASE_URL=${BASE_URL%/}
DB_PASS=$(openssl rand -hex 24)
SESSION=$(openssl rand -hex 32)
ADMIN_USER="qa_$(openssl rand -hex 6)"
ADMIN_PASS=$(openssl rand -hex 18)
SLUG="qa-$(openssl rand -hex 6)"

cleanup() {
  set +e
  docker compose -p "$PROJECT" --env-file "$ENV_FILE" down -v --remove-orphans >/dev/null 2>&1
  rm -rf "$TMP"
}
trap cleanup EXIT INT TERM

cat >"$ENV_FILE" <<EOF
APP_ENV=development
POSTGRES_DB=blog
POSTGRES_USER=blog
POSTGRES_PASSWORD=$DB_PASS
SESSION_SECRET=$SESSION
COOKIE_SECURE=false
TRUSTED_ORIGINS=http://127.0.0.1:3000,http://localhost:3000
FRONTEND_BIND=127.0.0.1
FRONTEND_HOST_PORT=${QA_HOST_PORT:-3000}
FRONTEND_PORT=3000
BACKEND_PORT=8080
API_INTERNAL_URL=http://backend:8080/api/v1
POSTGRES_VOLUME_NAME=${PROJECT}_postgres_data
EOF

compose() { docker compose -p "$PROJECT" --env-file "$DOCKER_ENV_FILE" "$@"; }
request() { curl -sS -D "$CURL_HEADERS" -o "$CURL_RESPONSE" -w '%{http_code}' -b "$CURL_COOKIE_JAR" -c "$CURL_COOKIE_JAR" "$@"; }
expect() {
  actual=$1; expected=$2; label=$3
  [ "$actual" = "$expected" ] || { printf 'FAIL %s: expected %s got %s\n' "$label" "$expected" "$actual" >&2; exit 1; }
  printf 'PASS %s (HTTP %s)\n' "$label" "$actual"
}

printf 'Starting disposable Compose project %s\n' "$PROJECT"
compose up -d --build
READY_BODY="$TMP/ready"
for n in 1 2 3 4 5 6 7 8 9 10 11 12; do
  code=$(curl -sS -o "$READY_BODY" -w '%{http_code}' "$BASE_URL/" || true)
  [ "$code" = 200 ] && break
  sleep 2
  [ "$n" = 12 ] && { printf 'FAIL frontend did not become ready\n' >&2; exit 1; }
done

# The admin helper is intentionally interactive; credentials are piped, never printed.
printf '%s\n%s\n' "$ADMIN_PASS" "$ADMIN_PASS" | compose exec -T backend /app/create-admin --username "$ADMIN_USER" >/dev/null

code=$(request "$BASE_URL/api/v1/posts?page=1&limit=10"); expect "$code" 200 'public list starts empty'
code=$(request -X POST -H 'Content-Type: application/json' -d '{"title":"x","content":"x","status":"draft"}' "$BASE_URL/api/v1/admin/posts"); expect "$code" 401 'admin endpoint rejects anonymous'
code=$(request -X POST -H 'Content-Type: application/json' -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASS\"}" "$BASE_URL/api/v1/auth/login"); expect "$code" 200 'login'
csrf=$(awk -F'[=;]' 'tolower($1) ~ /^set-cookie: blog_csrf/ { print $2; exit }' "$CURL_HEADERS" | tr -d '\r')
[ -n "$csrf" ] || { printf 'FAIL CSRF cookie missing\n' >&2; sed -n '1,80p' "$CURL_HEADERS" >&2; exit 1; }
code=$(request -X POST -H 'Content-Type: application/json' -d '{"title":"CSRF","content":"blocked","status":"draft"}' "$BASE_URL/api/v1/admin/posts"); expect "$code" 403 'CSRF required'
body=$(printf '{"title":"QA draft","slug":"%s","excerpt":"private","content":"private body","status":"draft"}' "$SLUG")
code=$(request -X POST -H 'Content-Type: application/json' -H "X-CSRF-Token: $csrf" -d "$body" "$BASE_URL/api/v1/admin/posts"); expect "$code" 201 'create draft'
id=$(sed -n 's/.*"id":"\([^"]*\)".*/\1/p' "$CURL_RESPONSE")
[ -n "$id" ] || { printf 'FAIL create response did not contain an id\n' >&2; exit 1; }
code=$(request "$BASE_URL/api/v1/posts/$SLUG"); expect "$code" 404 'draft hidden by slug'
code=$(request "$BASE_URL/api/v1/admin/posts?status=draft"); expect "$code" 200 'draft visible to admin'
code=$(request -X PUT -H 'Content-Type: application/json' -H "X-CSRF-Token: $csrf" -d "{\"title\":\"QA published\",\"slug\":\"$SLUG\",\"excerpt\":\"public\",\"content\":\"public body\",\"status\":\"published\"}" "$BASE_URL/api/v1/admin/posts/$id"); expect "$code" 200 'publish draft'
code=$(request "$BASE_URL/api/v1/posts/$SLUG"); expect "$code" 200 'published post public'
code=$(request -X PUT -H 'Content-Type: application/json' -H "X-CSRF-Token: $csrf" -d "{\"title\":\"QA draft again\",\"slug\":\"$SLUG\",\"excerpt\":\"private\",\"content\":\"private body\",\"status\":\"draft\"}" "$BASE_URL/api/v1/admin/posts/$id"); expect "$code" 200 'unpublish'
code=$(request "$BASE_URL/api/v1/posts/$SLUG"); expect "$code" 404 'unpublished post hidden'
code=$(request -X POST -H "X-CSRF-Token: $csrf" "$BASE_URL/api/v1/auth/logout"); expect "$code" 204 'logout'
code=$(request "$BASE_URL/api/v1/admin/posts"); expect "$code" 401 'logout invalidates session'
compose restart >/dev/null
for n in 1 2 3 4 5 6 7 8 9 10 11 12; do
  code=$(curl -sS -o "$CURL_RESPONSE" -w '%{http_code}' "$BASE_URL/api/v1/posts?page=1&limit=10" || true)
  [ "$code" = 200 ] && break
  sleep 2
  [ "$n" = 12 ] && { printf 'FAIL public API unavailable after restart\n' >&2; exit 1; }
done
expect "$code" 200 'public API after restart'
printf 'Integration QA passed; disposable project and volume removed.\n'
