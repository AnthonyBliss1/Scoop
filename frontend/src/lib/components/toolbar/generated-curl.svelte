<script lang="ts">
  import { getAppState } from "$lib/store/AppState.svelte";
  import Clipboard from "@lucide/svelte/icons/clipboard";
  import { toast } from "svelte-sonner";

  const app = getAppState();

  let { showCurl = $bindable<boolean>(false) } = $props<{ showCurl: boolean }>();

  let curlCommand = $state("");

  let btnEl: HTMLButtonElement | null = $state(null);

  $effect(() => {
    if (showCurl) {
      btnEl?.focus();
    }
  });

  $effect(() => {
    if (app.curlCommand !== "") {
      curlCommand = app.curlCommand;
    }
  });

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(curlCommand);
      toast.success("Copied cURL to clipboard");
    } catch (err) {
      console.error("Failed to copy: ", err);
      toast.success("Failed to copy");
    }
  }
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex h-full items-center justify-center text-lg">Generated cURL Command</p>

  <div class="flex flex-row items-center gap-5">
    <button
      class="flex h-full w-10 items-center justify-center rounded-sm border-green-300 p-1 hover:border hover:text-green-300 focus:outline-none"
      bind:this={btnEl}
      onclick={copyToClipboard}
    >
      <Clipboard class="h-full" size={28} />
    </button>

    <textarea
      class="focus:ring-offset-background bg-accent border-border h-13 min-w-0
    resize-none rounded-sm border px-2 leading-12 text-green-300
    shadow-md focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none md:min-w-[450px]"
      placeholder="Generating cURL command ..."
      readonly
      wrap="off"
      bind:value={curlCommand}
    >
    </textarea>
  </div>
</div>
