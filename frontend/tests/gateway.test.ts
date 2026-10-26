import { afterEach, describe, expect, it, vi } from 'vitest';
vi.mock('$env/dynamic/private', () => ({ env: { API_INTERNAL_URL: 'http://backend:8080/api/v1' } }));
import { handle } from '../src/hooks.server';
afterEach(() => vi.unstubAllGlobals());
function event(path = '/api/v1/auth/logout', options: RequestInit = {}) {
  const url = new URL(path, 'http://blog.test');
  return { url, request: new Request(url, options) } as any;
}
const resolve = vi.fn();
describe('same-origin gateway', () => {
  it('rejects foreign-origin login before contacting the backend', async () => {
    const fetch = vi.fn(); vi.stubGlobal('fetch', fetch);
    const response = await handle({ event: event('/api/v1/auth/login', { method: 'POST', headers: { origin: 'https://attacker.test' } }), resolve } as any);
    expect(response.status).toBe(403);
    expect(fetch).not.toHaveBeenCalled();
  });
  it('allows same-origin login even when the browser reports cross-site fetch metadata', async () => {
    const fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify({ admin: { id: '1', username: 'admin' } }), { status: 200 }));
    vi.stubGlobal('fetch', fetch);
    const response = await handle({ event: event('/api/v1/auth/login', { method: 'POST', headers: { origin: 'http://blog.test', 'sec-fetch-site': 'cross-site' } }), resolve } as any);
    expect(response.status).toBe(200);
    expect(fetch).toHaveBeenCalledOnce();
  });
  it('forwards session and CSRF and preserves separate Set-Cookie headers', async () => {
    const headers = new Headers();
    headers.append('set-cookie', 'blog_session=session; Path=/; HttpOnly');
    headers.append('set-cookie', 'blog_csrf=token; Path=/');
    const fetch = vi.fn().mockResolvedValue(new Response(null, { status: 204, headers }));
    vi.stubGlobal('fetch', fetch);
    const response = await handle({ event: event(undefined, { method: 'POST', headers: { origin: 'http://blog.test', cookie: 'blog_session=session; blog_csrf=token', 'X-CSRF-Token': 'token' } }), resolve } as any);
    expect(fetch.mock.calls[0][0]).toBe('http://backend:8080/api/v1/auth/logout');
    expect(fetch.mock.calls[0][1].headers.get('cookie')).toContain('blog_session=session');
    expect(fetch.mock.calls[0][1].headers.get('X-CSRF-Token')).toBe('token');
    expect(response.headers.getSetCookie()).toHaveLength(2);
  });
  it('returns a controlled unavailable response', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('ECONNREFUSED')));
    const response = await handle({ event: event('/api/v1/posts'), resolve } as any);
    expect(response.status).toBe(502);
    expect((await response.json()).error.code).toBe('upstream_unavailable');
  });
});
