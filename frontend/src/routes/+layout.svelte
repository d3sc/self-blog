<script lang="ts">
  import "../app.css";
  import { onMount } from "svelte";
  import { api } from "$lib/api";
  import type { Admin } from "$lib/types";
  let admin: Admin | null = null;
  onMount(() =>
    api
      .me()
      .then((r) => (admin = r.admin))
      .catch(() => {}),
  );
  async function logout() {
    await api.logout();
    location.href = "/";
  }
</script>

<header class="border-b bg-white">
  <nav class="mx-auto flex max-w-5xl items-center justify-between px-5 py-4">
    <a href="/" class="text-xl font-bold tracking-tight">Field Notes</a>
    <div class="flex gap-4 text-sm">
      <a href="/">Journal</a><a href="/about">About</a>{#if admin}<a
          href="/admin">Admin</a
        ><button on:click={logout}>Log out</button>{:else}<a href="/login"
          >Log in</a
        >{/if}
    </div>
  </nav>
</header>
<main class="mx-auto min-h-[70vh] max-w-5xl px-5 py-10"><slot /></main>
<footer class="border-t py-8 text-center text-sm text-stone-500">
  Thoughts, experiments, and useful things.
</footer>
