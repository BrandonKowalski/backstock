<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Stock, Item } from '@/types'
import { formatLocation, formatQuantity } from '@/utils/format'
import { useStock } from '@/composables/useStock'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const props = defineProps<{ stock: Stock; item: Item; locations: string[] }>()
const emit = defineEmits<{ close: []; moved: [] }>()

const { moveStock } = useStock()

const toLocation = ref('')
const quantity = ref(props.stock.quantity)

const otherLocations = computed(() => props.locations.filter(l => l !== props.stock.location))

async function handleMove() {
  if (!toLocation.value || quantity.value <= 0) return
  await moveStock(props.stock.id, { to_location: toLocation.value, quantity: quantity.value })
  emit('moved')
}
</script>

<template>
  <Dialog>
    <DialogContent @close="emit('close')">
      <DialogHeader>
        <DialogTitle>Move Stock</DialogTitle>
      </DialogHeader>
      <p class="text-sm text-muted-foreground">
        From {{ formatLocation(stock.location) }} ({{ formatQuantity(stock.quantity, item.unit?.abbreviation) }} available)
      </p>

      <div class="space-y-4 mt-4">
        <div class="space-y-2">
          <Label>To</Label>
          <select v-model="toLocation" aria-label="Destination location" class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring">
            <option value="">Select location...</option>
            <option v-for="loc in otherLocations" :key="loc" :value="loc">{{ formatLocation(loc) }}</option>
          </select>
        </div>
        <div class="space-y-2">
          <Label>Quantity</Label>
          <Input v-model="quantity" type="number" min="0.1" :max="stock.quantity" step="any" />
        </div>
      </div>

      <div class="flex gap-2 mt-6">
        <Button variant="outline" class="flex-1" @click="emit('close')">Cancel</Button>
        <Button class="flex-1" @click="handleMove">Move</Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
