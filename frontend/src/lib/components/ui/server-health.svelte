<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import { getAppState } from "$lib/store/AppState.svelte";
  import { SyncServer } from "../../../../bindings/changeme";
  import { onDestroy, onMount } from "svelte";

  const app = getAppState();

  type Connection = "Offline" | "Online";
  let status: Connection = $state("Offline");

  let onHealthCheck: undefined | (() => void);

  async function openSyncServer() {
    try {
      app.currentServer = await SyncServer.OpenSyncServer();
      console.info(`Server set! ${app.currentServer.name}`);
    } catch (error) {
      console.error(error);
    }
  }

  onMount(async () => {
    await openSyncServer();

    onHealthCheck = Events.On("serverHealth", async (event: any) => {
      const ok = event.data as Connection;
      status = ok;
    });

    // make sure server is actually set before initiating the health check poll on backend
    if (app.currentServer.name !== "" && app.currentServer.url !== "") {
      Events.Emit("initiateHealthCheck", app.currentServer);
    }
  });

  onDestroy(() => {
    onHealthCheck?.();
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
