<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Item } from '@/types'
import { formatQuantity, formatDate, formatLocation, isExpiringSoon, isExpired, isPastBestBy, isStaleStock } from '@/utils/format'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Pencil, Minus, Plus } from 'lucide-vue-next'
import { api } from '@/api/client'

const props = defineProps<{
  item: Item
  pendingChanges: Map<number, number>
  locations: string[]
}>()
const emit = defineEmits<{ edit: []; 'update-qty': [stockId: number, qty: number]; restocked: [] }>()

const isLow = computed(() => props.item.low_quantity_threshold != null && props.item.total_quantity <= props.item.low_quantity_threshold)
const expiring = computed(() => isExpiringSoon(props.item.expiration_date))
const expired = computed(() => isExpired(props.item.expiration_date))
const pastBestBy = computed(() => !expired.value && isPastBestBy(props.item.best_by_date))
const hasStale = computed(() => (props.item.stock ?? []).some(s => s.quantity > 0 && isStaleStock(s.date_added)))
const isOutOfStock = computed(() => props.item.total_quantity === 0)

const warnings = computed(() => {
  const w: { label: string; badge: string; card: string }[] = []
  if (isOutOfStock.value) w.push({ label: 'Out of stock', badge: 'bg-status-oos/15 text-status-oos border-status-oos/30', card: 'border-status-oos/40 bg-status-oos/5' })
  if (expired.value) w.push({ label: 'Expired', badge: 'bg-status-expired/15 text-status-expired border-status-expired/30', card: 'border-status-expired/40 bg-status-expired/5' })
  else if (expiring.value) w.push({ label: 'Expiring soon', badge: 'bg-status-expiring/15 text-status-expiring border-status-expiring/30', card: 'border-status-expiring/40 bg-status-expiring/5' })
  if (pastBestBy.value) w.push({ label: 'Past best by', badge: 'bg-status-bestby/15 text-status-bestby border-status-bestby/30', card: 'border-status-bestby/40 bg-status-bestby/5' })
  if (hasStale.value) w.push({ label: 'Over 1 year old', badge: 'bg-status-stale/15 text-status-stale border-status-stale/30', card: 'border-status-stale/40 bg-status-stale/5' })
  if (isLow.value && !isOutOfStock.value) w.push({ label: 'Replenish', badge: 'bg-status-replenish/15 text-status-replenish border-status-replenish/30', card: 'border-status-replenish/40 bg-status-replenish/5' })
  return w
})

const cardClass = computed(() => warnings.value.length ? warnings.value[0].card : '')
const nonZeroStock = computed(() => (props.item.stock ?? []).filter(s => s.quantity > 0))

// Quick restock for out-of-stock items
const restockLocation = ref('')
const restockQty = ref(1)
const restocking = ref(false)

async function quickRestock() {
  if (!restockLocation.value || restockQty.value <= 0) return
  restocking.value = true
  try {
    await api.addStock(props.item.id, { location: restockLocation.value, quantity: restockQty.value })
    restockLocation.value = ''
    restockQty.value = 1
    emit('restocked')
  } finally {
    restocking.value = false
  }
}

function getDisplayQty(stockId: number, original: number): number {
  return props.pendingChanges.has(stockId) ? props.pendingChanges.get(stockId)! : original
}

function adjust(stockId: number, original: number, delta: number) {
  const current = getDisplayQty(stockId, original)
  const next = Math.max(0, current + delta)
  emit('update-qty', stockId, next)
}

function onManualInput(stockId: number, value: string) {
  const num = parseFloat(value)
  if (!isNaN(num) && num >= 0) {
    emit('update-qty', stockId, num)
  }
}
</script>

