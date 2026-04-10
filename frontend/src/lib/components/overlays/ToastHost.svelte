<script lang="ts">
  import { fly } from 'svelte/transition'

  import { removeToast, toasts } from '$lib/stores/toast'
</script>

<div class="toast toast-end fixed bottom-6 right-6 z-50 space-y-3">
  {#each $toasts as toast (toast.id)}
    <div
      class={`alert shadow-lg ${toast.kind === 'error' ? 'alert-error' : toast.kind === 'info' ? 'alert-info' : 'alert-success'}`}
      transition:fly={{ x: 20, duration: 160 }}
    >
      <span class="text-sm font-medium">{toast.message}</span>
      <button type="button" class="btn btn-ghost btn-xs" onclick={() => removeToast(toast.id)}>✕</button>
    </div>
  {/each}
</div>