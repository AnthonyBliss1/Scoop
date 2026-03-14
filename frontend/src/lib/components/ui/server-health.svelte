<script lang="ts">
  import { getAppState } from "$lib/store/AppState.svelte";
  import { SyncServer } from "../../../../bindings/changeme";
  import { onMount } from "svelte";

  const app = getAppState();

  type Connection = "Offline" | "Online";

  let status: Connection = $state("Offline");

  async function openSyncServer() {
    try {
      app.currentServer = await SyncServer.OpenSyncServer();
      console.info(`Server set! ${app.currentServer.name}`);
    } catch (error) {
      console.error(error);
    }
  }

  async function checkServerHealth() {
    try {
      const ok = await SyncServer.CheckServerHealth(app.currentServer);

      // set the server status accordingly
      if (ok) {
        status = "Online";
      } else {
        status = "Offline";
      }
    } catch (error) {
      console.error(error);
    }
  }

  onMount(async () => {
    await openSyncServer();

    // make sure server is actually set before calling /health
    if (app.currentServer.name !== "" && app.currentServer.url !== "") {
      checkServerHealth();
    }
  });
</script>

<span class="relative flex h-3 w-3">
  <span
    class={status === "Online"
      ? `absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400/75`
      : `absolute inline-flex h-full w-full rounded-full bg-red-400/75`}
  ></span>
  <span
    class={status === "Online"
      ? `relative inline-flex h-3 w-3 rounded-full bg-green-600`
      : `relative inline-flex h-3 w-3 rounded-full bg-red-900`}
  ></span>
</span>
<span class={status === "Online" ? `text-green-500/90` : `text-red-500/75`}>{status}</span>
