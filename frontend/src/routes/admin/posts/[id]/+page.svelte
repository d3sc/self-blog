<script lang="ts">
  import { page } from "$app/state";
  import { onMount } from "svelte";
  import { api } from "$lib/api";
  import type { PostInput } from "$lib/types";
  import { renderMarkdown } from "$lib/markdown";
  let edit = !!page.params.id;
  let input: PostInput = {
    title: "",
    slug: null,
    excerpt: "",
    content: "",
    status: "draft",
  };
  let error = "";
  let saved = "";
  let loading = true;
  let saving = false;
  let loadFailed = false;
  onMount(async () => {
    if (edit) {
      try {
        const p = await api.adminPost(page.params.id!);
        input = {
          title: p.title,
          slug: p.slug,
          excerpt: p.excerpt,
          content: p.content,
          status: p.status,
        };
      } catch (e) {
        loadFailed = true;
        error = e instanceof Error ? e.message : "Unable to load post";
      }
    }
    loading = false;
  });
  async function submit() {
    if (saving) return;
    saving = true;
    error = "";
    saved = "";
    try {
      if (edit) await api.update(page.params.id!, input);
      else await api.create(input);
      saved = "Saved successfully.";
      if (!edit) setTimeout(() => (location.href = "/admin"), 400);
    } catch (e) {
      error = e instanceof Error ? e.message : "Unable to save";
    } finally {
      saving = false;
    }
  }
</script>

<svelte:head
  ><title>{edit ? "Edit post" : "New post"} — Field Notes</title></svelte:head
>{#if loading}<p>Loading…</p>{:else}<div
    class="flex items-center justify-between"
  >
    <div>
      <a class="text-sm underline" href="/admin">← Posts</a>
      <h1 class="mt-3 text-4xl font-bold">{edit ? "Edit post" : "New post"}</h1>
    </div>
  </div>
  {#if error}<p class="mt-5 rounded bg-red-50 p-3 text-red-800">
      {error}
    </p>{/if}{#if saved}<p class="mt-5 rounded bg-green-50 p-3 text-green-800">
      {saved}
    </p>{/if}{#if !loadFailed}<form
      class="mt-8 grid gap-6 lg:grid-cols-2"
      on:submit|preventDefault={submit}
    >
      <div class="space-y-4">
        <label class="block"
          >Title<input
            class="mt-1 w-full rounded border p-3"
            bind:value={input.title}
            required
          /></label
        ><label class="block"
          >Slug <span class="text-sm text-stone-500">(optional)</span><input
            class="mt-1 w-full rounded border p-3"
            bind:value={input.slug}
            placeholder="generated from title"
          /></label
        ><label class="block"
          >Excerpt<textarea
            class="mt-1 w-full rounded border p-3"
            rows="3"
            bind:value={input.excerpt}></textarea></label
        ><label class="block"
          >Status<select
            class="mt-1 w-full rounded border p-3"
            bind:value={input.status}
            ><option value="draft">Draft</option><option value="published"
              >Published</option
            ></select
          ></label
        ><button
          disabled={saving}
          class="rounded bg-stone-900 px-5 py-3 text-white disabled:opacity-50"
          >{saving ? "Saving…" : edit ? "Save changes" : "Create post"}</button
        >
      </div>
      <div>
        <label class="block"
          >Markdown content<textarea
            class="mt-1 h-96 w-full rounded border p-3 font-mono text-sm"
            bind:value={input.content}></textarea></label
        >
        <h2 class="mt-4 font-semibold">Preview</h2>
        <div class="prose mt-2 rounded border bg-white p-4">
          {@html renderMarkdown(input.content)}
        </div>
      </div>
    </form>{/if}{/if}
