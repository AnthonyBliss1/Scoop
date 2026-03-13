<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import { onDestroy, onMount } from "svelte";
  import { Toaster } from "$lib/components/ui/sonner/index.js";
  import { toast } from "svelte-sonner";

  import { Method, ScoopService, Scoop, Response } from "../bindings/changeme";
  import { setAppState } from "$lib/store/AppState.svelte";

  import KvInput from "$lib/components/ui/kv-input.svelte";
  import RawInput from "$lib/components/ui/raw-input.svelte";
  import DotSpinner from "$lib/components/ui/dot-spinner.svelte";
  import ResponseViewer from "$lib/components/ui/response-viewer.svelte";
  import CmdPalette from "$lib/components/command-palette/command-palette.svelte";
  import RenameScoop from "$lib/components/ui/rename-scoop.svelte";
  import GeneratedCurl from "$lib/components/toolbar/generated-curl.svelte";
  import HelpKeybindings from "$lib/components/toolbar/help-keybindings.svelte";

  import Package from "@lucide/svelte/icons/package";
  import Info from "@lucide/svelte/icons/info";
  import File from "@lucide/svelte/icons/file-braces";
  import BodyInput from "$lib/components/ui/body-input.svelte";

  // encapsulates all the important reactive vars
  // easier to share between components this way
  const appState = setAppState();

  // events emitted from backend
  let onRspMsg: undefined | (() => void);
  let onErrMsg: undefined | (() => void);

  let showCmdPalette: boolean = $state(false);
  let showRenameScoop: boolean = $state(false);
  let showHelp: boolean = $state(false);
  let showCurl: boolean = $state(false);
  let reqParamsHidden: boolean = $state(false);

  let loading: boolean = $state(false);

  type TType = "raw" | "key-value" | "json";

  let headerTType: TType = $state("raw");
  let qParamTType: TType = $state("raw");

  // TODO: support body payloads in request on backend
  let bodyTType: TType = $state("json");

  function methodColor(method: string): string {
    switch (method) {
      case "GET":
        return `text-blue-500`;

      case "POST":
        return `text-orange-500`;

      case "PUT":
        return `text-purple-500`;

      case "PATCH":
        return `text-yellow-500`;

      case "DELETE":
        return `text-red-500`;

      default:
        return "text-green-300";
    }
  }

  function hydrateFormFromRequest(s: Scoop) {
    appState.method = (s.request.method ?? Method.Empty) as Method;
    appState.url = s.request.url ?? "";
    appState.response = s.response;
    appState.headers = s.request.headers;
    appState.queryParams = s.request.query_params;
    appState.body = s.request.body;
  }

  function persistFormToRequest(s: Scoop) {
    s.request.method = appState.method;
    s.request.url = appState.url;
    s.response = appState.response;
    s.request.headers = appState.headers;
    s.request.query_params = appState.queryParams;
    s.request.body = appState.body;
  }

  $effect(() => {
    hydrateFormFromRequest(appState.currentScoop);
  });

  async function onSend(method: Method, url: string) {
    // need to reset the response on send
    appState.resetResponse();

    let inputErr: string = "";

    if (method === "") {
      inputErr = "Please select a valid Request Method";
    } else if (url === "") {
      inputErr = "Please enter a URL";
    }

    if (inputErr !== "") {
      console.error(inputErr);
      toast.error(inputErr);
      return;
    }

    appState.currentScoop.response = new Response();
    loading = true;

    persistFormToRequest(appState.currentScoop);

    try {
      appState.currentScoop.request = await ScoopService.ModelIntializer(
        method,
        url,
        appState.headers,
        appState.queryParams,
        appState.body,
      );
      await ScoopService.SubmitRequest(appState.currentScoop);
    } catch (error) {
      console.error(error);
      loading = false;
    }
  }

  async function onSwitchRequest(): Promise<boolean> {
    try {
      persistFormToRequest(appState.currentScoop);
      const ok = await ScoopService.SaveScoop(appState.currentScoop, appState.currentCollection);
      return ok;
    } catch (error) {
      console.log(error);
      return false;
    }
  }

  async function generateCurl() {
    if (appState.currentScoop.name === "temp") {
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

  const openCmdPalette = (event: KeyboardEvent) => {
    if (event.ctrlKey && event.shiftKey && (event.key === "P" || event.code === "KeyP")) {
      showCmdPalette = true;
    }
  };

  const hideReqParams = (event: KeyboardEvent) => {
    if (event.ctrlKey && (event.key === "E" || event.code === "KeyE")) {
      reqParamsHidden = !reqParamsHidden;
    }
  };

  const renameScoop = (event: KeyboardEvent) => {
    if (event.ctrlKey && (event.key === "R" || event.code === "KeyR")) {
      if (appState.currentScoop.name === "temp") {
        toast.warning("Cannot rename a temporary Scoop");
      } else {
        showRenameScoop = !showRenameScoop;
      }
    }
  };

  const onEscape = (event: KeyboardEvent) => {
    if (event.key === "Escape" && (showCmdPalette || showRenameScoop || showHelp)) {
      switch (true) {
        case showCmdPalette:
          showCmdPalette = false;
          break;
        case showRenameScoop:
          showRenameScoop = false;
          break;
        case showHelp:
          showHelp = false;
          break;
      }
    }
  };

  const switchRequest = async (event: KeyboardEvent) => {
    if (event.ctrlKey && event.key > "0" && event.key <= "9") {
      const n: number = Number(event.key);

      if (n > appState.allScoops.length) return;

      console.log(`Switch Request Fired! Ctrl + ${n}`);

      const ok = await onSwitchRequest();
      if (ok) {
        appState.currentScoop = appState.allScoops[n - 1];
        appState.response = appState.currentScoop.response ?? "";
        appState.queryParams = appState.currentScoop.request.query_params;
      }
    }
  };

  onMount(() => {
    onRspMsg = Events.On("respMsg", async (event: any) => {
      const s = event.data as Scoop;

      // want to use reactive plain objects in UI since Svelete reactivity doesnt like classes
      appState.response = s.response;
      appState.url = s.request.url;
      persistFormToRequest(appState.currentScoop);

      loading = false;
    });

    onErrMsg = Events.On("errMsg", async (event: any) => {
      console.error(event.data);
      toast.error(event.data);
      loading = false;
    });

    document.addEventListener("keydown", openCmdPalette);
    document.addEventListener("keydown", renameScoop);
    document.addEventListener("keydown", onEscape);
    document.addEventListener("keydown", hideReqParams);
    document.addEventListener("keydown", switchRequest);
  });

  // cleanup events on destroy
  onDestroy(() => {
    onRspMsg?.();
    onErrMsg?.();

    document.removeEventListener("keydown", openCmdPalette);
    document.removeEventListener("keydown", renameScoop);
    document.removeEventListener("keydown", onEscape);
    document.removeEventListener("keydown", hideReqParams);
    document.removeEventListener("keydown", switchRequest);
  });
</script>

<Toaster />

<div class="relative flex min-h-screen min-w-screen flex-col items-center justify-center gap-5">
  <!-->App Title<-->
  <p class="text-md text-green-500">Scoop v1.0</p>

  <!-->Outer Card<-->
  <div class="border-border flex h-[90vh] w-[90vw] flex-col gap-10 rounded-sm border p-10 pb-0">
    <!-->Request Section<-->
    <div class="flex-rows flex w-full gap-8">
      <div class="flex min-w-0 flex-row items-center gap-2">
        <p>Method:</p>
        <div class="relative w-full max-w-[8rem] min-w-0">
          <select
            bind:value={appState.method}
            class={`border-border bg-accent h-9 w-full min-w-0
            appearance-none rounded-sm border px-3 pr-8
            ${methodColor(appState.method)}
            focus:ring-offset-background focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none`}
          >
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="PATCH">PATCH</option>
            <option value="DELETE">DELETE</option>
          </select>

          <!-->Caret Icon<-->
          <svg
            class="pointer-events-none absolute top-1/2 right-2 h-4 w-4 -translate-y-1/2 opacity-70"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M5.23 7.21a.75.75 0 011.06.02L10 10.94l3.71-3.71a.75.75 0 111.06 1.06l-4.24 4.24a.75.75 0 01-1.06 0L5.21 8.29a.75.75 0 01.02-1.08z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
      </div>
      <div class="flex min-w-0 flex-1 flex-row items-center gap-2">
        <p>URL:</p>
        <input
          class="focus:ring-offset-background w-full min-w-0 text-green-300 focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
          placeholder="https:// ..."
          bind:value={appState.url}
        />
      </div>
      <button
        class="border-border bg-accent text-foreground focus:ring-offset-background ml-auto inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
        disabled={loading}
        onclick={() => {
          onSend(appState.method, appState.url);
        }}
      >
        Send
      </button>
    </div>

    <!-->Request Headers, Query Params, and Request Body<-->
    {#if !reqParamsHidden}
      <div class="min-w-0 overflow-x-auto">
        <div class="flex flex-nowrap gap-10">
          <!-->Request Headers<-->
          <div
            class="border-border min-h-[25vh] min-w-[25rem] flex-1 shrink-0 rounded-sm border p-2"
          >
            <div class="flex w-full flex-row justify-between px-3">
              <p class="text-sm underline underline-offset-3">Headers</p>

              <!--Radio button group for text editing option-->
              <div class="flex items-center gap-4">
                <label for="header-raw" class="inline-flex items-center gap-2 text-sm">
                  <input
                    type="radio"
                    id="header-raw"
                    value="raw"
                    class="align-middle"
                    bind:group={headerTType}
                  />
                  raw
                </label>
                <label for="header-key-value" class="inline-flex items-center gap-2 text-sm">
                  <input
                    type="radio"
                    id="header-key-value"
                    value="key-value"
                    class="align-middle"
                    bind:group={headerTType}
                  />
                  key-value
                </label>
              </div>
            </div>
            {#if headerTType === "raw"}
              <RawInput bind:content={appState.headers} inputMode={"isHeader"} />
            {:else}
              <KvInput bind:content={appState.headers} inputMode={"isHeader"} />
            {/if}
          </div>

          <!-->Query Parameters<-->
          <div class="border-border min-h-[25vh] min-w-[25rem] flex-1 rounded-sm border p-2">
            <div class="flex w-full flex-row justify-between px-3">
              <p class="text-sm underline underline-offset-3">Parameters</p>

              <!--Radio button group for text editing option-->
              <div class="flex items-center gap-4">
                <label for="qp-raw" class="inline-flex items-center gap-2 text-sm">
                  <input
                    type="radio"
                    id="qp-raw"
                    value="raw"
                    class="align-middle"
                    bind:group={qParamTType}
                  />
                  raw
                </label>
                <label for="qp-key-value" class="inline-flex items-center gap-2 text-sm">
                  <input
                    type="radio"
                    id="qp-key-value"
                    value="key-value"
                    class="align-middle"
                    bind:group={qParamTType}
                  />
                  key-value
                </label>
              </div>
            </div>
            {#if qParamTType === "raw"}
              <RawInput bind:content={appState.queryParams} inputMode={"isQParam"} />
            {:else}
              <KvInput bind:content={appState.queryParams} inputMode={"isQParam"} />
            {/if}
          </div>

          <!-->Request Body<-->
          <div class="border-border min-h-[25vh] min-w-[25rem] flex-1 rounded-sm border p-2">
            <div class="flex w-full flex-row justify-between px-3">
              <p class="text-sm underline underline-offset-3">Body</p>

              <!--Radio button group for text editing option-->
              <div class="flex items-center gap-4">
                <label for="body-json" class="inline-flex items-center gap-2 text-sm">
                  <input
                    type="radio"
                    id="body-json"
                    value="json"
                    class="align-middle"
                    bind:group={bodyTType}
                  />
                  JSON
                </label>
              </div>
            </div>
            <!-->TODO: Would like this to be Monaco or CodeMirror Editor<-->
            <BodyInput />
          </div>
        </div>
      </div>
    {/if}

    <!-->Response Section<-->
    <div class="border-border bg-accent flex min-h-0 w-full flex-1 flex-col rounded-sm border p-2">
      <!-->Response Title Section<-->
      <div class="mb-3 flex flex-row gap-2">
        <p class="px-3 text-sm underline underline-offset-3">Response</p>
        {#if appState.response.body !== ""}
          <p class="border-border border px-2 text-sm">{appState.response.status}</p>
          <p class="border-border border px-2 text-sm">{`${appState.response.size}B`}</p>
          <p class="border-border border px-2 text-sm">{appState.response.duration} ms</p>
        {:else if loading}
          <DotSpinner />
        {/if}
        <button
          class="border-border mr-4 ml-auto inline-flex h-5 w-5 items-center justify-center rounded-sm border text-2xl hover:bg-green-400 hover:text-black focus:outline-none"
          onclick={() => {
            reqParamsHidden = !reqParamsHidden;
          }}
        >
          {reqParamsHidden === false ? "+" : "-"}
        </button>
      </div>
      <ResponseViewer
        value={appState.response.body ?? ""}
        contentType={appState.response.content_type ?? ""}
      />
    </div>
    <div class="-mx-10 h-8 items-center rounded-b-sm bg-green-950/30">
      <div class="flex h-full flex-row items-center gap-5 px-10 text-sm text-green-500/90">
        <div class="flex flex-row gap-2">
          <Package
            class={appState.currentCollection.name === "temp" ? `text-blue-500/90` : ``}
            size={20}
          />
          <p
            class={appState.currentCollection.name === "temp"
              ? `text-blue-500/90`
              : `text-green-400`}
          >
            {appState.currentCollection.name}
          </p>
        </div>

        {#if appState.currentScoop.name !== "temp"}
          <!-->TODO: add some logic to handle overflow<-->
          {#each appState.allScoops as scoop, i}
            <div class="flex flex-row gap-1">
              <p
                class={appState.currentScoop.name === scoop.name
                  ? `text-blue-500`
                  : `text-green-400`}
              >
                [{i + 1}]
              </p>
              <p
                class={appState.currentScoop.name === scoop.name
                  ? `text-blue-500`
                  : `text-green-400`}
              >
                {scoop.name}
              </p>
            </div>
          {/each}
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
  </div>

  <!-->Handful of Different Overlays<-->
  {#if showCmdPalette || showRenameScoop || showHelp || showCurl}
    <div
      class="fixed inset-0 z-100 flex min-h-screen min-w-screen items-center justify-center p-3 sm:p-8"
      aria-modal="true"
      role="dialog"
    >
      <!--Backdrop as button for click away-->
      <button
        class="absolute inset-0 bg-black/60"
        aria-hidden="true"
        onclick={() => {
          switch (true) {
            case showCmdPalette:
              showCmdPalette = false;
              break;
            case showRenameScoop:
              showRenameScoop = false;
              break;
            case showHelp:
              showHelp = false;
              break;
            case showCurl:
              showCurl = false;
              break;
          }
        }}
      ></button>
      {#if showCmdPalette}
        <!--CmdPalette-->
        <div class=" relative z-101 w-full max-w-xl shadow-lg">
          <CmdPalette />
        </div>
      {:else if showRenameScoop}
        <!--Rename Scoop-->
        <div class=" relative z-101 w-full max-w-xl shadow-lg">
          <RenameScoop bind:showRenameScoop />
        </div>
      {:else if showHelp}
        <!--Keybindings Help Component-->
        <div class=" relative z-101 w-full max-w-xl shadow-lg">
          <HelpKeybindings />
        </div>
      {:else if showCurl}
        <!--Generated CURL Command-->
        <div class=" relative z-101 w-full max-w-xl shadow-lg">
          <GeneratedCurl />
        </div>
      {/if}
    </div>
  {/if}
</div>
