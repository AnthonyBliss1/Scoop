<script lang="ts">
  let { cmd = $bindable("Create New Request"), request = $bindable("") } = $props<{
    cmd: any;
    request: string;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);

  let newRequest: string = $state("");

  function createRequest() {
    if (newRequest === "") return;

    request = newRequest;
    cmd = null;
  }

  $effect(() => {
    if (cmd === "Create New Request") {
      inputEl?.focus();
    }
  });
</script>

<div class="border-border bg-background flex flex-col gap-5 rounded-sm border p-5">
  <p class="flex items-center justify-center">Create New Request</p>
  <input
    class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none md:min-w-[450px]"
    bind:value={newRequest}
    bind:this={inputEl}
    placeholder="Enter Request Name..."
  />

  <div class="flex w-full flex-row items-center justify-center gap-10">
    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:cursor-pointer hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={createRequest}>Create</button
    >

    <button
      class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:cursor-pointer hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
      onclick={() => {
        cmd = null;
      }}>Cancel</button
    >
  </div>
</div>
