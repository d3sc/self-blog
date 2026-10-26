<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api";
  import type { Post, Paginated } from "$lib/types";
  let data: Paginated<Post> | null = null;
  let status: "all" | "draft" | "published" = "all";
  let error = "";
  let page = 1;
  let success = "";
  async function load() {
    try {
      error = "";
      data = await api.adminPosts(page, 20, status);
    } catch (e) {
      error = e instanceof Error ? e.message : "Unable to load posts";
    }
  }
  onMount(load);
  async function remove(post: Post) {
    if (confirm(`Delete “${post.title}”? This cannot be undone.`)) {
      try {
        await api.remove(post.id);
        success = "Post deleted.";
        await load();
      } catch (e) {
        error = e instanceof Error ? e.message : "Delete failed";
      }
    }
  }
</script>

<svelte:head><title>Admin — Field Notes</title></svelte:head>
<div class="flex flex-wrap items-center justify-between gap-4">
  <div>
    <p class="text-sm uppercase tracking-widest text-amber-700">Admin</p>
    <h1 class="text-4xl font-bold">Your posts</h1>
  </div>
  <a class="rounded bg-stone-900 px-4 py-2 text-white" style="color: white;" href="/admin/posts/new"
    >New post</a
  >
</div>
<div class="mt-8 flex gap-2">
  {#each ["all", "draft", "published"] as option}<button
      class="rounded border px-3 py-1 text-sm {status === option
        ? 'bg-stone-900 text-white'
        : ''}"
      on:click={() => {
        status = option as typeof status;
        page = 1;
        load();
      }}>{option}</button
    >{/each}
</div>
{#if error}<p class="mt-5 rounded bg-red-50 p-3 text-red-800">
    {error}
  </p>{/if}{#if !data && !error}<p class="mt-8">
    Loading posts…
  </p>{:else if data?.items.length === 0}<p class="mt-8 text-stone-600">
    No posts found.
  </p>{:else}<div class="mt-6 divide-y border-y">
    {#each data?.items ?? [] as post}<div
        class="flex flex-wrap items-center justify-between gap-3 py-4"
      >
        <div>
          <a class="font-semibold underline" href={"/admin/posts/" + post.id}
            >{post.title}</a
          >
          <p class="text-sm text-stone-500">
            {post.status} · updated {new Date(
              post.updatedAt,
            ).toLocaleDateString()}
          </p>
        </div>
        <button
          class="text-sm text-red-700 underline"
          on:click={() => remove(post)}>Delete</button
        >
      </div>{/each}
  </div>{/if}{#if success}<p class="mt-4 text-green-800" role="status">
    {success}
  </p>{/if}{#if data && data.pagination.totalPages > 1}<nav
    aria-label="Post pages"
    class="mt-6 flex gap-5"
  >
    <button
      disabled={page <= 1}
      on:click={() => {
        page--;
        load();
      }}>Previous</button
    ><span>Page {page} of {data.pagination.totalPages}</span><button
      disabled={page >= data.pagination.totalPages}
      on:click={() => {
        page++;
        load();
      }}>Next</button
    >
  </nav>{/if}
