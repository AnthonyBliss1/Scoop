<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import { Server, SyncServer } from "../../../../bindings/changeme";
  import { getAppState } from "$lib/store/AppState.svelte";

  const app = getAppState();

  let { cmd = $bindable("Set Sync Server") } = $props<{
    cmd: any;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);

  let serverKey: string = $derived(app.currentServer.key);
  let serverURL: string = $derived(app.currentServer.url);

  async function setSyncServer() {
    if (serverKey === "" || serverURL === "") {
      toast.error("Please enter a valid key and URL");
      return;
    }

    try {
      app.currentServer = new Server({ key: serverKey, url: serverURL });

      const ok = await SyncServer.SetSyncServer(app.currentServer);

      if (ok) {
        toast.success("Scoop Server Successfully Set");
        console.log("Scoop Server Successfully Set");
      }
    } catch (error) {
      console.error(error);
    } finally {
      cmd = null;
    }
  }

  async function openSyncServer() {
    try {
      app.currentServer = await SyncServer.OpenSyncServer();
    } catch (error) {
      console.error(error);
    }
  }

  $effect(() => {
    if (cmd === "Set Sync Server") {
      inputEl?.focus();
    }
  });

  onMount(() => {
    if (app.currentServer.key === "" && app.currentServer.url === "") {
      openSyncServer();
    }
  });
</script>

<div class="border-border bg-background flex min-h-[20vh] flex-col gap-5 rounded-sm border p-5">
  <div class="flex h-full flex-col items-center justify-center gap-10">
    <p class="text-lg underline decoration-1 underline-offset-3">Server Credentials</p>
    <div class="flex h-full flex-row items-center justify-center gap-5">
      <input
        class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
        bind:value={serverKey}
        bind:this={inputEl}
        placeholder="Enter server API key ..."
      />
      <input
        class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
        bind:value={serverURL}
        placeholder="Enter URL ..."
      />
    </div>

    <div class="flex w-full flex-row items-center justify-center gap-5">
      <button
        class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
        onclick={setSyncServer}>Create</button
      >

      <button
        class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
        onclick={() => {
          cmd = null;
        }}>Cancel</button
      >
    </div>
  </div>
</div>
