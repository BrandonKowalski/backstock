<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Item } from '@/types'
import { api } from '@/api/client'
import { Button } from '@/components/ui/button'
import { Check, X } from 'lucide-vue-next'
import ItemCard from './ItemCard.vue'

const props = defineProps<{ items: Item[]; loading: boolean; locations: string[] }>()
const emit = defineEmits<{ edit: [item: Item]; saved: []; 'pending-status': [has: boolean] }>()

const pendingChanges = ref<Map<number, number>>(new Map())
const saving = ref(false)

watch(() => pendingChanges.value.size, (size) => {
  emit('pending-status', size > 0)
})

function onUpdateQty(stockId: number, qty: number) {
  for (const item of props.items) {
    const stock = item.stock?.find(s => s.id === stockId)
    if (stock) {
      if (qty === stock.quantity) {
        pendingChanges.value.delete(stockId)
      } else {
        pendingChanges.value.set(stockId, qty)
      }
      pendingChanges.value = new Map(pendingChanges.value)
      return
    }
  }
}

async function saveAll() {
  saving.value = true
  try {
    const promises = Array.from(pendingChanges.value.entries()).map(([stockId, qty]) =>
      api.updateStock(stockId, qty)
    )
    await Promise.all(promises)
    pendingChanges.value = new Map()
    emit('saved')
  } finally {
    saving.value = false
  }
}

function discard() {
  pendingChanges.value = new Map()
}
</script>

<template>
  <div class="mt-4 space-y-2">
    <div v-if="loading" class="space-y-2">
      <div v-for="i in 4" :key="i" class="rounded-xl border bg-card p-4 animate-pulse">
        <div class="flex items-start justify-between">
          <div class="flex-1 space-y-2">
            <div class="h-4 w-32 rounded bg-muted"></div>
            <div class="flex gap-1">
              <div class="h-4 w-16 rounded bg-muted"></div>
              <div class="h-4 w-20 rounded bg-muted"></div>
            </div>
          </div>
          <div class="h-8 w-16 rounded bg-muted"></div>
        </div>
        <div class="mt-3 h-10 rounded-lg bg-muted"></div>
      </div>
    </div>
    <div v-else-if="items.length === 0" class="text-center py-12 text-muted-foreground">
      No items found.
    </div>
    <ItemCard
      v-for="item in items"
      :key="item.id"
      :item="item"
      :pending-changes="pendingChanges"
      :locations="locations"
      @edit="$emit('edit', item)"
      @update-qty="onUpdateQty"
      @restocked="$emit('saved')"
    />

    <!-- Floating save bar -->
    <Teleport to="body">
        <div
          v-if="pendingChanges.size > 0"
          class="fixed bottom-0 left-0 right-0 z-50 border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80 shadow-lg px-4 py-3 flex items-center justify-between gap-3"
          :style="{ paddingBottom: 'max(0.75rem, env(safe-area-inset-bottom))' }"
        >
          <span class="text-sm text-muted-foreground">
            {{ pendingChanges.size }} change{{ pendingChanges.size === 1 ? '' : 's' }}
          </span>
          <div class="flex gap-2">
            <Button variant="outline" size="sm" @click="discard">
              <X class="h-3.5 w-3.5" />
              Discard
            </Button>
            <Button size="sm" :disabled="saving" @click="saveAll">
              <Check class="h-3.5 w-3.5" />
              {{ saving ? 'Saving...' : 'Save' }}
            </Button>
          </div>
        </div>
    </Teleport>
  </div>
</template>
