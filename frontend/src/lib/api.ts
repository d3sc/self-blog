import type {
  Admin,
  ApiError,
  Paginated,
  Post,
  PostInput,
  PostSummary,
  Status,
} from "./types";
const base = import.meta.env.VITE_API_BASE_URL || "/api/v1";
async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.body) headers.set("content-type", "application/json");
  if (
    typeof document !== "undefined" &&
    init.method &&
    !["GET", "HEAD"].includes(init.method)
  ) {
    const token = document.cookie
      .split("; ")
      .find((c) => c.startsWith("blog_csrf="))
      ?.slice("blog_csrf=".length);
    if (token) headers.set("X-CSRF-Token", decodeURIComponent(token));
  }
  const res = await fetch(`${base}${path}`, {
    ...init,
    headers,
    credentials: "include",
  });
  if (!res.ok) {
    let detail: ApiError | undefined;
    try {
      detail = await res.json();
    } catch {}
    throw new Error(detail?.error?.message || `Request failed (${res.status})`);
  }
  return res.status === 204 ? (undefined as T) : ((await res.json()) as T);
}
export const api = {
  login: (username: string, password: string) =>
    request<{ admin: Admin }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    }),
  logout: () => request<void>("/auth/logout", { method: "POST" }),
  me: () => request<{ admin: Admin }>("/auth/me"),
  publicPosts: (page = 1, limit = 10) =>
    request<Paginated<PostSummary>>(`/posts?page=${page}&limit=${limit}`),
  post: (slug: string) => request<Post>(`/posts/${encodeURIComponent(slug)}`),
  adminPosts: (page = 1, limit = 20, status: "all" | Status = "all") =>
    request<Paginated<Post>>(
      `/admin/posts?page=${page}&limit=${limit}&status=${status}`,
    ),
  adminPost: (id: string) =>
    request<Post>(`/admin/posts/${encodeURIComponent(id)}`),
  create: (input: PostInput) =>
    request<Post>("/admin/posts", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  update: (id: string, input: PostInput) =>
    request<Post>(`/admin/posts/${encodeURIComponent(id)}`, {
      method: "PUT",
      body: JSON.stringify(input),
    }),
  remove: (id: string) =>
    request<void>(`/admin/posts/${encodeURIComponent(id)}`, {
      method: "DELETE",
    }),
};
