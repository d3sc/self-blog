export type Status = "draft" | "published";
export interface Admin {
  id: string;
  username: string;
}
export interface PostSummary {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  status: "published";
  createdAt: string;
  updatedAt: string;
  publishedAt: string | null;
}
export interface Post extends Omit<PostSummary, "status"> {
  status: Status;
  content: string;
}
export interface Pagination {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}
export interface Paginated<T> {
  items: T[];
  pagination: Pagination;
}
export interface ApiError {
  error: { code: string; message: string; fields?: Record<string, string> };
}
export interface PostInput {
  title: string;
  slug: string | null;
  excerpt: string;
  content: string;
  status: Status;
}
