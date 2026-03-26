<script lang="ts">
  import { toast } from "svelte-sonner";
  import { ScoopService, Scoop } from "../../../../bindings/changeme";
  import { getAppState } from "$lib/store/AppState.svelte";

  const app = getAppState();

  type Item = "Collection" | "Scoop";

  let { showRename = $bindable<boolean>(), mode = $bindable<Item>() } = $props<{
    showRename: boolean;
    mode: Item;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);

  let newName: string = $derived.by(() => {
    if (mode === "Scoop") {
      return app.currentScoop.name;
    } else {
      return app.currentCollection.name;
    }
  });

  const renameFunc: () => void = $derived.by(() => {
    return mode === "Scoop" ? renameScoop : renameCollection;
  });

  async function renameScoop() {
    if (newName === "" || newName === "temp") {
      toast.error("Please enter a valid name");
      return;
    }

    try {
      replaceScoop();
      const ok = await ScoopService.SaveScoop(app.currentScoop, app.currentCollection);

      if (ok) {
        toast.success(`Renamed Request: ${app.currentScoop.name}`);
        console.log(`Renamed Request: ${app.currentScoop.name}`);
      }
    } catch (error) {
      console.error(error);
    } finally {
      showRename = false;
    }
  }

  async function renameCollection() {
    if (newName === "" || newName === "temp") {
      toast.error("Please enter a valid name");
      return;
    }

    const updatedCollection = {
      ...app.currentCollection,
      name: newName,
    };

    try {
      const ok = await ScoopService.SaveCollection(updatedCollection);

      if (ok) {
        app.currentCollection = updatedCollection;
        toast.success(`Renamed Collection: ${app.currentCollection.name}`);
        console.log(`Renamed Collection: ${app.currentCollection.name}`);
      }
    } catch (error) {
      console.error(error);
    } finally {
      showRename = false;
    }
  }

  function replaceScoop() {
    // find the index of the scoop to modify (currentScoop)
    const idx = app.allScoops.findIndex((s: Scoop) => s.name === app.currentScoop.name);
    if (idx === -1) return;

    // need to create temp array
    // (crete new references for reactivity)
    let tempAllScoop = [...app.allScoops];
    tempAllScoop[idx] = { ...tempAllScoop[idx], name: newName };

    // overwrite the main arrays
    // (making sure to create new references)
    app.allScoops = tempAllScoop;
    app.currentCollection = { ...app.currentCollection, scoops: tempAllScoop };

    // overwrite name of the currentScoop
    app.currentScoop.name = newName;
  }

  $effect(() => {
    if (showRename === true) {
      inputEl?.focus();
    }
  });
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex items-center justify-center">
    {mode === "Scoop" ? "Rename Scoop" : "Rename Collection"}
  </p>
  <input
    class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none md:min-w-[450px]"
    bind:value={newName}
    bind:this={inputEl}
    placeholder={mode === "Scoop" ? "Enter New Scoop Name..." : "Enter New Collection Name..."}
  />

  <div class="flex w-full flex-row items-center justify-center gap-10">
    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={renameFunc}>Rename</button
    >

    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={() => {
        showRename = false;
      }}>Cancel</button
    >
  </div>
</div>
