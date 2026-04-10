export const ROW_HEIGHT = 46
export const PAGE_CHROME = 220

export function calculatePageSize(viewportHeight: number) {
  return Math.max(5, Math.floor((viewportHeight - PAGE_CHROME) / ROW_HEIGHT))
}

export function setupResponsivePageSize(onPageSize: (next: number) => void, onReload: () => void) {
  onPageSize(calculatePageSize(window.innerHeight))

  const resize = () => {
    onPageSize(calculatePageSize(window.innerHeight))
    onReload()
  }

  window.addEventListener('resize', resize)
  return () => window.removeEventListener('resize', resize)
}
