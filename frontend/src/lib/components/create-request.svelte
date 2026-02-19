<script lang="ts">
  import { toast } from "svelte-sonner";
  import { Backend, Collection, Request, Scoop } from "../../../bindings/changeme";

  // TODO: Think about renaming this to create-scoop
  // originally was request but recently created the scoop structure (consisting of Request + Response)
  // not sure how committed i am to this

  let {
    cmd = $bindable("Create New Request"),
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

  let tempScoop: Scoop = $state(new Scoop({ request: new Request() }));
  let newScoop: string = $state("");

  async function createScoop() {
    if (newScoop === "" || newScoop === "temp") {
      toast.error("Please enter a valid name");
      return;
    }

    tempScoop.name = newScoop;

    try {
      const ok = await Backend.CreateScoop(collection, tempScoop);

      if (ok) {
        allScoops.push(tempScoop);

        // only override currentScoop if created scoop is the first in the collection
        if (currentScoop.name === "temp") currentScoop = allScoops[0];

        // this seems to work here but not sure why, need to investigate
        collection.scoops.push(tempScoop);
        console.log(`Created Request: ${tempScoop.name}`);
      }
    } catch (error) {
      console.error(error);
    } finally {
      cmd = null;
    }
  }

  $effect(() => {
    if (cmd === "Create New Request") {
      inputEl?.focus();
    }
  });
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex items-center justify-center">Create New Request</p>
  <input
    class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none md:min-w-[450px]"
    bind:value={newScoop}
    bind:this={inputEl}
    placeholder="Enter Request Name..."
  />

  <div class="flex w-full flex-row items-center justify-center gap-10">
    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:cursor-pointer hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={createScoop}>Create</button
    >

    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:cursor-pointer hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={() => {
        cmd = null;
      }}>Cancel</button
    >
  </div>
</div>
