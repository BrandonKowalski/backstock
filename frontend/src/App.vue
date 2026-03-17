<script setup lang="ts">
import { provide, onMounted, ref, computed, watch } from 'vue'
import { useItems, ItemsKey } from './composables/useItems'
import { useCategories } from './composables/useCategories'
import { useUnits } from './composables/useUnits'
import { useLocations } from './composables/useLocations'
import type { Item } from './types'
import { Button } from '@/components/ui/button'
import { Plus, PackageX, AlertTriangle } from 'lucide-vue-next'
import AppHeader from './components/AppHeader.vue'
import SearchBar from './components/SearchBar.vue'
import FilterBar from './components/FilterBar.vue'
import ItemList from './components/ItemList.vue'
import ItemDrawer from './components/ItemDrawer.vue'
import SettingsDrawer from './components/SettingsDrawer.vue'
import AddItemDialog from './components/AddItemDialog.vue'
import ToastContainer from './components/ToastContainer.vue'

const itemsCtx = useItems()
provide(ItemsKey, itemsCtx)

const { categories, fetchCategories } = useCategories()
const { units, fetchUnits } = useUnits()
const { locations, flatLocationNames, foodLocationNames, nonFoodLocationNames, excludedLocationNames, fetchLocations } = useLocations()

const foodMode = ref(true)

const foodCategories = computed(() => categories.value.filter(c => c.is_food))
const nonFoodCategories = computed(() => categories.value.filter(c => !c.is_food))
const modeCategories = computed(() => foodMode.value ? foodCategories.value : nonFoodCategories.value)
const modeLocations = computed(() => foodMode.value ? foodLocationNames.value : nonFoodLocationNames.value)
const modeUnits = computed(() => units.value.filter(u => u.is_food === null || u.is_food === foodMode.value))

const showItemDrawer = ref(false)
const editingItem = ref<Item | null>(null)
const showSettings = ref(false)
const showAddDialog = ref(false)
const showOutOfStock = ref(false)
const showReplenish = ref(false)
const hasPendingChanges = ref(false)

function itemMatchesMode(item: Item): boolean {
  return item.is_food === foodMode.value
}

function modeItems() {
  return itemsCtx.items.value.filter(i => itemMatchesMode(i))
}

function isItemExcluded(item: Item): boolean {
  // Only exclude if no location filter is active and ALL stock is in excluded locations
  if (itemsCtx.filters.value.location) return false
  if (!item.stock?.length) return false
  return item.stock.every(s => excludedLocationNames.value.has(s.location))
}

const displayItems = computed(() => {
  let items = modeItems()

  if (showReplenish.value) {
    return modeItems().filter(i => i.low_quantity_threshold != null && i.total_quantity <= i.low_quantity_threshold && i.total_quantity > 0)
  }

  if (showOutOfStock.value) {
    return modeItems().filter(i => i.total_quantity === 0)
  }

  return items.filter(i => i.total_quantity > 0 && !isItemExcluded(i))
})

const outOfStockCount = computed(() => modeItems().filter(i => i.total_quantity === 0).length)
const replenishCount = computed(() => modeItems().filter(i => i.low_quantity_threshold != null && i.total_quantity <= i.low_quantity_threshold && i.total_quantity > 0).length)

watch(outOfStockCount, (n) => { if (n === 0 && showOutOfStock.value) showOutOfStock.value = false })
watch(replenishCount, (n) => { if (n === 0 && showReplenish.value) showReplenish.value = false })

function switchMode(food: boolean) {
  foodMode.value = food
  itemsCtx.filters.value.location = ''
  itemsCtx.filters.value.category = ''
  itemsCtx.filters.value.search = ''
  showOutOfStock.value = false
  showReplenish.value = false
  itemsCtx.fetchItems()
}

function openNewItem(name?: string) {
  editingItem.value = name ? { name } as Item : null
  showItemDrawer.value = true
  showAddDialog.value = false
}

function openEditItem(item: Item) {
  editingItem.value = item
  showItemDrawer.value = true
  showAddDialog.value = false
}

function closeItemDrawer() {
  showItemDrawer.value = false
  editingItem.value = null
}

