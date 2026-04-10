<script lang="ts">
  import { scale } from 'svelte/transition'
  import type { Snippet } from 'svelte'

  type Props = {
    open: boolean
    title: string
    onClose: () => void
    size?: 'sm' | 'md' | 'lg' | 'xl'
    closable?: boolean
    children?: Snippet
  }

  let { open, title, onClose, size = 'md', closable = false, children }: Props = $props()

  function portal(node: HTMLDivElement) {
    if (typeof document === 'undefined') {
      return
    }

    document.body.appendChild(node)

    return {
      destroy() {
        if (node.parentNode === document.body) {
          document.body.removeChild(node)
        }
      }
    }
  }

  $effect(() => {
    if (typeof document === 'undefined') {
      return
    }

    if (open) {
      document.body.dataset.ssisModalOpen = 'true'
    } else {
      delete document.body.dataset.ssisModalOpen
    }

    return () => {
      delete document.body.dataset.ssisModalOpen
    }
  })

  let sizeClass = $derived(
    {
      sm: 'max-w-md',
      md: 'max-w-xl',
      lg: 'max-w-3xl',
      xl: 'max-w-5xl'
    }[size]
  )
</script>

{#if open}
  <div class="modal modal-open ssis-modal-portal" use:portal>
    <div
      class="modal-box ssis-modal-box relative w-full overflow-visible {sizeClass} p-0"
      transition:scale={{ start: 0.98, duration: 140 }}
    >
      {#if closable}
        <button
          type="button"
          class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4 text-(--text-muted) hover:bg-(--bg-hover) hover:text-(--text-primary)"
          onclick={onClose}
          aria-label="Close dialog"
        >
          ✕
        </button>
      {/if}

      <div class="ssis-modal-header px-7 pt-6 pb-2">
        <h2 class="ssis-modal-title font-lora text-(--text-primary)">{title}</h2>
      </div>

      <div class="ssis-modal-content overflow-visible px-7 pb-6 pt-3">
        {@render children?.()}
      </div>
    </div>

    <div class="modal-backdrop ssis-modal-backdrop" onclick={onClose} aria-hidden="true"></div>
  </div>
{/if}