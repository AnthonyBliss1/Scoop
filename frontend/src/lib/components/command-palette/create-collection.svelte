<script lang="ts">
  import { toast } from "svelte-sonner";
  import { ScoopService, Collection, Scoop } from "../../../../bindings/changeme";

  let {
    cmd = $bindable("Create Collection"),
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

  let tempCollection: Collection = $state(new Collection());
  let newCollection: string = $state(""); // using string since user is just prompted for a name

  async function createCollection() {
    if (newCollection === "" || newCollection === "temp") {
      toast.error("Please enter a valid name");
      return;
    }

    tempCollection.name = newCollection;

    try {
      const ok = await ScoopService.CreateCollection(tempCollection);

      if (ok) {
        collection = tempCollection;
        currentScoop = new Scoop({ name: "temp" });
        allScoops = [];
        console.log(`Created Collection: ${collection.name}`);
      }
    } catch (error) {
      console.error(error);
    } finally {
      cmd = "Create New Request";
    }
  }

  $effect(() => {
    if (cmd === "Create Collection") {
      inputEl?.focus();
    }
  });
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex items-center justify-center">Create Collection</p>
  <input
    class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none md:min-w-[450px]"
    bind:value={newCollection}
    bind:this={inputEl}
    placeholder="Enter Collection Name..."
  />

  <div class="flex w-full flex-row items-center justify-center gap-10">
    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={createCollection}>Create</button
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
