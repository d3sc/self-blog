# API Contract v1

Base path: `/api/v1`. JSON uses camelCase. Errors use `{ "error": { "code": "...", "message": "...", "fields": { "field": "..." } } }`; `fields` is optional. IDs are UUID strings and timestamps are RFC3339.

## Authentication

`POST /auth/login`
Request: `{ "username": string, "password": string }`
Success 200: `{ "admin": { "id": string, "username": string } }`; sets secure HTTP-only SameSite=Lax session cookie and CSRF cookie/token as applicable.
401 invalid credentials; 429 rate limited; 400 validation.

`POST /auth/logout`
Auth required. Success 204; clears session.

`GET /auth/me`
Auth optional. Authenticated success 200: `{ "admin": { "id": string, "username": string } }`; unauthenticated 401.

## Public posts

`GET /posts?page=1&limit=10`
Returns published posts only, newest publication first. Success 200: `{ "items": [PostSummary], "pagination": { "page": 1, "limit": 10, "total": 0, "totalPages": 0 } }`.

`GET /posts/:slug`
Returns a published post only. Success 200: `Post`; unpublished/missing both return 404.

`PostSummary`: `{ "id": string, "title": string, "slug": string, "excerpt": string, "status": "published", "createdAt": string, "updatedAt": string, "publishedAt": string|null }`.
`Post`: summary fields plus `content: string`.

## Admin posts

All admin endpoints require an authenticated administrator and CSRF protection for POST/PUT/DELETE.

`GET /admin/posts?page=1&limit=20&status=all|draft|published`
Success 200: `{ "items": [Post], "pagination": ... }`.

`POST /admin/posts`
Request: `{ "title": string, "slug": string|null, "excerpt": string, "content": string, "status": "draft"|"published" }`. Empty/null slug is generated from title. Success 201: `Post`. 409 slug conflict; 400 validation.

`GET /admin/posts/:id`
Success 200: `Post`; 404 missing.

`PUT /admin/posts/:id`
Same request shape as create; slug may be manually edited. Publishing sets `publishedAt` if absent; unpublishing clears it. Success 200: `Post`; 409 slug conflict; 400 validation; 404 missing.

`DELETE /admin/posts/:id`
Success 204; 404 missing.

## Rules

- Public endpoints never return drafts or draft content; missing and unpublished are intentionally indistinguishable.
- `status` is the source of publication truth. `publishedAt` is set on first publication and retained while published; it is cleared on unpublish only if the backend chooses a simpler policy, but response must be consistent.
- Pagination limits are bounded server-side (recommended max 100). Invalid values return 400 or are normalized consistently.
- Parameterized SQL is required. Backend authorization is authoritative; frontend route guards are UX only.
