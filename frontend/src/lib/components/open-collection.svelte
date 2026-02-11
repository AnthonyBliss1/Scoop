<script lang="ts">
  import Send from "@lucide/svelte/icons/send";
  import * as Command from "$lib/components/ui/command/index.js";
  import { Backend, Collection, Request } from "../../../bindings/changeme";
  import { onMount } from "svelte";

  // binding current collection and request to component
  // will change this to bind a list of requests instead of single request
  let {
    cmd = $bindable("Open Collection"),
    request = $bindable(),
    collection = $bindable<Collection>(),
  } = $props<{
    cmd: any;
    request: Request;
    collection: Collection;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);
  let availCollections: Collection[] = $state<Collection[]>([]);

  async function OpenCollections() {
    try {
      availCollections = await Backend.OpenCollections();

      if (availCollections.length > 0) {
        console.log(`Loaded ${availCollections.length} Collections`);
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

    <Command.Group heading="Collections">
      {#each availCollections as c}
        <Command.Item
          class="hover:cursor-pointer"
          value={c.name}
          onclick={() => {
            collection = c;
            request = c.requests[0];
            cmd = null;
          }}
        >
          <p class="text-green-500">{c.requests.length}</p>
          <Send class="text-green-500" />
          <span class="text-green-300">{c.name}</span>
        </Command.Item>
      {/each}
    </Command.Group>
  </Command.List>
</Command.Root>
