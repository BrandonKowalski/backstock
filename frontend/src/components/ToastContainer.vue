<script setup lang="ts">
import { useToast } from '@/composables/useToast'
import { X } from 'lucide-vue-next'

const { toasts } = useToast()

function dismiss(id: number) {
  toasts.value = toasts.value.filter(t => t.id !== id)
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-4 left-1/2 -translate-x-1/2 z-[100] flex flex-col gap-2 w-[90%] max-w-md pointer-events-none">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto flex items-center gap-2 rounded-lg border px-4 py-3 text-sm shadow-lg"
        :class="toast.type === 'error'
          ? 'bg-destructive/10 border-destructive/30 text-destructive'
          : 'bg-primary/10 border-primary/30 text-foreground'"
      >
        <span class="flex-1">{{ toast.message }}</span>
        <button @click="dismiss(toast.id)" class="shrink-0 opacity-70 hover:opacity-100" aria-label="Dismiss">
          <X class="h-3.5 w-3.5" />
        </button>
      </div>
    </div>
  </Teleport>
</template>
