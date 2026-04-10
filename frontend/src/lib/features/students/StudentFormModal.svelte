<script lang="ts">
  import Modal from '$lib/components/overlays/Modal.svelte'
  import Combobox from '$lib/components/inputs/Combobox.svelte'

  import type { StudentGender } from '$lib/types/schema'

  type CollegeOption = {
    code: string
    name: string
  }

  type ProgramOption = {
    code: string
    name: string
    college_code?: string
  }

  type YearOption = {
    code: string
    name: string
  }

  type GenderOption = {
    code: StudentGender
    name: string
  }

  type Props = {
    open: boolean
    title: string
    editStudentId: string
    studentIdYear: string
    studentIdSeq: string
    firstName: string
    lastName: string
    year: string
    gender: StudentGender | ''
    collegeOptions: CollegeOption[]
    programOptions: ProgramOption[]
    selectedCollegeCode: string
    selectedCollegeLabel: string
    selectedProgramCode: string
    selectedProgramLabel: string
    formError: string
    onStudentIdYearInput: (value: string) => void
    onStudentIdSeqInput: (value: string) => void
    onFirstNameInput: (value: string) => void
    onLastNameInput: (value: string) => void
    onYearInput: (value: string) => void
    onGenderInput: (value: StudentGender | '') => void
    onSelectCollege: (option: CollegeOption | null) => void
    onSelectProgram: (option: ProgramOption | null) => void
    onSave: () => void
    onClose: () => void
  }

  let {
    open,
    title,
    editStudentId,
    studentIdYear,
    studentIdSeq,
    firstName,
    lastName,
    year,
    gender,
    collegeOptions,
    programOptions,
    selectedCollegeCode,
    selectedCollegeLabel,
    selectedProgramCode,
    selectedProgramLabel,
    formError,
    onStudentIdYearInput,
    onStudentIdSeqInput,
    onFirstNameInput,
    onLastNameInput,
    onYearInput,
    onGenderInput,
    onSelectCollege,
    onSelectProgram,
    onSave,
    onClose
  }: Props = $props()

  const yearOptions: YearOption[] = [
    { code: '1', name: '1st Year' },
    { code: '2', name: '2nd Year' },
    { code: '3', name: '3rd Year' },
    { code: '4', name: '4th Year' }
  ]

  const genderOptions: GenderOption[] = [
    { code: 'Male', name: 'Male' },
    { code: 'Female', name: 'Female' },
    { code: 'Other', name: 'Other' }
  ]

  let selectedYearLabel = $derived(yearOptions.find((option) => option.code === year)?.name ?? '')
  let selectedGenderLabel = $derived(genderOptions.find((option) => option.code === gender)?.name ?? '')
</script>

<Modal {open} {title} onClose={onClose} size="lg">
  <div class="space-y-4">
    <input type="hidden" value={editStudentId} />

    <div class="grid gap-4 sm:grid-cols-[120px_minmax(0,1fr)]">
      <label class="form-control w-full gap-2">
        <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">Student ID Year</span>
        <input
          class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
          value={studentIdYear}
          maxlength="4"
          placeholder="2026"
          oninput={(event) => {
            onStudentIdYearInput((event.currentTarget as HTMLInputElement).value)
          }}
        />
      </label>

      <label class="form-control w-full gap-2">
        <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">Student ID Sequence</span>
        <input
          class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
          value={studentIdSeq}
          maxlength="4"
          placeholder="0001"
          oninput={(event) => {
            onStudentIdSeqInput((event.currentTarget as HTMLInputElement).value)
          }}
        />
      </label>
    </div>

    <p class="text-xs text-(--text-muted)">Format: YYYY-NNNN</p>

    <div class="grid gap-4 sm:grid-cols-2">
      <label class="form-control w-full gap-2">
        <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">First Name</span>
        <input
          class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
          value={firstName}
          placeholder="First name"
          oninput={(event) => {
            onFirstNameInput((event.currentTarget as HTMLInputElement).value)
          }}
        />
      </label>

      <label class="form-control w-full gap-2">
        <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">Last Name</span>
        <input
          class="input input-bordered w-full bg-(--bg-base) border-(--border) text-(--text-primary) focus:border-(--accent) focus:outline-none"
          value={lastName}
          placeholder="Last name"
          oninput={(event) => {
            onLastNameInput((event.currentTarget as HTMLInputElement).value)
          }}
        />
      </label>
    </div>

    <div class="grid gap-4 sm:grid-cols-2">
      <Combobox
        label="College"
        placeholder="Type to search college..."
        options={collegeOptions}
        selectedCode={selectedCollegeCode}
        selectedLabel={selectedCollegeLabel}
        onSelect={onSelectCollege}
      />

      <Combobox
        label="Program"
        placeholder="Type to search program..."
        options={programOptions}
        selectedCode={selectedProgramCode}
        selectedLabel={selectedProgramLabel}
        onSelect={onSelectProgram}
      />
    </div>

    <div class="grid gap-4 sm:grid-cols-2">
      <Combobox
        label="Year"
        placeholder="Select year"
        options={yearOptions}
        selectedCode={year}
        selectedLabel={selectedYearLabel}
        onSelect={(option) => {
          onYearInput(option?.code ?? '')
        }}
      />

      <Combobox
        label="Gender"
        placeholder="Select gender"
        options={genderOptions}
        selectedCode={gender}
        selectedLabel={selectedGenderLabel}
        onSelect={(option) => {
          onGenderInput((option?.code as StudentGender | undefined) ?? '')
        }}
      />
    </div>

    {#if formError}
      <div class="alert alert-error rounded-box border border-(--danger-dim) bg-(--danger-dim) text-(--danger)">{formError}</div>
    {/if}

    <div class="modal-action justify-start gap-3">
      <button type="button" class="btn btn-primary" onclick={onSave}>Save</button>
      <button type="button" class="btn btn-ghost" onclick={onClose}>Cancel</button>
    </div>
  </div>
</Modal>
