import { error, redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import type { Admin } from '$lib/types';

export const load: LayoutServerLoad = async ({ fetch }) => {
  try {
    const response = await fetch('/api/v1/auth/me');
    if (response.status === 401 || response.status === 403) redirect(303, '/login');
    if (!response.ok) error(503, 'Authentication service is temporarily unavailable.');
    const body = (await response.json()) as { admin: Admin };
    return { admin: body.admin };
  } catch (e) {
    if (e && typeof e === 'object' && 'status' in e) throw e;
    error(503, 'Authentication service is temporarily unavailable.');
  }
};
