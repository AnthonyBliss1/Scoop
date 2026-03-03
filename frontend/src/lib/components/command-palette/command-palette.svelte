<script lang="ts">
  import PackagePlus from "@lucide/svelte/icons/package-plus";
  import PackageOpen from "@lucide/svelte/icons/package-open";
  import Send from "@lucide/svelte/icons/send";
  import ServerIcon from "@lucide/svelte/icons/server";
  import Database from "@lucide/svelte/icons/database";
  import Cloud from "@lucide/svelte/icons/cloud-backup";

  import * as Command from "$lib/components/ui/command/index.js";
  import { onDestroy, onMount } from "svelte";
  import CreateCollection from "./create-collection.svelte";
  import CreateRequest from "./create-request.svelte";
  import type { Collection, Scoop, Server } from "../../../../bindings/changeme";
  import OpenCollection from "./open-collection.svelte";
  import DnsOverride from "./dns-override.svelte";
  import SetSyncServer from "./set-sync-server.svelte";
  import RunSync from "./run-sync.svelte";

  type Command =
    | "Create New Request"
    | "Create Collection"
    | "Open Collection"
    | "Configure DNS Override"
    | "Set Sync Server"
    | "Run Sync";

  let cmd: Command = $derived.by((): Command => {
    if (collection.name === "temp") {
      return "Create Collection";
    } else {
      return "Create New Request";
    }
  });
  let executeCmd = $state<Command | null>(null);

  let inputEl: HTMLInputElement | null = $state(null);

  let {
    collection = $bindable<Collection>(),
    allScoops = $bindable<Scoop[]>(),
    currentScoop = $bindable<Scoop>(),
    currentServer = $bindable<Server>(),
  } = $props<{
    collection: Collection;
    allScoops: Scoop[];
    currentScoop: Scoop;
    currentServer: Server;
  }>();

  $effect(() => {
    if (!executeCmd) {
      inputEl?.focus();
    }
  });

  onMount(() => {
    inputEl?.focus();
  });

  onDestroy(() => {
    executeCmd = null;
  });
</script>

{#if !executeCmd}
  <Command.Root
    class="border-border rounded-sm  border shadow-md md:min-w-[450px]"
    bind:value={cmd}
  >
    <Command.Input bind:ref={inputEl} placeholder="Type a command or search..." />

    <Command.List>
      <Command.Empty>No results found</Command.Empty>

      <Command.Group heading="Suggested">
        <Command.Item
          value="Create New Request"
          disabled={collection.name === "temp" ? true : false}
          onclick={() => {
            executeCmd = "Create New Request";
          }}
        >
          <Send class="text-green-500" />
          <span class="text-green-300">Create New Request</span>
          <Command.Shortcut class="text-green-500">⌘R</Command.Shortcut>
        </Command.Item>
      </Command.Group>

      <Command.Group heading="Collections">
        <Command.Item
          value="Create Collection"
          onclick={() => {
            executeCmd = "Create Collection";
          }}
        >
          <PackagePlus class="text-green-500" />
          <span class="text-green-300">Create Collection</span>
          <Command.Shortcut class="text-green-500">⌘N</Command.Shortcut>
        </Command.Item>
        <Command.Item
          value="Open Collection"
          onclick={() => {
            executeCmd = "Open Collection";
          }}
        >
          <PackageOpen class="text-green-500" />
          <span class="text-green-300">Open Collection</span>
          <Command.Shortcut class="text-green-500">⌘O</Command.Shortcut>
        </Command.Item>
      </Command.Group>

      <Command.Group heading="Cloud">
        <Command.Item
          value="Create DNS Alias"
          onclick={() => {
            executeCmd = "Configure DNS Override";
          }}
        >
          <ServerIcon class="text-green-500" />
          <span class="text-green-300">Configure DNS Override</span>
          <Command.Shortcut class="text-green-500">⌘D</Command.Shortcut>
        </Command.Item>
        <Command.Item
          value="Set Sync Server"
          onclick={() => {
            executeCmd = "Set Sync Server";
          }}
        >
          <Database class="text-green-500" />
          <span class="text-green-300">Set Sync Server</span>
          <Command.Shortcut class="text-green-500">⌘C</Command.Shortcut>
        </Command.Item>
        <Command.Item
          value="Run Sync"
          onclick={() => {
            executeCmd = "Run Sync";
          }}
        >
          <Cloud class="text-green-500" />
          <span class="text-green-300">Run Sync</span>
          <Command.Shortcut class="text-green-500">⌘S</Command.Shortcut>
        </Command.Item>
      </Command.Group>
    </Command.List>
  </Command.Root>
{:else if executeCmd === "Create New Request"}
  <CreateRequest bind:cmd={executeCmd} bind:allScoops bind:collection bind:currentScoop />
{:else if executeCmd === "Create Collection"}
  <CreateCollection bind:cmd={executeCmd} bind:allScoops bind:collection bind:currentScoop />
{:else if executeCmd === "Open Collection"}
  <OpenCollection bind:cmd={executeCmd} bind:allScoops bind:collection bind:currentScoop />
{:else if executeCmd === "Configure DNS Override"}
  <DnsOverride bind:cmd={executeCmd} />
{:else if executeCmd === "Set Sync Server"}
  <SetSyncServer bind:cmd={executeCmd} bind:currentServer />
{:else if executeCmd === "Run Sync"}
  <RunSync bind:cmd={executeCmd} bind:currentServer />
{/if}
