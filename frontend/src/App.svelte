<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import { onDestroy, onMount } from "svelte";
  import { Method, KV, Backend, Scoop, Response, Collection, Request } from "../bindings/changeme";
  import KvInput from "$lib/components/kv-input.svelte";
  import RawInput from "$lib/components/raw-input.svelte";
  import DotSpinner from "$lib/components/dot-spinner.svelte";
  import ResponseViewer from "$lib/components/response-viewer.svelte";
  import CmdPalette from "$lib/components/command-palette.svelte";
  import { Toaster } from "$lib/components/ui/sonner/index.js";
  import { toast } from "svelte-sonner";
  import Package from "@lucide/svelte/icons/package";

  // events emitted from backend
  let onRspMsg: undefined | (() => void);
  let onErrMsg: undefined | (() => void);

  let showCmdPalette: boolean = $state(false);
  let reqParamsHidden: boolean = $state(false);

  // set default to temp
  let currentCollection: Collection = $state(new Collection({ name: "temp" }));
  let currentRequest: Request = $state(new Request({ name: "temp" }));

  let scoop: Scoop | null = $state(null);
  let method: Method = $state(Method.Empty);
  let url: string = $state("");

  let headers: KV[] = $state([]);
  let queryParams: KV[] = $state([]);

  // need this response value, since Svelte reactivity does not play nice with mutated class instances
  let response: Response | undefined = $state();
  let loading: boolean = $state(false);

  type TType = "raw" | "key-value" | "json";

  // TODO test swapping out the RawInput component for a monaco-editor
  let headerTType: TType = $state("raw");
  let headerRawContent: string = $state("");

  let qParamTType: TType = $state("raw");
  let qParamRawContent: string = $state("");

  // TODO support body in request on backend
  let bodyTType: TType = $state("json");
  let bodyRawContent: string = $state("");

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

  async function onSend(method: Method, url: string) {
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

    response = undefined;
    loading = true;

    rawToHeaders();
    rawToQParams();

    try {
      scoop = await Backend.ModelIntializer(method, url, headers, queryParams);
      await Backend.SubmitRequest(scoop);
    } catch (error) {
      console.error(error);
      loading = false;
    }
  }

  function rawToHeaders() {
    const r: string[] = headerRawContent.split("\n");
    headers = [];

    for (const row of r) {
      if (row === "") return;

      const idx: number = row.indexOf(":");

      const key: string = (idx === -1 ? row : row.slice(0, idx)).trim();
      const val: string = (idx === -1 ? "" : row.slice(idx + 1)).trim();

      const newRow: KV = { key: key, value: val };
      console.log(`Headers: ${newRow.key}:${newRow.value}`);
      headers.push(newRow);
    }
  }

  function rawToQParams() {
    const r: string[] = qParamRawContent.split("\n");
    queryParams = [];

    for (const row of r) {
      if (row === "") return;

      const idx: number = row.indexOf("=");

      const key: string = (idx === -1 ? row : row.slice(0, idx)).trim();
      const val: string = (idx === -1 ? "" : row.slice(idx + 1)).trim();

      const newRow: KV = { key: key, value: val };
      console.log(`Query Param: ${newRow.key}=${newRow.value}`);
      queryParams.push(newRow);
    }
  }

  const openCmdPalette = (event: KeyboardEvent) => {
    if (event.ctrlKey && event.shiftKey && (event.key === "P" || event.code === "KeyP")) {
      showCmdPalette = true;
    } else if (event.key === "Escape" && showCmdPalette) {
      showCmdPalette = false;
    }
  };

  const hideReqParams = (event: KeyboardEvent) => {
    if (event.ctrlKey && (event.key === "E" || event.code === "KeyE")) {
      reqParamsHidden = !reqParamsHidden;
    }
  };

  onMount(() => {
    onRspMsg = Events.On("respMsg", async (event: any) => {
      scoop = event.data as Scoop;

      // want to use reactive plain objects in UI since Svelete reactivity doesnt like classes
      response = scoop.response;
      url = scoop?.request.url ?? "";

      loading = false;
    });

    onErrMsg = Events.On("errMsg", async (event: any) => {
      console.error(event.data);
      toast.error(event.data);
      loading = false;
    });

    document.addEventListener("keydown", openCmdPalette);
    document.addEventListener("keydown", hideReqParams);
  });

  // cleanup events on destroy
  onDestroy(() => {
    onRspMsg?.();
    onErrMsg?.();

    document.removeEventListener("keydown", openCmdPalette);
    document.removeEventListener("keydown", hideReqParams);
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
            bind:value={method}
            class={`border-border bg-accent h-9 w-full min-w-0 cursor-pointer
            appearance-none rounded-sm border px-3 pr-8
            ${methodColor(method)}
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
          placeholder="https://google.com"
          bind:value={url}
        />
      </div>
      <button
        class="border-border bg-accent text-foreground focus:ring-offset-background ml-auto inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:cursor-pointer hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
        disabled={loading}
        onclick={() => {
          onSend(method, url);
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
              <RawInput bind:content={headerRawContent} />
            {:else}
              <KvInput bind:rawContent={headerRawContent} inputMode={"isHeader"} />
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
              <RawInput bind:content={qParamRawContent} />
            {:else}
              <KvInput bind:rawContent={qParamRawContent} inputMode={"isQParam"} />
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
            <RawInput bind:content={bodyRawContent} />
          </div>
        </div>
      </div>
    {/if}

    <!-->Response Section<-->
    <div class="border-border bg-accent flex min-h-0 w-full flex-1 flex-col rounded-sm border p-2">
      <!-->Response Title Section<-->
      <div class="mb-3 flex flex-row gap-2">
        <p class="px-3 text-sm underline underline-offset-3">Response</p>
        {#if response}
          <p class="border-border border px-2 text-sm">{response.status}</p>
          <p class="border-border border px-2 text-sm">{response.content_type}</p>
          <p class="border-border border px-2 text-sm">{response?.duration} ms</p>
        {:else if loading}
          <DotSpinner />
        {/if}
        <button
          class="border-border mr-4 ml-auto inline-flex h-5 w-5 items-center justify-center rounded-sm border text-2xl hover:cursor-pointer hover:bg-green-400 hover:text-black focus:outline-none"
          onclick={() => {
            reqParamsHidden = !reqParamsHidden;
          }}
        >
          {reqParamsHidden === false ? "+" : "-"}
        </button>
      </div>
      <ResponseViewer value={response?.body ?? ""} contentType={response?.content_type ?? ""} />
    </div>
    <div class="-mx-10 h-8 items-center rounded-b-sm bg-green-950/30">
      <div class="flex h-full flex-row items-center gap-5 px-10 text-sm text-green-500/90">
        <div class="flex flex-row gap-2">
          <Package class={currentCollection.name === "temp" ? `text-blue-500/90` : ``} size={20} />
          <p class={currentCollection.name === "temp" ? `text-blue-500/90` : ``}>
            {currentCollection.name}
          </p>
        </div>

        <!-->I want this to resemble tmux sessions<-->
        <!-->TODO will need to display all requests in the collection not just one<-->
        {#if currentRequest.name !== "temp"}
          <p>/</p>
          <p>{currentRequest.name}</p>
        {/if}
      </div>
    </div>
  </div>

  <!-->CmdPalette Overlay<-->
  {#if showCmdPalette}
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
          showCmdPalette = false;
        }}
      ></button>

      <!--CmdPalette-->
      <div class=" relative z-101 w-full max-w-xl shadow-lg">
        <CmdPalette bind:collection={currentCollection} bind:request={currentRequest} />
      </div>
    </div>
  {/if}
</div>
