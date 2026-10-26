<script lang="ts">
  import { api } from "$lib/api";
  let username = "";
  let password = "";
  let error = "";
  let loading = false;
  async function submit() {
    loading = true;
    error = "";
    try {
      await api.login(username, password);
      location.href = "/admin";
    } catch (e) {
      error = e instanceof Error ? e.message : "Login failed";
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head><title>Log in — Field Notes</title></svelte:head>
<div class="mx-auto max-w-md">
  <h1 class="text-3xl font-bold">Admin login</h1>
  <p class="mt-2 text-stone-600">Sign in to manage your writing.</p>
  {#if error}<p class="mt-4 rounded bg-red-50 p-3 text-red-800">{error}</p>{/if}
  <form class="mt-8 space-y-4" on:submit|preventDefault={submit}>
    <label class="block"
      >Username<input
        class="mt-1 w-full rounded border p-3"
        bind:value={username}
        autocomplete="username"
        required
      /></label
    ><label class="block"
      >Password<input
        class="mt-1 w-full rounded border p-3"
        type="password"
        bind:value={password}
        autocomplete="current-password"
        required
      /></label
    ><button
      class="rounded bg-stone-900 px-5 py-3 text-white disabled:opacity-50"
      disabled={loading}>{loading ? "Signing in…" : "Sign in"}</button
    >
  </form>
</div>
