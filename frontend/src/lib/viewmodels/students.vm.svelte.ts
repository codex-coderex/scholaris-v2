import { College as CollegeModel, Program as ProgramModel, Student as StudentModel } from '../../../bindings/scholaris-v2/internal/features/students/models.js'
import { CreateStudent, DeleteStudent, GetStudents, UpdateStudent } from '../../../bindings/scholaris-v2/appservice.js'

import { pushToast } from '$lib/stores/toast'
import { fetchCollegeOptions, fetchProgramOptions } from '$lib/services/lookups.service'
import { validateName, validateStudentID } from '$lib/shared/validation'
import type { CollegeRow, ProgramRow, StudentGender, StudentRow } from '$lib/types/schema'
import { TableQueryViewModel } from '$lib/viewmodels/table-query.vm.svelte'

type StudentFormMode = 'create' | 'edit'

type CollegeOption = {
  code: string
  name: string
}

type ProgramOption = {
  code: string
  name: string
  college_code?: string
}

export class StudentsViewModel {
  query = new TableQueryViewModel('s.id')

  students = $state<StudentRow[]>([])
  programs = $state<ProgramRow[]>([])
  colleges = $state<CollegeRow[]>([])
  loading = $state(false)
  isPageTransitioning = $state(false)
  error = $state('')

  requestVersion = 0

  formMode = $state<StudentFormMode>('create')
  formOpen = $state(false)
  confirmOpen = $state(false)
  deleteTarget = $state<StudentRow | null>(null)
  editStudentId = $state('')

  studentIdYear = $state('')
  studentIdSeq = $state('')
  firstName = $state('')
  lastName = $state('')
  year = $state('1')
  gender = $state<StudentGender | ''>('')
  selectedCollegeCode = $state('')
  selectedProgramCode = $state('')
  formError = $state('')

  selectedCollege = $derived(this.colleges.find((college) => college.code === this.selectedCollegeCode) ?? null)
  selectedProgram = $derived(this.programs.find((program) => program.code === this.selectedProgramCode) ?? null)

  collegeOptions = $derived(this.colleges)
  programOptions = $derived(
    this.programs.filter((program) => !this.selectedCollegeCode || program.college_code === this.selectedCollegeCode)
  )

  selectedCollegeLabel = $derived(
    this.selectedCollege
      ? `${this.selectedCollege.code} — ${this.selectedCollege.name}`
      : this.selectedProgram?.college
        ? `${this.selectedProgram.college.code} — ${this.selectedProgram.college.name}`
        : ''
  )
  selectedProgramLabel = $derived(this.selectedProgram ? `${this.selectedProgram.code} — ${this.selectedProgram.name}` : '')

  openAdd() {
    this.formMode = 'create'
    this.editStudentId = ''
    this.resetForm()
    this.formOpen = true
  }

  openEdit(student: StudentRow) {
    this.formMode = 'edit'
    this.editStudentId = student.id

    const parts = student.id.split('-')
    this.studentIdYear = parts[0] ?? ''
    this.studentIdSeq = parts[1] ?? ''
    this.firstName = student.first_name
    this.lastName = student.last_name
    this.year = String(student.year)
    this.gender = (student.gender as StudentGender) || ''
    this.selectedProgramCode = student.program_code ?? ''
    this.selectedCollegeCode =
      student.program?.college_code ??
      student.program?.college?.code ??
      this.programs.find((program) => program.code === student.program_code)?.college_code ??
      ''
    this.formError = ''
    this.formOpen = true
  }

  askDelete(student: StudentRow) {
    this.deleteTarget = student
    this.confirmOpen = true
  }

  resetForm() {
    this.studentIdYear = `${new Date().getFullYear()}`
    this.studentIdSeq = ''
    this.firstName = ''
    this.lastName = ''
    this.year = '1'
    this.gender = ''
    this.selectedCollegeCode = ''
    this.selectedProgramCode = ''
    this.formError = ''
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
    const [programRows, collegeRows] = await Promise.all([fetchProgramOptions(), fetchCollegeOptions()])
    this.programs = programRows
    this.colleges = collegeRows
  }

