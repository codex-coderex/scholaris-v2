<script lang="ts">
  import type { AppMode } from '$lib/types/schema'

  type Props = {
    page: number
    totalPages: number
    total: number
    mode?: AppMode
    busy?: boolean
    onPrevious: () => void
    onNext: () => void
    onGoToPage?: (pageNumber: number) => void
    onJump?: (event: KeyboardEvent) => void
  }

  let { page, totalPages, total, mode = 'light', busy = false, onPrevious, onNext, onGoToPage, onJump }: Props = $props()
</script>

<div class="ssis-pagination-bar text-sm text-(--text-muted)">
  <span>
    <span class="font-medium text-(--text-primary)">{total.toLocaleString()}</span> total
  </span>

  <div class="ssis-pagination-controls">
    <button
      type="button"
      disabled={busy || page === 1}
      onclick={(e) => {
        ;(e.currentTarget as HTMLButtonElement).blur()
        onPrevious()
      }}
      class="ssis-pagination-nav"
      aria-label="Previous page"
      title="Previous page"
    >
      ←
    </button>

    <button
      type="button"
      disabled={busy}
      class="ssis-pagination-current"
      onclick={() => onGoToPage?.(page)}
      onkeydown={onJump}
      aria-label={`Current page ${page} of ${totalPages}`}
      title={`Page ${page} of ${totalPages}`}
    >
      {page}
      {#if busy}
        <svg class="ml-1.5 h-3 w-3 animate-spin opacity-80" fill="none" viewBox="0 0 24 24" aria-hidden="true">
          <circle class="opacity-30" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3"></circle>
          <path class="opacity-100" fill="currentColor" d="M4 12a8 8 0 018-8v3a5 5 0 00-5 5H4z"></path>
        </svg>
      {/if}
    </button>

    <button
      type="button"
      disabled={busy || page >= totalPages}
      onclick={(e) => {
        ;(e.currentTarget as HTMLButtonElement).blur()
        onNext()
      }}
      class="ssis-pagination-nav"
      aria-label="Next page"
      title="Next page"
    >
      →
    </button>
  </div>
</div>