<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api";
  let authorized = false;
  onMount(async () => {
    try {
      await api.me();
      authorized = true;
    } catch {
      location.href = "/login";
    }
  });
</script>

<svelte:head><meta name="robots" content="noindex,nofollow" /></svelte:head>
{#if authorized}<slot />{:else}<p role="status">Checking your session…</p>{/if}
