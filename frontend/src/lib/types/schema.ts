export type AppMode = 'light' | 'dark'

export type StudentGender = 'Male' | 'Female' | 'Other'

export interface CollegeRow {
  code: string
  name: string
}

export interface ProgramRow {
  code: string
  name: string
  college_code: string
  college?: CollegeRow
}

export interface StudentRow {
  id: string
  first_name: string
  last_name: string
  year: number
  gender: StudentGender | string
  program_code: string
  program?: ProgramRow
}
