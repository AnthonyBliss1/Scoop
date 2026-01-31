<script lang="ts">
  import KvInput from "$lib/components/kv-input.svelte";
  import RawInput from "$lib/components/raw-input.svelte";
  import { Method, KV, App, Scoop, Response } from "../bindings/changeme";

  let scoop: Scoop | null = $state(null);
  let method: Method = $state(Method.Empty);
  let url = $state("");
  let headers = $state(new Map<string, string>([]));
  let queryParams = $state(new Map<string, string>([]));
  let response: Response | undefined = $state<Response>();

  type TType = "raw" | "key-value" | "json";

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
        return "";
    }
  }

  async function onSubmit(method: Method, url: string, headers: KV[]) {
    try {
      scoop = await App.ModelIntializer(method, url, headers);
      response = await App.SubmitRequest(scoop);
    } catch (error) {
      console.error(error);
    }
  }
</script>

<div class="flex min-h-screen min-w-screen flex-col items-center justify-center gap-5">
  <!-->App Title<-->
  <p class="text-lg text-green-500">Scoop v1.0</p>

  <!-->Outer Card<-->
  <div class="border-border flex h-[90vh] w-[90vw] flex-col gap-10 rounded-sm border p-10">
    <!-->Request Section<-->
    <div class="flex-rows flex w-full gap-4">
      <div class="flex min-w-0 flex-row items-center gap-2">
        <p>Method:</p>
        <input
          placeholder="GET"
          list="reqMethods"
          bind:value={method}
          class={`w-full max-w-[8rem] min-w-0 ${methodColor(method)}`}
        />
        <datalist id="reqMethods">
          <option value="GET"></option>
          <option value="POST"></option>
          <option value="PUT"></option>
          <option value="PATCH"></option>
          <option value="DELETE"></option>
        </datalist>
      </div>
      <div class="flex min-w-0 flex-1 flex-row items-center gap-2">
        <p>URL:</p>
        <input class="w-full min-w-0" placeholder="https://google.com" bind:value={url} />
      </div>
      <button
        class="border-border ml-auto items-center rounded-sm border px-3 py-1 hover:cursor-pointer hover:bg-green-400 hover:text-black"
        onclick={void onSubmit(method, url, [])}>Submit</button
      >
    </div>

    <!-->Request Headers, Query Params, and Request Body<-->
    <div class="min-w-0 overflow-x-auto">
      <div class="flex flex-nowrap gap-10">
        <!-->Request Headers<-->
        <div class="border-border min-h-[25vh] min-w-[25rem] flex-1 shrink-0 rounded-sm border p-2">
          <div class="flex w-full flex-row justify-between px-3">
            <p class="text-sm underline underline-offset-3">Request Headers</p>

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
            <RawInput />
          {:else}
            <KvInput />
          {/if}
        </div>

        <!-->Query Parameters<-->
        <div class="border-border min-h-[25vh] min-w-[25rem] flex-1 rounded-sm border p-2">
          <div class="flex w-full flex-row justify-between px-3">
            <p class="text-sm underline underline-offset-3">Query Parameters</p>

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
            <RawInput />
          {:else}
            <KvInput />
          {/if}
        </div>

        <!-->Request Body<-->
        <div class="border-border min-h-[25vh] min-w-[25rem] flex-1 rounded-sm border p-2">
          <div class="flex w-full flex-row justify-between px-3">
            <p class="text-sm underline underline-offset-3">Request Body</p>

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
          <RawInput />
        </div>
      </div>
    </div>

    <!-->Response Section<-->
    <div class="border-border flex min-h-0 w-full flex-1 flex-col rounded-sm border p-2">
      <p class="text-sm underline underline-offset-3">Response</p>
      <pre class="min-h-0 flex-1 overflow-auto text-sm wrap-break-word whitespace-pre-wrap">
      {response ? JSON.stringify(response, null, 2) : ""}
      </pre>
    </div>
  </div>
</div>
