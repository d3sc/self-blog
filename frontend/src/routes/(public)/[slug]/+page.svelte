<script lang="ts">
  import { renderMarkdown } from "$lib/markdown";
  import type { PageData } from "./$types";
  export let data: PageData;
  $: post = data.post;
  let error = "";
</script>

<svelte:head
  ><title>{post?.title ?? "Article"} — Field Notes</title><meta
    name="description"
    content={post?.excerpt ?? "Read an article on Field Notes."}
  /></svelte:head
>{#if post}<article>
    <a class="text-sm text-stone-500 underline" href="/">← All notes</a>
    <p class="mt-10 text-sm text-stone-500">
      {new Date(post.publishedAt ?? post.createdAt).toLocaleDateString()}
    </p>
    <h1 class="mt-3 max-w-3xl text-5xl font-bold tracking-tight">
      {post.title}
    </h1>
    <p class="mt-5 text-xl text-stone-600">{post.excerpt}</p>
    <div class="prose mt-12">{@html renderMarkdown(post.content)}</div>
  </article>{:else if error}<h1 class="text-3xl font-bold">
    Article not found
  </h1>
  <p class="mt-3">This article may have moved or is not published.</p>
  <a class="mt-6 inline-block underline" href="/">Return home</a>{:else}<p>
    Loading article…
  </p>{/if}
