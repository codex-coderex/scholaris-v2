export type SortOrder = 'ASC' | 'DESC'

export class TableQueryViewModel {
  search = $state('')
  sortBy = $state('')
  order = $state<SortOrder>('ASC')
  page = $state(1)
  pageSize = $state(10)
  total = $state(0)

  totalPages = $derived(Math.max(1, Math.ceil(this.total / this.pageSize)))

  constructor(defaultSortBy: string, defaultOrder: SortOrder = 'ASC') {
    this.sortBy = defaultSortBy
    this.order = defaultOrder
  }

  setSort(column: string) {
    if (this.sortBy === column) {
      this.order = this.order === 'ASC' ? 'DESC' : 'ASC'
    } else {
      this.sortBy = column
      this.order = 'ASC'
    }

    this.page = 1
  }

  setSearch(value: string) {
    this.search = value
    this.page = 1
  }

  setPageSize(next: number) {
    this.pageSize = next
  }

  setTotal(next: number) {
    this.total = next
  }

  previousPage() {
    this.page = Math.max(1, this.page - 1)
  }

  nextPage() {
    this.page = Math.min(this.totalPages, this.page + 1)
  }

  goToPage(nextPageNumber: number) {
    this.page = Math.min(Math.max(nextPageNumber, 1), this.totalPages)
  }

  parsePageJump(event: KeyboardEvent) {
    if (event.key !== 'Enter') return null

    const input = event.currentTarget as HTMLInputElement
    const nextPageNumber = Number(input.value)

    if (!nextPageNumber || nextPageNumber < 1 || nextPageNumber > this.totalPages) {
      input.value = ''
      return null
    }

    return nextPageNumber
  }
}