  async load() {
    const requestId = ++this.requestVersion
    const hasExistingRows = this.students.length > 0

    this.loading = !hasExistingRows
    this.isPageTransitioning = hasExistingRows
    this.error = ''

    try {
      const [rows, total] = await GetStudents(
        this.query.search,
        this.query.sortBy,
        this.query.order,
        this.query.page,
        this.query.pageSize
      )

      if (requestId !== this.requestVersion) {
        return
      }

      this.students = rows ?? []
      this.query.setTotal(total)
    } catch (caught: unknown) {
      if (requestId !== this.requestVersion) {
        return
      }

      this.error = caught instanceof Error ? caught.message : 'Failed to load students'
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
    void this.load()
  }

  previousPage() {
    this.query.previousPage()
    void this.load()
  }

  nextPage() {
    this.query.nextPage()
    void this.load()
  }

  goToPage(nextPageNumber: number) {
    this.query.goToPage(nextPageNumber)
    void this.load()
  }

  handlePageJump(event: KeyboardEvent) {
    const nextPageNumber = this.query.parsePageJump(event)
    if (!nextPageNumber) return

    this.query.goToPage(nextPageNumber)
    void this.load()
  }

  async remove() {
    if (!this.deleteTarget) return

    try {
      await DeleteStudent(this.deleteTarget.id)
      pushToast(`Student ${this.deleteTarget.id} deleted`, 'success')
      this.closeConfirm()
      await this.load()
    } catch (caught: unknown) {
      pushToast(caught instanceof Error ? caught.message : 'Failed to delete student', 'error')
    }
  }

  async saveStudent() {
    const id = `${this.studentIdYear.trim()}-${this.studentIdSeq.trim()}`
    const idError = validateStudentID(id)
    const firstError = validateName(this.firstName.trim(), 'First name')
    const lastError = validateName(this.lastName.trim(), 'Last name')

    if (idError || firstError || lastError || !this.selectedProgramCode || !this.year || !this.gender) {
      this.formError = idError ?? firstError ?? lastError ?? 'Complete all required fields before saving'
      return
    }

    const selectedProgramRow = this.programs.find((program) => program.code === this.selectedProgramCode)
    const selectedCollegeRow = this.colleges.find((college) => college.code === this.selectedCollegeCode)

    const payload = new StudentModel({
      id,
      first_name: this.firstName.trim(),
      last_name: this.lastName.trim(),
      year: Number(this.year),
      gender: this.gender,
      program_code: this.selectedProgramCode,
      program: new ProgramModel({
        code: selectedProgramRow?.code ?? this.selectedProgramCode,
        name: selectedProgramRow?.name ?? '',
        college_code: selectedProgramRow?.college_code ?? this.selectedCollegeCode,
        college: new CollegeModel({
          code: selectedCollegeRow?.code ?? selectedProgramRow?.college_code ?? this.selectedCollegeCode,
          name: selectedCollegeRow?.name ?? selectedProgramRow?.college?.name ?? ''
        })
      })
    })

    try {
      if (this.formMode === 'edit' && this.editStudentId) {
        await UpdateStudent(payload)
        pushToast(`Student ${id} updated`, 'success')
      } else {
        await CreateStudent(payload)
        pushToast(`Student ${id} created`, 'success')
      }

      this.closeForm()
      await this.load()
    } catch (caught: unknown) {
      this.formError = caught instanceof Error ? caught.message : 'Failed to save student'
    }
  }

  handleCollegeSelect(option: CollegeOption | null) {
    this.selectedCollegeCode = option?.code ?? ''
    this.selectedProgramCode = ''
  }

  handleProgramSelect(option: ProgramOption | null) {
    this.selectedProgramCode = option?.code ?? ''

    if (option?.college_code) {
      this.selectedCollegeCode = option.college_code
    }
  }
}
