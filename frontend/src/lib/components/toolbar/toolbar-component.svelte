<script lang="ts">
  import { getAppState } from "$lib/store/AppState.svelte";
  import { ScoopService, Scoop } from "../../../../bindings/changeme";

  import Package from "@lucide/svelte/icons/package";
  import Info from "@lucide/svelte/icons/info";
  import File from "@lucide/svelte/icons/file-braces";
  import { toast } from "svelte-sonner";
  import { onDestroy, onMount } from "svelte";

  const appState = getAppState();

  // using this for the effect to track when the collection changes and reset the windowStart index
  let lastCollection = appState.currentCollection.name;

  let { showCurl = $bindable<boolean>(false), showHelp = $bindable<boolean>(false) } = $props<{
    showCurl: boolean;
    showHelp: boolean;
  }>();

  const windowSize: number = 5;
  let windowStart: number = $state(0); // window start index of allScoops

  // reset windowStart on new collection
  $effect(() => {
    const currentCollection = appState.currentCollection.name;

    if (currentCollection !== lastCollection) {
      lastCollection = currentCollection;
      windowStart = 0;
    }
  });

  // visibleScoop serves as the sliding window
  let visibleScoops: Scoop[] = $derived(
    appState.allScoops.slice(windowStart, windowStart + windowSize),
  );

  // ranges that can appear on the right or left side of the window as the window slides
  let hiddenLeftStart: number = $derived(1);
  let hiddenLeftEnd: number = $derived(windowStart);

  let hiddenRightStart: number = $derived(windowStart + windowSize + 1);
  let hiddenRightEnd: number = $derived(appState.allScoops.length);

  let hasLeftHidden: boolean = $derived(windowStart > 0);
  let hasRightHidden: boolean = $derived(windowStart + windowSize < appState.allScoops.length);

  function persistFormToRequest(s: Scoop) {
    s.request.method = appState.method;
    s.request.url = appState.url;
    s.response = appState.response;
    s.request.headers = appState.headers;
    s.request.query_params = appState.queryParams;
    s.request.body = appState.body;
  }

  async function onSwitchScoop(): Promise<boolean> {
    try {
      persistFormToRequest(appState.currentScoop);
      return await ScoopService.SaveScoop(appState.currentScoop, appState.currentCollection);
    } catch (error) {
      console.log(error);
      return false;
    }
  }

  async function generateCurl() {
    persistFormToRequest(appState.currentScoop);

    if (appState.currentScoop.request.url === "") {
      toast.warning("Please create a valid request first");
      return;
    }

    try {
      appState.curlCommand = await ScoopService.GenerateCurlCommand(appState.currentScoop);
    } catch (error) {
      console.error(error);
    } finally {
      showCurl = true;
    }
  }

  const switchRequest = async (event: KeyboardEvent) => {
    if (event.ctrlKey && event.key >= "1" && event.key <= "9") {
      const localIndex = Number(event.key) - 1;

      if (localIndex >= visibleScoops.length) return;

      const ok = await onSwitchScoop();
      if (!ok) return;

      const nextScoop = visibleScoops[localIndex];
      appState.currentScoop = nextScoop;
      appState.response = nextScoop.response ?? "";
      appState.queryParams = nextScoop.request.query_params;
    }
  };

  const moveWindow = (event: KeyboardEvent) => {
    if (event.ctrlKey && event.key == "]") {
      windowStart = Math.min(
        windowStart + windowSize,
        Math.max(0, appState.allScoops.length - windowSize),
      );
    }

    if (event.ctrlKey && event.key == "[") {
      windowStart = Math.max(0, windowStart - windowSize);
    }
  };

  onMount(() => {
    document.addEventListener("keydown", switchRequest);
    document.addEventListener("keydown", moveWindow);
  });

  onDestroy(() => {
    document.removeEventListener("keydown", switchRequest);
    document.removeEventListener("keydown", moveWindow);
  });
</script>

<div class="-mx-10 h-8 shrink-0 overflow-hidden rounded-b-sm bg-green-950/30">
  <div class="flex h-full min-w-0 items-center gap-5 px-10 text-sm text-green-500/90">
    <div class="flex shrink-0 flex-row gap-2">
      <Package
        class={appState.currentCollection.name === "temp" ? "text-blue-500/90" : ""}
        size={20}
      />
      <p class={appState.currentCollection.name === "temp" ? "text-blue-500/90" : "text-green-400"}>
        {appState.currentCollection.name}
      </p>
    </div>

    <!-->Left Range<-->
    {#if hasLeftHidden}
      <div class="text-green-400">
        <p>[{hiddenLeftStart} - {hiddenLeftEnd}]</p>
      </div>
    {/if}

    <!-->Visible Scoops (within the sliding window)<-->
    {#if appState.currentScoop.name !== "temp"}
      {#each visibleScoops as scoop, i (scoop.id)}
        <div class="flex flex-row gap-1">
          <p class={appState.currentScoop.id === scoop.id ? "text-blue-500" : "text-green-400"}>
            ({i + 1})
          </p>
          <p class={appState.currentScoop.id === scoop.id ? "text-blue-500" : "text-green-400"}>
            {scoop.name}
          </p>
        </div>
      {/each}

      <!-->Right Range<-->
      {#if hasRightHidden}
        <div class="text-green-400">
          <p>[{hiddenRightStart} - {hiddenRightEnd}]</p>
        </div>
      {/if}
    {/if}

    <div class="ml-auto flex h-8 items-center justify-center gap-5">
      <button
        class="flex w-8 items-center justify-center rounded-sm border-blue-400/90 p-1 hover:border focus:outline-none"
        onclick={() => {
          generateCurl();
        }}
      >
        <File class="text-blue-400/90" size={20} />
      </button>
      <button
        class="flex w-8 items-center justify-center rounded-sm border-blue-500/90 p-1 hover:border focus:outline-none"
        onclick={() => {
          showHelp = true;
        }}
      >
        <Info class="text-blue-500/90" size={20} />
      </button>
    </div>
  </div>
</div>
