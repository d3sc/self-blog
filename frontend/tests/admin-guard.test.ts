import { describe, it, expect, vi } from 'vitest';
import { load } from '../src/routes/admin/+layout.server';
describe('server admin guard', () => {
  it('redirects an expired session before rendering admin routes', async () => {
    await expect(load({ fetch: vi.fn().mockResolvedValue(new Response(null, { status: 401 })) } as any)).rejects.toMatchObject({ status: 303, location: '/login' });
  });
  it('does not disguise backend outages as invalid credentials', async () => {
    await expect(load({ fetch: vi.fn().mockResolvedValue(new Response(null, { status: 502 })) } as any)).rejects.toMatchObject({ status: 503 });
  });
  it('returns authenticated administrator data', async () => {
    const admin = { id: 'admin', username: 'writer' };
    expect(await load({ fetch: vi.fn().mockResolvedValue(Response.json({ admin })) } as any)).toEqual({ admin });
  });
});
