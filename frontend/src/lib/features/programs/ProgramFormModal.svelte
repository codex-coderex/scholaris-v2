<script lang="ts">
  import Modal from '$lib/components/overlays/Modal.svelte'
  import Combobox from '$lib/components/inputs/Combobox.svelte'

  type Option = {
    code: string
    name: string
  }

  type Props = {
    open: boolean
    title: string
    code: string
    name: string
    colleges: Option[]
    selectedCollegeCode: string
    selectedCollegeLabel: string
    formError: string
    onCodeInput: (value: string) => void
    onNameInput: (value: string) => void
    onSelectCollege: (option: Option | null) => void
    onSave: () => void
    onClose: () => void
  }

  let {
    open,
    title,
    code,
    name,
    colleges,
    selectedCollegeCode,
    selectedCollegeLabel,
    formError,
    onCodeInput,
    onNameInput,
    onSelectCollege,
    onSave,
    onClose
  }: Props = $props()
</script>

<Modal {open} {title} onClose={onClose} size="lg">
  <div class="space-y-4">
    <label class="form-control w-full gap-2">
      <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">Program Code</span>
      <input
        class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
        value={code}
        placeholder="BSCS"
        oninput={(event) => {
          onCodeInput((event.currentTarget as HTMLInputElement).value)
        }}
      />
    </label>

    <label class="form-control w-full gap-2">
      <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">Program Name</span>
      <input
        class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
        value={name}
        placeholder="Bachelor of Science in Computer Science"
        oninput={(event) => {
          onNameInput((event.currentTarget as HTMLInputElement).value)
        }}
      />
    </label>

    <Combobox
      label="College"
      placeholder="Type to search college..."
      options={colleges}
      selectedCode={selectedCollegeCode}
      selectedLabel={selectedCollegeLabel}
      onSelect={onSelectCollege}
    />

    {#if formError}
      <div class="alert alert-error rounded-box border border-(--danger-dim) bg-(--danger-dim) text-(--danger)">{formError}</div>
    {/if}

    <div class="modal-action justify-start gap-3">
      <button type="button" class="btn btn-primary" onclick={onSave}>Save</button>
      <button type="button" class="btn btn-ghost" onclick={onClose}>Cancel</button>
    </div>
  </div>
</Modal>
