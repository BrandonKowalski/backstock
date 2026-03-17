<script setup lang="ts">
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

defineProps<{
  title?: string
  message: string
  confirmLabel?: string
  destructive?: boolean
}>()
const emit = defineEmits<{ confirm: []; cancel: [] }>()
</script>

<template>
  <Dialog>
    <DialogContent @close="emit('cancel')">
      <DialogHeader>
        <DialogTitle>{{ title ?? 'Confirm' }}</DialogTitle>
      </DialogHeader>
      <p class="text-sm text-muted-foreground mt-2">{{ message }}</p>
      <div class="flex gap-2 mt-6">
        <Button variant="outline" class="flex-1" @click="emit('cancel')">Cancel</Button>
        <Button :variant="destructive ? 'destructive' : 'default'" class="flex-1" @click="emit('confirm')">
          {{ confirmLabel ?? 'Confirm' }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
