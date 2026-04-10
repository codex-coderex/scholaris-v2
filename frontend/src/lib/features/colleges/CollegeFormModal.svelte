<script lang="ts">
  import Modal from '$lib/components/overlays/Modal.svelte'

  type Props = {
    open: boolean
    title: string
    code: string
    name: string
    formError: string
    onCodeInput: (value: string) => void
    onNameInput: (value: string) => void
    onSave: () => void
    onClose: () => void
  }

  let {
    open,
    title,
    code,
    name,
    formError,
    onCodeInput,
    onNameInput,
    onSave,
    onClose
  }: Props = $props()

  let isEditMode = $derived(title.toLowerCase().includes('edit'))
</script>

<Modal {open} {title} onClose={onClose} size="md">
  <div class="space-y-4">
    <label class="form-control w-full gap-2">
      <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">College Code</span>
      <input
        class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
        value={code}
        disabled={isEditMode}
        placeholder="CCS"
        oninput={(event) => {
          onCodeInput((event.currentTarget as HTMLInputElement).value)
        }}
      />
    </label>

    <label class="form-control w-full gap-2">
      <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">College Name</span>
      <input
        class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
        value={name}
        placeholder="College of Computer Studies"
        oninput={(event) => {
          onNameInput((event.currentTarget as HTMLInputElement).value)
        }}
      />
    </label>

    {#if formError}
      <div class="alert alert-error rounded-box border border-(--danger-dim) bg-(--danger-dim) text-(--danger)">{formError}</div>
    {/if}

    <div class="modal-action justify-start gap-3">
      <button type="button" class="btn btn-primary" onclick={onSave}>Save</button>
      <button type="button" class="btn btn-ghost" onclick={onClose}>Cancel</button>
    </div>
  </div>
</Modal>
