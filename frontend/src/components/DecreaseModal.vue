<script setup lang="ts">
import { ref } from 'vue'
import type { Stock, Item } from '@/types'
import { formatLocation, formatQuantity } from '@/utils/format'
import { useStock } from '@/composables/useStock'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const props = defineProps<{ stock: Stock; item: Item }>()
const emit = defineEmits<{ close: []; decreased: [] }>()

const { updateStock } = useStock()

const amount = ref(1)

async function handleDecrease() {
  if (amount.value <= 0) return
  const newQty = props.stock.quantity - amount.value
  await updateStock(props.stock.id, Math.max(0, newQty))
  emit('decreased')
}
</script>

<template>
  <Dialog>
    <DialogContent @close="emit('close')">
      <DialogHeader>
        <DialogTitle>Decrease Stock</DialogTitle>
      </DialogHeader>
      <p class="text-sm text-muted-foreground">
        {{ formatLocation(stock.location) }}: {{ formatQuantity(stock.quantity, item.unit?.abbreviation) }} available
      </p>

      <div class="space-y-2 mt-4">
        <Label>Amount to remove</Label>
        <Input v-model="amount" type="number" min="0.1" :max="stock.quantity" step="any" />
      </div>

      <div class="flex gap-2 mt-6">
        <Button variant="outline" class="flex-1" @click="emit('close')">Cancel</Button>
        <Button variant="destructive" class="flex-1" @click="handleDecrease">Decrease</Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
