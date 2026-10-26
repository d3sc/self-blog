import type { PageLoad } from "./$types";
export const load: PageLoad = async ({ fetch, url }) => {
  const raw = Number(url.searchParams.get("page") || 1);
  const page = Number.isSafeInteger(raw) && raw > 0 ? raw : 1;
  try {
    const response = await fetch(`/api/v1/posts?page=${page}&limit=10`);
    if (!response.ok)
      return {
        posts: null,
        message: "Unable to load posts. Please try again.",
        page,
      };
    return {
      posts: (await response.json()) as import("$lib/types").Paginated<
        import("$lib/types").PostSummary
      >,
      message: "",
      page,
    };
  } catch {
    return {
      posts: null,
      message: "Unable to load posts. Please try again.",
      page,
    };
  }
};
