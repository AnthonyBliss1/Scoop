<script lang="ts">
  import { getAppState } from "$lib/store/AppState.svelte";
  import { ScoopService } from "../../../../bindings/changeme";
  import { toast } from "svelte-sonner";

  const app = getAppState();

  let { showDeleteCollection = $bindable<boolean>() } = $props<{
    showDeleteCollection: boolean;
  }>();

  let deleting: boolean = $state(false);

  let btnEl: HTMLButtonElement | null = $state(null);

  $effect(() => {
    if (showDeleteCollection) {
      btnEl?.focus();
    }
  });

  async function deleteCollection() {
    const target = app.currentCollection;

    // this shouldnt fire, just in case
    if (target.name === "temp") {
      toast.warning("Cannot delete 'temp' collection");
      return;
    }

    deleting = true;

    try {
      const ok = await ScoopService.DeleteCollection(target);

      if (ok) {
        deleting = false;
        showDeleteCollection = false;
        app.reset();

        toast.success(`Deleted '${target.name}' Collection`);
      }
    } catch (error) {
      console.error(error);
      deleting = false;
    }
  }
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex items-center justify-center">
    Are you sure you want to delete '{app.currentCollection.name}'?
  </p>

  <div class="flex w-full flex-row items-center justify-center gap-10">
    <button
      class="bg-accent focus:ring-offset-background inline-flex h-9 items-center justify-center rounded-sm
      border border-red-500 px-3 text-sm
      text-red-500 hover:bg-red-400 hover:text-black focus:ring-2
      focus:ring-red-400/20 focus:ring-offset-2 focus:outline-none
      disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-red-500"
      disabled={deleting}
      onclick={deleteCollection}>Delete</button
    >

    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
      justify-center rounded-sm border px-3
      text-sm hover:bg-green-400 hover:text-black focus:ring-2
      focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
      disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      bind:this={btnEl}
      disabled={deleting}
      onclick={() => {
        showDeleteCollection = false;
      }}>Cancel</button
    >
  </div>
</div>
