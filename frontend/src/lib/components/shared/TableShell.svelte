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
    <td colspan={colspan} class="py-8 text-center text-sm text-slate-500">
      {loadingMessage}
    </td>
  </tr>
{:else if !hasData}
  <tr>
    <td colspan={colspan} class="py-8 text-center text-sm text-slate-500">
      {emptyMessage}
    </td>
  </tr>
{:else}
  {@render children?.()}
{/if}
