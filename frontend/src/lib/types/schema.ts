export type AppMode = 'light' | 'dark'

export type StudentGender = 'Male' | 'Female' | 'Other'

export interface CollegeRow {
  code: string
  name: string
  created_at?: string
  updated_at?: string
}

export interface ProgramRow {
  code: string
  name: string
  college_code: string
  college?: CollegeRow
  created_at?: string
  updated_at?: string
}

export interface StudentRow {
  id: string
  first_name: string
  last_name: string
  year: number
  gender: StudentGender | string
  program_code: string
  program?: ProgramRow
  created_at?: string
  updated_at?: string
}

export interface ListQuery {
  search: string
  sort_by: string
  order: 'ASC' | 'DESC'
  page: number
  page_size: number
}
