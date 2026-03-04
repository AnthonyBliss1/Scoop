<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import ServerUpload from "@lucide/svelte/icons/cloud-upload";
  import ServerDownload from "@lucide/svelte/icons/cloud-download";
  import { SyncServer } from "../../../../bindings/changeme";
  import { getAppState } from "$lib/store/AppState.svelte";

  const app = getAppState();
  let { cmd = $bindable("Run Sync") } = $props<{
    cmd: any;
  }>();

  let loading: boolean = $state(false);

  async function openSyncServer() {
    try {
      app.currentServer = await SyncServer.OpenSyncServer();
    } catch (error) {
      console.error(error);
    }
  }

  async function sendToServer() {
    loading = true;
    try {
      const ok = await SyncServer.SendToServer(app.currentServer);

      if (ok) {
        toast.success(`App data sucessfully pushed to "${app.currentServer.name}"`);
        loading = false;
        cmd = null;
      }
    } catch (error) {
      loading = false;
      console.error(error);
    }
  }

  // TODO: add the receive from server function

  onMount(() => {
    if (app.currentServer.name !== "" && app.currentServer.url !== "") {
      openSyncServer();
    } else {
      toast.warning("No server has been set");
    }
  });
</script>

<div
  class="border-border bg-background flex min-h-[20vh] flex-col items-center justify-center gap-5 rounded-sm border p-5"
>
  <div class="flex h-full flex-col items-center justify-center gap-10">
    <p class="text-lg underline decoration-1 underline-offset-3">Server Operations</p>

    <div class="flex w-full flex-row items-center justify-center gap-10">
      <div class="flex flex-col items-center justify-center gap-3">
        <button
          class="border-border bg-accent focus:ring-offset-background inline-flex h-9 w-50 items-center justify-center
        rounded-sm border text-sm
        text-green-300 hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
          onclick={sendToServer}
          disabled={loading}
        >
          <ServerUpload class="mr-2" size={20} />
          Push To Server
        </button>
        <p class="text-xs">* Sync Remote Data</p>
      </div>

      <div class="flex flex-col items-center justify-center gap-3">
        <button
          class="border-border bg-accent focus:ring-offset-background inline-flex h-9 w-50 items-center justify-center
        rounded-sm border text-sm
        text-green-300 hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
          onclick={() => {
            cmd = null;
          }}
          disabled={loading}
        >
          <ServerDownload class="mr-2" size={20} />
          Pull From Server</button
        >
        <p class="text-xs">* Sync Local Data</p>
      </div>
    </div>
  </div>
</div>
