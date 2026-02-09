<script lang="ts">
  import PackagePlus from "@lucide/svelte/icons/package-plus";
  import PackageOpen from "@lucide/svelte/icons/package-open";
  import Send from "@lucide/svelte/icons/send";
  import Server from "@lucide/svelte/icons/server";
  import Database from "@lucide/svelte/icons/database";
  import CloudDownload from "@lucide/svelte/icons/cloud-download";

  import * as Command from "$lib/components/ui/command/index.js";
  import { onMount } from "svelte";

  let cmd: string = $state("Create New Request");
  let inputEl: HTMLInputElement | null = $state(null);

  onMount(() => {
    inputEl?.focus();
  });
</script>

<Command.Root bind:value={cmd} class="border-border rounded-sm  border shadow-md md:min-w-[450px]">
  <Command.Input bind:ref={inputEl} placeholder="Type a command or search..." />

  <Command.List>
    <Command.Empty>No results found.</Command.Empty>

    <Command.Group heading="Suggested">
      <Command.Item value="Create New Request">
        <Send class="text-green-500" />
        <span class="text-green-300">Create New Request</span>
        <Command.Shortcut class="text-green-500">⌘R</Command.Shortcut>
      </Command.Item>
    </Command.Group>

    <Command.Group heading="Collections">
      <Command.Item value="Create Collection">
        <PackagePlus class="text-green-500" />
        <span class="text-green-300">Create Collection</span>
        <Command.Shortcut class="text-green-500">⌘N</Command.Shortcut>
      </Command.Item>
      <Command.Item value="Open Collection">
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
