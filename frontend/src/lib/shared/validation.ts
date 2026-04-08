export function validateStudentID(id: string) {
  if (!id) return 'Student ID is required'
  if (!/^\d{4}-\d{4}$/.test(id)) return 'Student ID must be in YYYY-NNNN format (e.g. 2024-0001)'
  return null
}

export function validateName(name: string, field: string) {
  if (!name) return `${field} is required`
  if (name.length > 64) return `${field} must be 64 characters or fewer`
  if (!/^[a-zA-ZÀ-ÖØ-öø-ÿñÑ\s\-'.]+$/.test(name)) return `${field} contains invalid characters`
  return null
}

export function validateCode(code: string, field: string) {
  if (!code) return `${field} is required`
  if (code.length > 16) return `${field} must be 16 characters or fewer`
  if (!/^[A-Za-z0-9\-]+$/.test(code)) return `${field} must only contain letters, numbers, or hyphens only`
  return null
}

export function validateLabel(name: string, field: string, maxLen = 128) {
  if (!name) return `${field} is required`
  if (name.length > maxLen) return `${field} must be ${maxLen} characters or fewer`
  if (!/^[a-zA-Z0-9\s\-'.]+$/.test(name)) return `${field} contains invalid characters`
  return null
}
