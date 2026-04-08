<script lang="ts">
  import type { Snippet } from 'svelte'

  type Props = {
    loading: boolean
    hasData: boolean
    colspan: number
    loadingMessage?: string
    emptyMessage?: string
    children?: Snippet
  }

  let {
    loading,
    hasData,
    colspan,
    loadingMessage = 'Loading...',
    emptyMessage = 'No data found',
    children
  }: Props = $props()
</script>

{#if loading}
  <tr>
    <td {colspan} class="py-12 text-center">
      <div class="inline-flex items-center gap-2 text-sm text-slate-500">
        <svg class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-20" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-90 text-cyan-500" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
        {loadingMessage}
      </div>
    </td>
  </tr>
{:else if !hasData}
  <tr>
    <td {colspan} class="py-12 text-center text-sm text-slate-500">
      {emptyMessage}
    </td>
  </tr>
{:else}
  {@render children?.()}
{/if}