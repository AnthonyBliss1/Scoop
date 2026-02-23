<script lang="ts">
  import Send from "@lucide/svelte/icons/send";
  import * as Command from "$lib/components/ui/command/index.js";
  import { Backend, Collection, Scoop } from "../../../bindings/changeme";
  import { onMount } from "svelte";

  // binding current collection and request to component
  // will change this to bind a list of requests instead of single request
  let {
    cmd = $bindable("Open Collection"),
    allScoops = $bindable<Scoop[]>(),
    collection = $bindable<Collection>(),
    currentScoop = $bindable<Scoop>(),
  } = $props<{
    cmd: any;
    allScoops: Scoop[];
    collection: Collection;
    currentScoop: Scoop;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);
  let availCollections: Collection[] | null = $state<Collection[] | null>(null);

  async function OpenCollections() {
    try {
      availCollections = await Backend.OpenCollections();

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
        {#each availCollections as c}
          <Command.Item
            value={c.name}
            onclick={() => {
              collection = c;
              allScoops = c.scoops;
              currentScoop = allScoops[0];
              console.log(`allScoops length : ${allScoops.length}`);
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
