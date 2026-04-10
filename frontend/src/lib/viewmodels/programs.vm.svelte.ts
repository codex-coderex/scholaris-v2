import { College as CollegeModel, Program as ProgramModel } from '../../../bindings/scholaris-v2/internal/features/programs/models.js'
import { CreateProgram, DeleteProgram, GetPrograms, UpdateProgram } from '../../../bindings/scholaris-v2/appservice.js'

import { pushToast } from '$lib/stores/toast'
import { fetchCollegeOptions } from '$lib/services/lookups.service'
import { validateCode, validateLabel } from '$lib/shared/validation'
import type { CollegeRow, ProgramRow } from '$lib/types/schema'
import { TableQueryViewModel } from '$lib/viewmodels/table-query.vm.svelte'

type ProgramFormMode = 'create' | 'edit'

type CollegeOption = {
  code: string
  name: string
}

export class ProgramsViewModel {
  readonly noCollegeFilterCode = '__NO_COLLEGE__'

  query = new TableQueryViewModel('p.code')

  programs = $state<ProgramRow[]>([])
  colleges = $state<CollegeRow[]>([])
  loading = $state(false)
  isPageTransitioning = $state(false)
  error = $state('')

  requestVersion = 0

  formMode = $state<ProgramFormMode>('create')
  formOpen = $state(false)
  confirmOpen = $state(false)
  deleteTarget = $state<ProgramRow | null>(null)

  code = $state('')
  name = $state('')
  selectedCollegeCode = $state('')
  filterCollegeCode = $state('')
  formError = $state('')

  searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

  selectedCollege = $derived(this.colleges.find((college) => college.code === this.selectedCollegeCode) ?? null)
  selectedCollegeLabel = $derived(
    this.selectedCollege
      ? `${this.selectedCollege.code} — ${this.selectedCollege.name}`
      : this.formMode === 'edit' && !this.selectedCollegeCode
        ? '(No College)'
        : ''
  )

  openAdd() {
    this.formMode = 'create'
    this.code = ''
    this.name = ''
    this.selectedCollegeCode = ''
    this.formError = ''
    this.formOpen = true
  }

  openEdit(program: ProgramRow) {
    this.formMode = 'edit'
    this.code = program.code
    this.name = program.name
    this.selectedCollegeCode = program.college_code ?? program.college?.code ?? ''
    this.formError = ''
    this.formOpen = true
  }

  askDelete(program: ProgramRow) {
    this.deleteTarget = program
    this.confirmOpen = true
  }

  closeForm() {
    this.formOpen = false
    this.formError = ''
  }

  closeConfirm() {
    this.confirmOpen = false
    this.deleteTarget = null
  }

  async loadOptions() {
    try {
      this.colleges = await fetchCollegeOptions()
    } catch (caught: unknown) {
      this.colleges = []
      this.error = caught instanceof Error ? caught.message : 'Failed to load college options'
    }
  }

  async load() {
    const requestId = ++this.requestVersion
    const hasExistingRows = this.programs.length > 0

    this.loading = !hasExistingRows
    this.isPageTransitioning = hasExistingRows
    this.error = ''

    try {
      const [rows, total] = await GetPrograms(
        this.query.search,
        this.query.sortBy,
        this.query.order,
        this.query.page,
        this.query.pageSize,
        this.filterCollegeCode
      )

      if (requestId !== this.requestVersion) {
        return
      }

      this.programs = rows ?? []
      this.query.setTotal(total)
    } catch (caught: unknown) {
      if (requestId !== this.requestVersion) {
        return
      }

      this.error = caught instanceof Error ? caught.message : 'Failed to load programs'
    } finally {
      if (requestId === this.requestVersion) {
        this.loading = false
        this.isPageTransitioning = false
      }
    }
  }

  setSort(column: string) {
    this.query.setSort(column)
    void this.load()
  }

  handleSearch(value: string) {
    this.query.setSearch(value)

    if (this.searchDebounceTimer) {
      clearTimeout(this.searchDebounceTimer)
    }

    this.searchDebounceTimer = setTimeout(() => {
      this.searchDebounceTimer = null
      void this.load()
    }, 250)
  }

  previousPage() {
    this.query.previousPage()
    void this.load()
  }

  nextPage() {
    this.query.nextPage()
    void this.load()
  }

  handleCollegeFilter(value: string) {
    this.filterCollegeCode = value
    this.query.page = 1
    void this.load()
  }

  async remove() {
    if (!this.deleteTarget) return

    try {
      await DeleteProgram(this.deleteTarget.code)
      pushToast(`Program ${this.deleteTarget.code} deleted`, 'success')
      this.closeConfirm()
      await this.load()
    } catch (caught: unknown) {
      pushToast(caught instanceof Error ? caught.message : 'Failed to delete program', 'error')
    }
  }

  async saveProgram() {
    const codeError = validateCode(this.code.trim(), 'Program code')
    const nameError = validateLabel(this.name.trim(), 'Program name')

    if (codeError || nameError) {
      this.formError = codeError ?? nameError ?? 'Complete all required fields before saving'
      return
    }

    const selectedCollegeRow = this.colleges.find((college) => college.code === this.selectedCollegeCode)

    const payload = new ProgramModel({
      code: this.code.trim(),
      name: this.name.trim(),
      college_code: this.selectedCollegeCode,
      college: new CollegeModel({
        code: selectedCollegeRow?.code ?? this.selectedCollegeCode,
        name: selectedCollegeRow?.name ?? ''
      })
    })

    try {
      if (this.formMode === 'edit') {
        await UpdateProgram(payload)
        pushToast(`Program ${this.code.trim()} updated`, 'success')
      } else {
        await CreateProgram(payload)
        pushToast(`Program ${this.code.trim()} created`, 'success')
      }

      this.closeForm()
      await this.load()
    } catch (caught: unknown) {
      this.formError = caught instanceof Error ? caught.message : 'Failed to save program'
    }
  }

  handleCollegeSelect(option: CollegeOption | null) {
    this.selectedCollegeCode = option?.code ?? ''
  }
}
