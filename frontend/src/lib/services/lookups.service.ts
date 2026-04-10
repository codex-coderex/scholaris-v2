import { GetColleges, GetPrograms } from '../../../bindings/scholaris-v2/appservice.js'

import type { CollegeRow, ProgramRow } from '$lib/types/schema'

export async function fetchCollegeOptions(limit = 1000): Promise<CollegeRow[]> {
  const [rows] = await GetColleges('', 'name', 'ASC', 1, limit)
  return rows ?? []
}

export async function fetchProgramOptions(limit = 1000): Promise<ProgramRow[]> {
  const [rows] = await GetPrograms('', 'p.code', 'ASC', 1, limit, '')
  return rows ?? []
}
