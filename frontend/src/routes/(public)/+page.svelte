<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import type { PageData } from "./$types";
  export let data: PageData;
  $: posts = data.posts;
  $: error = data.message;
  $: page = data.page;
  $: totalPages = posts?.pagination.totalPages ?? 0;
</script>

<svelte:head
  ><title>Field Notes — Personal blog</title><meta
    name="description"
    content="A personal blog of thoughts and experiments."
  /></svelte:head
>
<section class="mb-12 max-w-2xl">
  <p
    class="mb-3 text-sm font-semibold uppercase tracking-widest text-amber-700"
  >
    Personal blog
  </p>
  <h1 class="text-5xl font-bold tracking-tight">Ideas worth keeping.</h1>
  <p class="mt-5 text-lg text-stone-600">
    Notes on building, learning, and making sense of the world.
  </p>
</section>
{#if !posts && !error}<p aria-live="polite">
    Loading posts…
  </p>{:else if error}<div
    class="rounded border border-red-200 bg-red-50 p-4 text-red-800"
  >
    {error}<button class="ml-3 underline" on:click={() => invalidateAll()}
      >Retry</button
    >
  </div>{:else if posts?.items.length === 0}<p class="text-stone-600">
    No published posts yet.
  </p>{:else}<div class="grid gap-8 md:grid-cols-2">
    {#each posts?.items ?? [] as post}<article
        class="border-t-2 border-stone-900 pt-4"
      >
        <p class="text-sm text-stone-500">
          {new Date(post.publishedAt ?? post.createdAt).toLocaleDateString()}
        </p>
        <h2 class="mt-2 text-2xl font-semibold">
          <a class="hover:text-amber-700" href={"/" + post.slug}>{post.title}</a
          >
        </h2>
        <p class="mt-2 text-stone-600">{post.excerpt}</p>
        <a
          class="mt-4 inline-block text-sm font-semibold underline"
          href={"/" + post.slug}>Read article →</a
        >
      </article>{/each}
  </div>
  {#if totalPages > 1}<div class="mt-12 flex items-center gap-4">
      <a
        aria-disabled={page <= 1}
        href={page > 1 ? `/?page=${page - 1}` : undefined}>← Newer</a
      ><span>Page {page} of {totalPages}</span><a
        aria-disabled={page >= totalPages}
        href={page < totalPages ? `/?page=${page + 1}` : undefined}>Older →</a
      >
    </div>{/if}{/if}
