<script lang="ts">
  import Library from "@lucide/svelte/icons/library-big";
  import Send from "@lucide/svelte/icons/send";
  import Server from "@lucide/svelte/icons/server";
  import Database from "@lucide/svelte/icons/database";
  import CloudDownload from "@lucide/svelte/icons/cloud-download";

  import * as Command from "$lib/components/ui/command/index.js";
</script>

<Command.Root class="border-border rounded-lg  border shadow-md md:min-w-[450px]">
  <Command.Input placeholder="Type a command or search..." />
  <Command.List>
    <Command.Empty>No results found.</Command.Empty>
    <Command.Group heading="Suggestions">
      <Command.Item>
        <Library class="text-green-500" />
        <span class="text-green-300">Create Collection</span>
        <Command.Shortcut class="text-green-500">⌘N</Command.Shortcut>
      </Command.Item>
      <Command.Item>
        <Send class="text-green-500" />
        <span class="text-green-300">Create New Request</span>
        <Command.Shortcut class="text-green-500">⌘R</Command.Shortcut>
      </Command.Item>
    </Command.Group>
    <Command.Separator />
    <Command.Group heading="Cloud">
      <Command.Item>
        <Server class="text-green-500" />
        <span class="text-green-300">Create DNS Alias</span>
        <Command.Shortcut class="text-green-500">⌘D</Command.Shortcut>
      </Command.Item>
      <Command.Item>
        <Database class="text-green-500" />
        <span class="text-green-300">Configure Sync</span>
        <Command.Shortcut class="text-green-500">⌘C</Command.Shortcut>
      </Command.Item>
      <Command.Item>
        <CloudDownload class="text-green-500" />
        <span class="text-green-300">Run Sync</span>
        <Command.Shortcut class="text-green-500">⌘S</Command.Shortcut>
      </Command.Item>
    </Command.Group>
  </Command.List>
</Command.Root>
