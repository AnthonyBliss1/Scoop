<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import KvInput from "$lib/components/kv-input.svelte";
  import RawInput from "$lib/components/raw-input.svelte";
  import { Method, KV, Backend, Scoop, Response } from "../bindings/changeme";
  import DotSpinner from "$lib/components/dot-spinner.svelte";
  import ResponseViewer from "$lib/components/response-viewer.svelte";
  import { onDestroy, onMount } from "svelte";
  import { Toaster } from "$lib/components/ui/sonner/index.js";
  import { toast } from "svelte-sonner";

  // events emitted from backend
  let onRspMsg: undefined | (() => void);
  let onErrMsg: undefined | (() => void);

  let scoop: Scoop | null = $state(null);
  let method: Method = $state(Method.Empty);
  let url: string = $state("");

  let headers: KV[] = $state([]);
  let queryParams: KV[] = $state([]);

  let response: Response | undefined = $state();

  let reqDuration: number | undefined = $state();
  let loading: boolean = $state(false);

  type TType = "raw" | "key-value" | "json";

  let reqParamsHidden = $state(false);

  // TODO test swapping out the RawInput component for a monaco-editor
  let headerTType: TType = $state("raw");
  let headerRawContent: string = $state("");

  let qParamTType: TType = $state("raw");
  let qParamRawContent: string = $state("");

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

  async function onSubmit(method: Method, url: string) {
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
    reqDuration = undefined;
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

  onMount(() => {
    onRspMsg = Events.On("respMsg", async (event: any) => {
      console.log("Reveived Response!");
      response = event.data as Response;
      reqDuration = response.duration;
      loading = false;
    });

    onErrMsg = Events.On("errMsg", async (event: any) => {
      console.error(event.data);
      toast.error(event.data);
      loading = false;
    });
  });

  // cleanup events on destroy
  onDestroy(() => {
    onRspMsg?.();
    onErrMsg?.();
  });
</script>

<Toaster />

<div class="flex min-h-screen min-w-screen flex-col items-center justify-center gap-5">
  <!-->App Title<-->
  <p class="text-lg text-green-500">Scoop v1.0</p>

  <!-->Outer Card<-->
  <div class="border-border flex h-[90vh] w-[90vw] flex-col gap-10 rounded-sm border p-10">
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
          onSubmit(method, url);
        }}
      >
        Submit
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
        {#if reqDuration}
          {#if response}
            <p class="border-border border px-2 text-sm">{response.status}</p>
            <p class="border-border border px-2 text-sm">{response.content_type}</p>
          {/if}

          <p class="border-border border px-2 text-sm">{reqDuration} ms</p>
        {:else if loading}
          <DotSpinner />
        {/if}
        <button
          class="border-border mr-4 ml-auto inline-flex h-5 w-5 items-center justify-center rounded-sm border text-2xl hover:bg-green-400 hover:text-black"
          onclick={() => {
            reqParamsHidden = !reqParamsHidden;
          }}
        >
          {reqParamsHidden === false ? "+" : "-"}
        </button>
      </div>
      <ResponseViewer value={response?.body ?? ""} contentType={response?.content_type ?? ""} />
    </div>
  </div>
</div>
