<script setup lang="ts">
import { DialogPortal, DialogOverlay, DialogContent } from 'reka-ui'
import { X } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

withDefaults(defineProps<{ class?: string }>(), {})
const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <DialogPortal>
    <DialogOverlay class="fixed inset-0 z-50 bg-black/50" />
    <DialogContent
      :class="cn(
        'fixed left-1/2 top-1/2 z-50 w-[95%] max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border bg-background p-6 shadow-lg',
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
