<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import { getAppState } from "$lib/store/AppState.svelte";
  import { SyncServer } from "../../../../bindings/changeme";
  import Alert from "@lucide/svelte/icons/triangle-alert";
  import { onDestroy, onMount } from "svelte";

  const app = getAppState();

  type Connection = "Offline" | "Online";
  let status: Connection = $state("Offline");

  let verWarning: boolean = $state(false);

  let onHealthCheck: undefined | (() => void);

  async function openSyncServer() {
    try {
      app.currentServer = await SyncServer.OpenSyncServer();
      await openSyncData();
    } catch (error) {
      console.error(error);
    }
  }

  async function openSyncData() {
    try {
      app.syncData = await SyncServer.OpenSyncData();
    } catch (error) {
      console.error(error);
    }
  }

  async function checkServerSyncData() {
    if (status === "Offline") return;

    try {
      const sData = await SyncServer.CheckServerSyncData(app.currentServer);

      if (sData.versionNum != 0) {
        verWarning = sData.versionNum !== app.syncData.versionNum;
      }
    } catch (error) {
      console.error(error);
    }
  }

  $effect(() => {
    if (app.currentServer.key !== "" && app.currentServer.url !== "") {
      Events.Emit("initiateHealthCheck", app.currentServer);

      checkServerSyncData();
    }
  });

  // this effect is in place to run if app state is reset after a server pull,
  // so if syncData is reset then open it and check it against the server
  $effect(() => {
    if (
      app.syncData.versionNum === 0 &&
      app.currentServer.key !== "" &&
      app.currentServer.url !== ""
    ) {
      openSyncData();
      checkServerSyncData();
    }
  });

  onMount(async () => {
    await openSyncServer();

    onHealthCheck = Events.On("serverHealth", async (event: any) => {
      const ok = event.data as Connection;
      status = ok;
    });

    // make sure server is actually set before initiating the health check poll on backend
    if (app.currentServer.key !== "" && app.currentServer.url !== "") {
      Events.Emit("initiateHealthCheck", app.currentServer);
    }
  });

  onDestroy(() => {
    onHealthCheck?.();
  });
</script>

{#if app.currentServer.key !== "" && app.currentServer.url != ""}
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
  <div class="flex flex-row gap-3.5">
    <span class={status === "Online" ? `text-green-500/90` : `text-red-500/75`}>{status}</span>

    <!-->Show warning to indicat to user data is stale<-->
    {#if verWarning}
      <span
        title="Local data is stale, please pull current data from server"
        class="bg-background inline-flex text-green-500"
      >
        <Alert class="text-yellow-300" size={22} />
      </span>

      <!-->Should include the lastUpdated Date here<-->
    {/if}
  </div>
{/if}
