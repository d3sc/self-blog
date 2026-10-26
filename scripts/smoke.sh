#!/usr/bin/env sh
set -eu

base_url=${1:-http://127.0.0.1:3000}
base_url=${base_url%/}

check_status() {
  url=$1
  expected=$2
  actual=$(curl -sS -o /tmp/blog-cms-smoke-body -w '%{http_code}' "$url")
  if [ "$actual" != "$expected" ]; then
    printf 'FAIL %s: expected HTTP %s, got %s\n' "$url" "$expected" "$actual" >&2
    sed -n '1,8p' /tmp/blog-cms-smoke-body >&2 || true
    exit 1
  fi
  printf 'PASS %s (HTTP %s)\n' "$url" "$actual"
}

check_status "$base_url/" 200
check_status "$base_url/api/v1/posts?page=1&limit=10" 200
printf 'Smoke checks passed (read-only).\n'
