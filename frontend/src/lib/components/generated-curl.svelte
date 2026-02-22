<script lang="ts">
  let { curlCommand = $bindable<string>("") } = $props();
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex items-center justify-center">Generated cURL Command</p>
  <textarea
    class="focus:ring-offset-background bg-accent border-border h-8 w-full
    min-w-0 resize-none rounded-sm border px-2 text-green-300
    shadow-md focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none md:min-w-[450px]"
    placeholder="Generating cURL command ..."
    readonly
    wrap="off"
    bind:value={curlCommand}
  >
  </textarea>
</div>
