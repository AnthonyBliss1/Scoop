<script lang="ts">
  import { onMount } from "svelte";
  import Send from "@lucide/svelte/icons/send";
  import * as Command from "$lib/components/ui/command/index.js";
  import { ScoopService, Collection, Scoop } from "../../../../bindings/changeme";
  import { getAppState } from "$lib/store/AppState.svelte";

  const app = getAppState();

  let { cmd = $bindable("Open Collection") } = $props<{
    cmd: any;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);
  let availCollections: Collection[] | null = $state<Collection[] | null>(null);

  async function OpenCollections() {
    try {
      availCollections = await ScoopService.OpenCollections();

      if (availCollections.length > 0) {
        console.log(`Loaded ${availCollections.length} Collection(s)`);
      }
    } catch (error) {
      console.error(error);
    }
  }

  onMount(() => {
    inputEl?.focus();

    OpenCollections();
  });
</script>

<Command.Root class="border-border rounded-sm  border shadow-md md:min-w-[450px]">
  <Command.Input bind:ref={inputEl} placeholder="Search collections..." />

  <Command.List>
    <Command.Empty>No collections found</Command.Empty>

    {#if availCollections && availCollections.length > 0}
      <Command.Group heading="Collections">
        {#each availCollections as c (c.id)}
          <Command.Item
            value={c.id}
            onclick={() => {
              app.currentCollection = c;
              app.allScoops = c.scoops;

              if (app.allScoops.length > 0) {
                app.currentScoop = app.allScoops[0];
              } else {
                app.currentScoop = new Scoop({ name: "temp" });
              }
              console.log(`allScoops length : ${app.allScoops.length}`);
              cmd = null;
            }}
          >
            <p class="text-green-500">{c.scoops.length}</p>
            <Send class="text-green-500" />
            <span class="text-green-300">{c.name}</span>
          </Command.Item>
        {/each}
      </Command.Group>
    {/if}
  </Command.List>
</Command.Root>
