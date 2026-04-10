export const ROW_HEIGHT = 46
export const PAGE_CHROME = 220

export function calculatePageSize(viewportHeight: number) {
  return Math.max(5, Math.floor((viewportHeight - PAGE_CHROME) / ROW_HEIGHT))
}

export function setupResponsivePageSize(onPageSize: (next: number) => void, onReload: () => void) {
  let currentPageSize = calculatePageSize(window.innerHeight)
  onPageSize(currentPageSize)
  let resizeDebounceTimer: ReturnType<typeof setTimeout> | null = null

  const resize = () => {
    const nextPageSize = calculatePageSize(window.innerHeight)
    if (nextPageSize === currentPageSize) {
      return
    }

    currentPageSize = nextPageSize

    if (resizeDebounceTimer) {
      clearTimeout(resizeDebounceTimer)
    }

    resizeDebounceTimer = setTimeout(() => {
      resizeDebounceTimer = null
      onPageSize(currentPageSize)
      onReload()
    }, 180)
  }

  window.addEventListener('resize', resize)
  return () => {
    if (resizeDebounceTimer) {
      clearTimeout(resizeDebounceTimer)
    }
    window.removeEventListener('resize', resize)
  }
}
