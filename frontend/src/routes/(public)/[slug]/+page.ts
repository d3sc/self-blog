import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import type { Post } from "$lib/types";
export const load: PageLoad = async ({ fetch, params }) => {
  const response = await fetch(
    `/api/v1/posts/${encodeURIComponent(params.slug)}`,
  );
  if (response.status === 404)
    error(404, "This article does not exist or is not published.");
  if (!response.ok)
    error(
      response.status >= 500 ? 503 : response.status,
      "Unable to load this article. Please try again.",
    );
  return { post: (await response.json()) as Post };
};
