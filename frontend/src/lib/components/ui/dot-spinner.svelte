<script lang="ts">
  import { onMount, onDestroy } from "svelte";

  const BRAILLE_MINI_DOT = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];

  export let label: string = "";
  export let intervalMs: number = 80;
  export let running: boolean = true;
  export let frames: string[] = BRAILLE_MINI_DOT;

  let i = 0;
  let timer: ReturnType<typeof setInterval> | null = null;

  function start() {
    stop();
    timer = setInterval(() => {
      i = (i + 1) % frames.length;
    }, intervalMs);
  }

  function stop() {
    if (timer) clearInterval(timer);
    timer = null;
  }

  onMount(() => {
    if (running) start();

    const mql = window.matchMedia?.("(prefers-reduced-motion: reduce)");
    const handle = () => {
      if (mql?.matches) stop();
      else if (running) start();
    };

    mql?.addEventListener?.("change", handle);
    handle();

    return () => mql?.removeEventListener?.("change", handle);
  });

  onDestroy(stop);
</script>

<span class="inline-flex items-center gap-2 font-mono" aria-live="polite" aria-busy={running}>
  <span aria-hidden="true">{frames[i]}</span>
  {#if label}<span>{label}</span>{/if}
</span>
