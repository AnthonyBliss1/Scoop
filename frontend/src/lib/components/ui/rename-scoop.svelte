<script lang="ts">
  import { toast } from "svelte-sonner";
  import { ScoopService, Scoop } from "../../../../bindings/changeme";
  import { getAppState } from "$lib/store/AppState.svelte";

  const app = getAppState();

  let { showRenameScoop = $bindable<boolean>() } = $props<{
    showRenameScoop: boolean;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);

  let newScoopName: string = $state(app.currentScoop.name);

  // TODO: add check to avoid duplicate scoop names

  async function renameScoop() {
    if (newScoopName === "" || newScoopName === "temp") {
      toast.error("Please enter a valid name");
      return;
    }

    try {
      replaceScoop();
      const ok = await ScoopService.SaveScoop(app.currentScoop, app.currentCollection);

      if (ok) {
        console.log(`Renamed Request: ${app.currentScoop.name}`);
      }
    } catch (error) {
      console.error(error);
    } finally {
      showRenameScoop = false;
    }
  }

  function replaceScoop() {
    // find the index of the scoop to modify (currentScoop)
    const idx = app.allScoops.findIndex((s: Scoop) => s.name === app.currentScoop.name);
    if (idx === -1) return;

    // need to create temp array
    // (crete new references for reactivity)
    let tempAllScoop = [...app.allScoops];
    tempAllScoop[idx] = { ...tempAllScoop[idx], name: newScoopName };

    // overwrite the main arrays
    // (making sure to create new references)
    app.allScoops = tempAllScoop;
    app.currentCollection = { ...app.currentCollection, scoops: tempAllScoop };

    // overwrite name of the currentScoop
    app.currentScoop.name = newScoopName;
  }

  $effect(() => {
    if (showRenameScoop === true) {
      inputEl?.focus();
    }
  });
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex items-center justify-center">Rename Scoop</p>
  <input
    class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none md:min-w-[450px]"
    bind:value={newScoopName}
    bind:this={inputEl}
    placeholder="Enter New Scoop Name..."
  />

  <div class="flex w-full flex-row items-center justify-center gap-10">
    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={renameScoop}>Rename</button
    >

    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={() => {
        showRenameScoop = false;
      }}>Cancel</button
    >
  </div>
</div>
