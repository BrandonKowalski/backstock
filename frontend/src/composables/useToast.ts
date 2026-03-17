import { ref } from 'vue'

export interface Toast {
  id: number
  message: string
  type: 'error' | 'success'
}

const toasts = ref<Toast[]>([])
let nextId = 0

function addToast(message: string, type: 'error' | 'success' = 'error', duration = 4000) {
  const id = nextId++
  toasts.value.push({ id, message, type })
  setTimeout(() => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }, duration)
}

export function useToast() {
  return { toasts, addToast }
}
