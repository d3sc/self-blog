<script lang="ts">
  import { page } from "$app/state";
</script>

<svelte:head
  ><title>{page.status} — Field Notes</title><meta
    name="robots"
    content="noindex"
  /></svelte:head
>
<div class="py-16 text-center">
  <p class="text-sm uppercase tracking-widest text-amber-700">{page.status}</p>
  <h1 class="mt-3 text-4xl font-bold">
    {page.status === 404 ? "That page wandered off." : "Something went wrong."}
  </h1>
  <p class="mt-4 text-stone-600">{page.error?.message}</p>
  <a class="mt-6 inline-block underline" href="/">Return home</a>
</div>