<template>
  <Card :class="'p-4 ' + cardClass">
    <div class="flex items-start justify-between">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2">
          <h3 class="font-medium leading-none truncate">{{ item.name }}</h3>
        </div>
        <div v-if="item.categories?.length || warnings.length" class="mt-1.5 flex flex-wrap gap-1">
          <Badge v-for="cat in item.categories" :key="cat.id" variant="secondary" class="text-xs px-1.5 py-0 font-normal">
            {{ cat.name }}
          </Badge>
          <span
            v-for="(w, i) in warnings"
            :key="'w' + i"
            class="inline-flex items-center text-xs font-medium px-1.5 py-0 rounded-md border"
            :class="w.badge"
          >{{ w.label }}</span>
        </div>
        <div v-if="item.expiration_date || item.best_by_date" class="mt-1 flex gap-3 text-xs text-muted-foreground">
          <span v-if="item.expiration_date">Exp: {{ formatDate(item.expiration_date) }}</span>
          <span v-if="item.best_by_date">BB: {{ formatDate(item.best_by_date) }}</span>
        </div>
      </div>
      <div class="flex items-center gap-2 ml-3 shrink-0">
        <span class="text-sm font-semibold tabular-nums">
          {{ formatQuantity(item.total_quantity, item.unit?.abbreviation) }}
        </span>
        <Button variant="outline" size="sm" class="h-8 w-8 p-0" @click="$emit('edit')" :aria-label="'Edit ' + item.name">
          <Pencil class="h-4 w-4" />
        </Button>
      </div>
    </div>

    <!-- Stock rows with +/- -->
    <div v-if="nonZeroStock.length" class="mt-3 grid gap-2 grid-cols-1">
      <div
        v-for="s in nonZeroStock"
        :key="s.id"
        class="flex items-center gap-1 rounded-lg border px-2 py-1.5"
        :class="pendingChanges.has(s.id) ? 'border-ring bg-primary/5' : 'border-input bg-muted/30'"
      >
        <span class="text-xs font-medium text-muted-foreground flex-1 truncate">{{ formatLocation(s.location) }}</span>
        <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="adjust(s.id, s.quantity, -1)" :aria-label="'Decrease quantity for ' + formatLocation(s.location)">
          <Minus class="h-3.5 w-3.5" />
        </Button>
        <input
          type="number"
          :value="getDisplayQty(s.id, s.quantity)"
          @change="onManualInput(s.id, ($event.target as HTMLInputElement).value)"
          min="0"
          step="any"
          data-1p-ignore
          :aria-label="'Quantity for ' + formatLocation(s.location)"
          class="w-14 h-8 px-1 rounded-md border border-input bg-background text-sm font-semibold text-foreground text-center tabular-nums focus:outline-none focus:ring-1 focus:ring-ring"
        />
        <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="adjust(s.id, s.quantity, 1)" :aria-label="'Increase quantity for ' + formatLocation(s.location)">
          <Plus class="h-3.5 w-3.5" />
        </Button>
        <span v-if="item.unit" class="text-xs text-muted-foreground w-6 text-right shrink-0">{{ item.unit.abbreviation }}</span>
      </div>
    </div>

    <!-- Quick restock for out-of-stock items -->
    <div v-else class="mt-3 flex items-center gap-2 rounded-lg border border-dashed border-input px-2 py-1.5">
      <select
        v-model="restockLocation"
        data-1p-ignore
        aria-label="Restock location"
        class="flex-1 h-8 rounded-md border border-input bg-background px-2 text-xs focus:outline-none focus:ring-1 focus:ring-ring"
      >
        <option value="">Add to...</option>
        <option v-for="loc in locations" :key="loc" :value="loc">{{ formatLocation(loc) }}</option>
      </select>
      <input
        v-model.number="restockQty"
        type="number"
        min="1"
        step="any"
        data-1p-ignore
        aria-label="Restock quantity"
        class="w-14 h-8 px-1 rounded-md border border-input bg-background text-sm font-semibold text-center tabular-nums focus:outline-none focus:ring-1 focus:ring-ring"
      />
      <Button size="sm" class="h-8" :disabled="!restockLocation || restocking" @click="quickRestock">
        <Plus class="h-3.5 w-3.5" />
        Add
      </Button>
    </div>
  </Card>
</template>
