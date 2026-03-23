<script lang="ts">
  import * as Command from "$lib/components/ui/command/index.js";
</script>

<Command.Root class="border-border rounded-sm  border shadow-md md:min-w-[450px]">
  <Command.List>
    <Command.Group heading="Keybindings">
      <Command.Item value="Open Command Palette" class="aria-selected:bg-background">
        <span class="text-green-400">Open Command Palette</span>
        <Command.Shortcut class="text-green-500">Shift + Ctrl + P</Command.Shortcut>
      </Command.Item>

      <Command.Item value="Go to Scoop" class="aria-selected:bg-background">
        <span class="text-green-400">Go to Scoop</span>
        <Command.Shortcut class="text-green-500">Ctrl + [Num]</Command.Shortcut>
      </Command.Item>

      <Command.Item value="Rename Scoop" class="aria-selected:bg-background">
        <span class="text-green-400">Rename Scoop</span>
        <Command.Shortcut class="text-green-500">Ctrl + R</Command.Shortcut>
      </Command.Item>

      <Command.Item value="Delete Current Collection" class="aria-selected:bg-background">
        <span class="text-green-400">Delete Current Collection</span>
        <Command.Shortcut class="text-green-500">Ctrl + D</Command.Shortcut>
      </Command.Item>

      <Command.Item value="Expand Response Body" class="aria-selected:bg-background">
        <span class="text-green-400">Expand Response Body</span>
        <Command.Shortcut class="text-green-500">Ctrl + E</Command.Shortcut>
      </Command.Item>
    </Command.Group>
  </Command.List>
</Command.Root>
