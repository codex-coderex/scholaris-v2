<script lang="ts">
  import type { AppMode } from '$lib/types/schema'

  type Props = {
    mode: AppMode
    onToggle: () => void
  }

  let { mode, onToggle }: Props = $props()
</script>

<button
  type="button"
  class="ssis-mode-switch"
  aria-label="Toggle theme"
  aria-pressed={mode === 'dark'}
  onclick={onToggle}
>
  <span class="ssis-mode-label">{mode === 'light' ? 'Light Mode' : 'Dark Mode'}</span>
  <span class="ssis-mode-track {mode === 'dark' ? 'is-dark' : ''}">
    <span class="ssis-mode-thumb {mode === 'dark' ? 'is-dark' : ''}"></span>
  </span>
</button>
