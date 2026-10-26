# Frontend

SvelteKit 2 + Svelte 5, TypeScript, Tailwind 4 and adapter-node.

```
npm ci
npm test
npm run check
npm run build
npm start
```

`API_INTERNAL_URL` is server-only and defaults to `http://localhost:8080/api/v1`. Set it to `http://backend:8080/api/v1` in Compose. The SvelteKit hook forwards `/api/v1/*` to this fixed backend, preserving session/CSRF cookies. Browsers use same-origin `/api/v1`. `VITE_API_BASE_URL` may override browser client calls at build time; never put credentials or secrets in public/VITE variables. Prefer the same-origin default.

Public pages use universal load functions, so article bodies and SEO metadata are server-rendered. Homepage pagination is URL-backed. Admin authentication is a client UX guard only: backend authorization and CSRF remain authoritative. Mutations read `blog_csrf` and send `X-CSRF-Token`, matching the backend implementation. Markdown is sanitized through DOMPurify before HTML insertion, including preview.

Unit tests cover CSRF forwarding and hostile Markdown; build/check validate Svelte and TypeScript. Real authenticated browser CRUD remains an integration responsibility and was not exercised by this isolated frontend task.

Dependency note: Kit 2 currently brings `cookie@0.6.0` (low npm advisory GHSA-pxg6-pf52-xh8x). We do not derive cookie names/options from user input. npm audit suggests a breaking Kit 3 upgrade; deliberately not forced. Production `npm audit --omit=dev` reports zero vulnerabilities. Vitest is upgraded to 5 to avoid critical worker-pool advisories.
