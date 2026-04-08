import { writable } from 'svelte/store'

export type ToastKind = 'success' | 'error' | 'info'

export interface ToastMessage {
  id: number
  message: string
  kind: ToastKind
}

const DEFAULT_TIMEOUT = 2800

let nextId = 1

export const toasts = writable<ToastMessage[]>([])

export function pushToast(message: string, kind: ToastKind = 'success', timeout = DEFAULT_TIMEOUT) {
  const id = nextId++

  toasts.update((current) => [...current, { id, message, kind }])

  if (timeout > 0) {
    window.setTimeout(() => {
      toasts.update((current) => current.filter((toast) => toast.id !== id))
    }, timeout)
  }

  return id
}

export function removeToast(id: number) {
  toasts.update((current) => current.filter((toast) => toast.id !== id))
}