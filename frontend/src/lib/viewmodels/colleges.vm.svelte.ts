import { College as CollegeModel } from '../../../bindings/scholaris-v2/internal/features/colleges/models.js'
import { CreateCollege, DeleteCollege, GetColleges, UpdateCollege } from '../../../bindings/scholaris-v2/appservice.js'

import { pushToast } from '$lib/stores/toast'
import type { CollegeRow } from '$lib/types/schema'
import { validateCode, validateLabel } from '$lib/shared/validation'
import { TableQueryViewModel } from '$lib/viewmodels/table-query.vm.svelte'

type CollegeFormMode = 'create' | 'edit'

export class CollegesViewModel {
  query = new TableQueryViewModel('code')

  colleges = $state<CollegeRow[]>([])
  loading = $state(false)
  isPageTransitioning = $state(false)
  error = $state('')

  requestVersion = 0

  formMode = $state<CollegeFormMode>('create')
  formOpen = $state(false)
  confirmOpen = $state(false)
  deleteTarget = $state<CollegeRow | null>(null)

  code = $state('')
  name = $state('')
  formError = $state('')

  searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

  openAdd() {
    this.formMode = 'create'
    this.code = ''
    this.name = ''
    this.formError = ''
    this.formOpen = true
  }

  openEdit(college: CollegeRow) {
    this.formMode = 'edit'
    this.code = college.code
    this.name = college.name
    this.formError = ''
    this.formOpen = true
  }

  askDelete(college: CollegeRow) {
    this.deleteTarget = college
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

  async load() {
    const requestId = ++this.requestVersion
    const hasExistingRows = this.colleges.length > 0

    this.loading = !hasExistingRows
    this.isPageTransitioning = hasExistingRows
    this.error = ''

    try {
      const [rows, total] = await GetColleges(
        this.query.search,
        this.query.sortBy,
        this.query.order,
        this.query.page,
        this.query.pageSize
      )

      if (requestId !== this.requestVersion) {
        return
      }

      this.colleges = rows ?? []
      this.query.setTotal(total)
    } catch (caught: unknown) {
      if (requestId !== this.requestVersion) {
        return
      }

      this.error = caught instanceof Error ? caught.message : 'Failed to load colleges'
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

  async remove() {
    if (!this.deleteTarget) return

    try {
      await DeleteCollege(this.deleteTarget.code)
      pushToast(`College ${this.deleteTarget.code} deleted`, 'success')
      this.closeConfirm()
      await this.load()
    } catch (caught: unknown) {
      pushToast(caught instanceof Error ? caught.message : 'Failed to delete college', 'error')
    }
  }

  async saveCollege() {
    const codeError = validateCode(this.code.trim(), 'College code')
    const nameError = validateLabel(this.name.trim(), 'College name')

    if (codeError || nameError) {
      this.formError = codeError ?? nameError ?? 'Please complete all fields'
      return
    }

    const payload = new CollegeModel({
      code: this.code.trim(),
      name: this.name.trim()
    })

    try {
      if (this.formMode === 'edit') {
        await UpdateCollege(payload)
        pushToast(`College ${this.code.trim()} updated`, 'success')
      } else {
        await CreateCollege(payload)
        pushToast(`College ${this.code.trim()} created`, 'success')
      }

      this.closeForm()
      await this.load()
    } catch (caught: unknown) {
      this.formError = caught instanceof Error ? caught.message : 'Failed to save college'
    }
  }
}
