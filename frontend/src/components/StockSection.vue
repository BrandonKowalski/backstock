<script setup lang="ts">
import { ref } from 'vue'
import type { Item, Stock } from '@/types'
import { formatQuantity, formatLocation } from '@/utils/format'
import { useStock } from '@/composables/useStock'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Minus, ArrowRightLeft, Trash2, Plus } from 'lucide-vue-next'
import StockMoveModal from './StockMoveModal.vue'
import DecreaseModal from './DecreaseModal.vue'
import ConfirmDialog from './ConfirmDialog.vue'

const props = defineProps<{ item: Item; locations: string[] }>()
const emit = defineEmits<{ changed: [] }>()

const { addStock, deleteStock } = useStock()

const addingLocation = ref('')
const addingQty = ref<number>(1)
const showAdd = ref(false)

const moveTarget = ref<Stock | null>(null)
const decreaseTarget = ref<Stock | null>(null)
const deleteTarget = ref<Stock | null>(null)

async function handleAdd() {
  if (!addingLocation.value || addingQty.value <= 0) return
  await addStock(props.item.id, { location: addingLocation.value, quantity: addingQty.value })
  showAdd.value = false
  addingLocation.value = ''
  addingQty.value = 1
  emit('changed')
}

async function handleDelete() {
  if (!deleteTarget.value) return
  await deleteStock(deleteTarget.value.id)
  deleteTarget.value = null
  emit('changed')
}
</script>

<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-3">
      <h3 class="font-semibold text-sm">Stock</h3>
      <Button variant="ghost" size="sm" @click="showAdd = !showAdd">
        <Plus v-if="!showAdd" class="h-3.5 w-3.5" />
        {{ showAdd ? 'Cancel' : 'Add Stock' }}
      </Button>
    </div>

    <div v-if="showAdd" class="mb-4 flex gap-2 items-center">
      <div class="flex-1">
        <select v-model="addingLocation" aria-label="Location" class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring">
          <option value="">Location...</option>
          <option v-for="loc in locations" :key="loc" :value="loc">{{ formatLocation(loc) }}</option>
        </select>
      </div>
      <Input v-model="addingQty" type="number" min="0.1" step="any" class="w-20" />
      <Button @click="handleAdd">Add</Button>
    </div>

    <p v-if="!item.stock?.length" class="text-sm text-muted-foreground py-2">No stock entries yet.</p>

    <table v-else class="w-full text-sm">
      <thead>
        <tr class="border-b text-muted-foreground">
          <th class="text-left font-medium py-2">Location</th>
          <th class="text-right font-medium py-2">Qty</th>
          <th class="text-right font-medium py-2 w-28">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in item.stock" :key="s.id" class="border-b last:border-0">
          <td class="py-2.5 font-medium">{{ formatLocation(s.location) }}</td>
          <td class="py-2.5 text-right tabular-nums">
            {{ formatQuantity(s.quantity, item.unit?.abbreviation) }}
          </td>
          <td class="py-2.5 text-right">
            <div class="inline-flex gap-1">
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="decreaseTarget = s" aria-label="Decrease stock">
                <Minus class="h-3.5 w-3.5" />
              </Button>
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="moveTarget = s" aria-label="Move stock">
                <ArrowRightLeft class="h-3.5 w-3.5" />
              </Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:text-destructive" @click="deleteTarget = s" aria-label="Remove stock">
                <Trash2 class="h-3.5 w-3.5" />
              </Button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <StockMoveModal
      v-if="moveTarget"
      :stock="moveTarget"
      :item="item"
      :locations="locations"
      @close="moveTarget = null"
      @moved="moveTarget = null; emit('changed')"
    />
    <DecreaseModal
      v-if="decreaseTarget"
      :stock="decreaseTarget"
      :item="item"
      @close="decreaseTarget = null"
      @decreased="decreaseTarget = null; emit('changed')"
    />
    <ConfirmDialog
      v-if="deleteTarget"
      title="Remove Stock"
      :message="'Remove all ' + formatQuantity(deleteTarget.quantity, item.unit?.abbreviation) + ' from ' + formatLocation(deleteTarget.location) + '?'"
      confirm-label="Remove"
      :destructive="true"
      @confirm="handleDelete()"
      @cancel="deleteTarget = null"
    />
  </div>
</template>
