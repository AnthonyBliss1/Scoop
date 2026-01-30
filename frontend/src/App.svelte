<script lang="ts">
  import { prefersReducedMotion } from "svelte/motion";
  import { Method, KV, App, Scoop, Response } from "../bindings/changeme";

  let scoop: Scoop | null = $state(null);
  let method: Method = $state(Method.Empty);
  let url = $state("");
  let headers = $state(new Map<string, string>([]));
  let queryParams = $state(new Map<string, string>([]));
  let response: Response | undefined = $state<Response>();

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
  <div
    class="border-border flex max-h-[90vh] min-h-[90vh] max-w-[90vw] min-w-[90vw] flex-col gap-10 rounded-sm border p-10"
  >
    <!-->Request Section<-->
    <div class="flex-rows flex w-full gap-20">
      <div class="flex flex-row items-center gap-2">
        <p>Method:</p>
        <input
          placeholder="GET"
          list="reqMethods"
          bind:value={method}
          class={methodColor(method)}
        />
        <datalist id="reqMethods">
          <option value="GET"></option>
          <option value="POST"></option>
          <option value="PUT"></option>
          <option value="PATCH"></option>
          <option value="DELETE"></option>
        </datalist>
      </div>
      <div class="flex flex-row items-center gap-2">
        <p>URL:</p>
        <input placeholder="https://google.com" bind:value={url} />
      </div>
      <button
        class="border-border ml-auto items-center rounded-sm border px-3 py-1 hover:cursor-pointer hover:bg-green-400 hover:text-black"
        onclick={void onSubmit(method, url, [])}>Submit</button
      >
    </div>

    <!-->Request Headers and Query Params<-->
    <div class="flex-rows flex gap-10">
      <div class="border-border min-h-[25vh] w-full flex-1 rounded-sm border p-2">
        <p class="text-sm underline underline-offset-3">Request Headers</p>
      </div>
      <div class="border-border min-h-[25vh] w-full flex-1 rounded-sm border p-2">
        <p class="text-sm underline underline-offset-3">Query Parameters</p>
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