async function onItemSaved() {
  await itemsCtx.fetchItems()
  closeItemDrawer()
}

async function onItemDeleted() {
  await itemsCtx.fetchItems()
  closeItemDrawer()
}

onMounted(async () => {
  await Promise.all([itemsCtx.fetchItems(), fetchCategories(), fetchUnits(), fetchLocations()])
})
</script>

<template>
  <div class="min-h-screen bg-background">
    <AppHeader @open-settings="showSettings = true" />

    <main class="max-w-2xl mx-auto px-4 pb-24">
      <!-- Food / Non-food toggle -->
      <div class="mt-4 flex rounded-lg border border-input p-0.5 bg-muted" role="tablist" aria-label="Item mode">
        <button
          @click="switchMode(true)"
          role="tab"
          :aria-selected="foodMode"
          class="flex-1 text-sm font-medium py-1.5 rounded-md transition-colors"
          :class="foodMode ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground'"
        >
          Food
        </button>
        <button
          @click="switchMode(false)"
          role="tab"
          :aria-selected="!foodMode"
          class="flex-1 text-sm font-medium py-1.5 rounded-md transition-colors"
          :class="!foodMode ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground'"
        >
          Non-Food
        </button>
      </div>

      <SearchBar
        v-model="itemsCtx.filters.value.search"
        @search="itemsCtx.fetchItems()"
      />
      <FilterBar
        v-model:location="itemsCtx.filters.value.location"
        v-model:category="itemsCtx.filters.value.category"
        v-model:sort="itemsCtx.filters.value.sort"
        :categories="modeCategories"
        :locations="modeLocations"
        @change="itemsCtx.fetchItems()"
      />

      <div class="mt-3 flex gap-2">
        <Button
          v-if="outOfStockCount > 0"
          :variant="showOutOfStock ? 'default' : 'outline'"
          size="sm"
          @click="showOutOfStock = !showOutOfStock; showReplenish = false"
        >
          <PackageX class="h-3.5 w-3.5" />
          Out of Stock ({{ outOfStockCount }})
        </Button>
        <Button
          v-if="replenishCount > 0"
          :variant="showReplenish ? 'default' : 'outline'"
          size="sm"
          @click="showReplenish = !showReplenish; showOutOfStock = false"
        >
          <AlertTriangle class="h-3.5 w-3.5" />
          Replenish ({{ replenishCount }})
        </Button>
      </div>

      <ItemList
        :items="displayItems"
        :loading="itemsCtx.loading.value"
        :locations="flatLocationNames"
        @edit="openEditItem"
        @saved="itemsCtx.fetchItems()"
        @pending-status="hasPendingChanges = $event"
      />
    </main>

    <!-- FAB -->
    <Button
      v-show="!hasPendingChanges"
      @click="showAddDialog = true"
      class="fixed right-6 h-14 w-14 rounded-full shadow-lg z-40"
      :style="{ bottom: 'max(1.5rem, env(safe-area-inset-bottom, 1.5rem))' }"
      size="icon"
      :aria-label="foodMode ? 'Add food item' : 'Add non-food item'"
    >
      <Plus class="h-6 w-6" />
    </Button>

    <AddItemDialog
      v-if="showAddDialog"
      :food-mode="foodMode"
      :excluded-location-names="excludedLocationNames"
      @close="showAddDialog = false"
      @select-existing="openEditItem"
      @create-new="openNewItem"
    />

    <ItemDrawer
      v-if="showItemDrawer"
      :item="editingItem"
      :food-mode="foodMode"
      :categories="modeCategories"
      :all-categories="categories"
      :units="modeUnits"
      :locations="modeLocations"
      :all-locations="flatLocationNames"
      @close="closeItemDrawer"
      @saved="onItemSaved"
      @deleted="onItemDeleted"
    />

    <SettingsDrawer
      v-if="showSettings"
      :food-mode="foodMode"
      :categories="categories"
      :units="units"
      :locations="locations"
      @close="showSettings = false"
      @categories-changed="fetchCategories(); itemsCtx.fetchItems()"
      @units-changed="fetchUnits(); itemsCtx.fetchItems()"
      @locations-changed="fetchLocations()"
    />

    <ToastContainer />
  </div>
</template>
