<script setup lang="ts">
import { DialogPortal, DialogOverlay, DialogContent } from 'reka-ui'
import { X } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

withDefaults(defineProps<{ class?: string; side?: 'bottom' | 'right' }>(), { side: 'bottom' })
const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <DialogPortal>
    <DialogOverlay class="fixed inset-0 z-50 bg-black/50" />
    <DialogContent
      :class="cn(
        'fixed z-50 bg-background shadow-lg',
        side === 'bottom' && 'inset-x-0 bottom-0 border-t rounded-t-2xl max-h-[92vh] overflow-y-auto',
        side === 'right' && 'inset-y-0 right-0 h-full w-full border-l sm:max-w-md',
        $props.class,
      )"
      @escape-key-down="emit('close')"
      @pointer-down-outside="emit('close')"
    >
      <slot />
      <button
        class="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background hover:opacity-100 focus:outline-none"
        @click="emit('close')"
        aria-label="Close"
      >
        <X class="h-4 w-4" />
      </button>
    </DialogContent>
  </DialogPortal>
</template>
