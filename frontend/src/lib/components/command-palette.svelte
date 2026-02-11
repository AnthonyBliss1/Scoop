<script lang="ts">
  import PackagePlus from "@lucide/svelte/icons/package-plus";
  import PackageOpen from "@lucide/svelte/icons/package-open";
  import Send from "@lucide/svelte/icons/send";
  import Server from "@lucide/svelte/icons/server";
  import Database from "@lucide/svelte/icons/database";
  import CloudDownload from "@lucide/svelte/icons/cloud-download";

  import * as Command from "$lib/components/ui/command/index.js";
  import { onDestroy, onMount } from "svelte";
  import CreateCollection from "./create-collection.svelte";
  import CreateRequest from "./create-request.svelte";
  import type { Collection, Request } from "bindings/changeme";
  import OpenCollection from "./open-collection.svelte";

  type Command =
    | "Create New Request"
    | "Create Collection"
    | "Open Collection"
    | "Create DNS Alias"
    | "Configure Sync"
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

  let { collection = $bindable<Collection>(), request = $bindable<Request>() } = $props<{
    collection: Collection;
    request: Request;
  }>();

  const onEnter = (event: KeyboardEvent) => {
    if (!executeCmd) return;

    if (event.key === "Enter") {
      switch (cmd) {
        case "Create New Request":
          executeCmd = "Create New Request";
          break;
        case "Create Collection":
          executeCmd = "Create Collection";
          break;
        default:
          break;
      }
    }
  };

  $effect(() => {
    if (!executeCmd) {
      inputEl?.focus();
    }
  });

  onMount(() => {
    inputEl?.focus();
    document.addEventListener("keydown", onEnter);
  });

  onDestroy(() => {
    executeCmd = null;
    document.removeEventListener("keydown", onEnter);
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
          class="hover:cursor-pointer"
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
          class="hover:cursor-pointer"
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
          class="hover:cursor-pointer"
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
        <Command.Item value="Create DNS Alias">
          <Server class="text-green-500" />
          <span class="text-green-300">Create DNS Alias</span>
          <Command.Shortcut class="text-green-500">⌘D</Command.Shortcut>
        </Command.Item>
        <Command.Item value="Configure Sync">
          <Database class="text-green-500" />
          <span class="text-green-300">Configure Sync</span>
          <Command.Shortcut class="text-green-500">⌘C</Command.Shortcut>
        </Command.Item>
        <Command.Item value="Run Sync">
          <CloudDownload class="text-green-500" />
          <span class="text-green-300">Run Sync</span>
          <Command.Shortcut class="text-green-500">⌘S</Command.Shortcut>
        </Command.Item>
      </Command.Group>
    </Command.List>
  </Command.Root>
{:else if executeCmd === "Create New Request"}
  <CreateRequest bind:cmd={executeCmd} bind:request bind:collection />
{:else if executeCmd === "Create Collection"}
  <CreateCollection bind:cmd={executeCmd} bind:collection />
{:else if executeCmd === "Open Collection"}
  <OpenCollection bind:cmd={executeCmd} bind:request bind:collection />
{/if}
