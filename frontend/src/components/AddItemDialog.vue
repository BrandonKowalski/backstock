<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Item } from '@/types'
import { formatQuantity } from '@/utils/format'
import { api } from '@/api/client'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Plus } from 'lucide-vue-next'

const props = defineProps<{
  foodMode: boolean
  excludedLocationNames: Set<string>
}>()
const emit = defineEmits<{ close: []; 'select-existing': [item: Item]; 'create-new': [name: string] }>()

const allItems = ref<Item[]>([])
const query = ref('')

onMounted(async () => {
  allItems.value = await api.listItems()
})

const modeItems = computed(() => allItems.value.filter(item => {
  if (item.is_food !== props.foodMode) return false
  if (item.stock?.length && item.stock.every(s => props.excludedLocationNames.has(s.location))) return false
  if (item.total_quantity > 0) return false
  return true
}))

const filtered = computed(() => {
  const q = query.value.toLowerCase().trim()
  if (!q) return modeItems.value
  return modeItems.value.filter(item => item.name.toLowerCase().includes(q))
})

const exactMatch = computed(() => {
  const q = query.value.toLowerCase().trim()
  return q && modeItems.value.some(item => item.name.toLowerCase() === q)
})
</script>

<template>
  <Dialog>
    <DialogContent class="max-h-[80vh] flex flex-col p-0 gap-0" @close="emit('close')">
      <DialogHeader class="p-6 pb-0">
        <DialogTitle>Add {{ foodMode ? 'Food' : 'Non-Food' }} Item</DialogTitle>
      </DialogHeader>
      <div class="px-6 pt-3 pb-2">
        <Input
          v-model="query"
          type="text"
          placeholder="Search or type new item name..."
          autofocus
        />
      </div>

      <div class="flex-1 overflow-y-auto px-2 pb-4">
        <button
          v-if="query.trim() && !exactMatch"
          @click="emit('create-new', query.trim())"
          class="w-full text-left px-4 py-2.5 rounded-md hover:bg-accent flex items-center gap-2 text-sm text-primary font-medium"
        >
          <Plus class="h-4 w-4" />
          Create "{{ query.trim() }}"
        </button>

        <button
          v-for="item in filtered"
          :key="item.id"
          @click="emit('select-existing', item)"
          class="w-full text-left px-4 py-2.5 rounded-md hover:bg-accent flex items-center justify-between text-sm"
        >
          <span>{{ item.name }}</span>
          <span class="text-muted-foreground text-xs tabular-nums">
            {{ formatQuantity(item.total_quantity, item.unit?.abbreviation) }}
          </span>
        </button>

        <p v-if="!filtered.length && !query.trim()" class="text-sm text-muted-foreground text-center py-6">Type a name to create a new item.</p>
      </div>
    </DialogContent>
  </Dialog>
</template>
