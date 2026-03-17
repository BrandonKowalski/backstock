<script setup lang="ts">
import { ref, inject, onMounted } from 'vue'
import type { Item, Category, Unit, ItemForm } from '@/types'
import { ItemsKey } from '@/composables/useItems'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Trash2, Pencil, ChevronDown, ChevronUp } from 'lucide-vue-next'
import StockSection from './StockSection.vue'
import ConfirmDialog from './ConfirmDialog.vue'

const props = defineProps<{
  item: Item | null
  foodMode: boolean
  categories: Category[]
  allCategories?: Category[]
  units: Unit[]
  locations: string[]
  allLocations?: string[]
}>()
const emit = defineEmits<{ close: []; saved: []; deleted: [] }>()

const ctx = inject(ItemsKey)!
const isEdit = !!props.item?.id
const showDetails = ref(!isEdit)
const showDeleteConfirm = ref(false)

const form = ref<ItemForm>({
  name: props.item?.name ?? '',
  is_food: props.item?.is_food ?? props.foodMode,
  unit_id: props.item?.unit_id ?? null,
  package_size: props.item?.package_size ?? null,
  expiration_date: props.item?.expiration_date ?? '',
  best_by_date: props.item?.best_by_date ?? '',
  low_quantity_threshold: props.item?.low_quantity_threshold ?? null,
  category_ids: props.item?.categories?.map(c => c.id) ?? [],
})

const fullItem = ref<Item | null>(null)
const saving = ref(false)

onMounted(async () => {
  if (props.item?.id) {
    fullItem.value = await ctx.getItem(props.item.id)
  }
})

async function save() {
  if (!form.value.name.trim()) return
  saving.value = true
  const payload = {
    ...form.value,
    expiration_date: form.value.expiration_date ?? null,
    best_by_date: form.value.best_by_date ?? null,
    package_size: form.value.package_size ?? null,
    low_quantity_threshold: form.value.low_quantity_threshold ?? null,
  }
  try {
    if (isEdit && props.item) {
      await ctx.updateItem(props.item.id, payload)
    } else {
      await ctx.createItem(payload)
    }
    emit('saved')
  } finally {
    saving.value = false
  }
}

async function remove() {
  if (!props.item) return
  await ctx.deleteItem(props.item.id)
  emit('deleted')
}

function toggleCategory(id: number) {
  const idx = form.value.category_ids.indexOf(id)
  if (idx >= 0) form.value.category_ids.splice(idx, 1)
  else form.value.category_ids.push(id)
}

async function refreshItem() {
  if (props.item?.id) {
    fullItem.value = await ctx.getItem(props.item.id)
  }
}
</script>

<template>
  <Sheet>
    <SheetContent side="bottom" @close="emit('close')">
      <SheetHeader class="flex-row items-center justify-between pr-10">
        <SheetTitle>{{ isEdit ? item?.name : 'New Item' }}</SheetTitle>
        <Button v-if="isEdit" variant="ghost" size="sm" @click="showDetails = !showDetails">
          <Pencil class="h-3.5 w-3.5" />
          {{ showDetails ? 'Hide' : 'Edit' }} Details
          <component :is="showDetails ? ChevronUp : ChevronDown" class="h-3.5 w-3.5" />
        </Button>
      </SheetHeader>

      <!-- Stock section (always visible when editing) -->
      <template v-if="isEdit && fullItem">
        <StockSection
          :item="fullItem"
          :locations="allLocations ?? locations"
          @changed="refreshItem(); $emit('saved')"
        />
        <div class="px-6 pb-4">
          <Button variant="destructive" class="w-full" @click="showDeleteConfirm = true">
            <Trash2 class="h-4 w-4" />
            Delete Item
          </Button>
        </div>
      </template>

      <!-- Item details form (collapsed by default when editing) -->
      <template v-if="showDetails">
        <Separator v-if="isEdit" />
        <form @submit.prevent="save" class="p-6 space-y-4">
          <div class="space-y-2">
            <Label>Name</Label>
            <Input v-model="form.name" type="text" required placeholder="Item name" />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-2">
              <Label>Unit</Label>
              <select v-model="form.unit_id" aria-label="Unit" class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring">
                <option :value="null">None</option>
                <option v-for="u in units" :key="u.id" :value="u.id">{{ u.name }} ({{ u.abbreviation }})</option>
              </select>
            </div>
            <div class="space-y-2">
              <Label>Package Size</Label>
              <Input v-model="form.package_size" type="number" step="any" min="0" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-2">
              <Label>Expiration</Label>
              <Input v-model="form.expiration_date" type="date" />
            </div>
            <div class="space-y-2">
              <Label>Best By</Label>
              <Input v-model="form.best_by_date" type="date" />
            </div>
          </div>

          <div class="space-y-2">
            <Label>Low Qty Threshold</Label>
            <Input v-model="form.low_quantity_threshold" type="number" step="any" min="0" />
          </div>

          <div v-if="categories.length" class="space-y-2">
            <Label>Categories</Label>
            <div class="flex flex-wrap gap-1.5" role="group" aria-label="Categories">
              <button
                v-for="cat in categories"
                :key="cat.id"
                type="button"
                :class="[
                  'inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors cursor-pointer select-none focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
                  form.category_ids.includes(cat.id)
                    ? 'border-transparent bg-primary text-primary-foreground shadow'
                    : 'text-foreground'
                ]"
                :aria-pressed="form.category_ids.includes(cat.id)"
                @click="toggleCategory(cat.id)"
              >
                {{ cat.name }}
              </button>
            </div>
          </div>

          <div class="pt-2">
            <Button type="submit" :disabled="saving" class="w-full">
              {{ saving ? 'Saving...' : (isEdit ? 'Update' : 'Create') }}
            </Button>
          </div>
        </form>
      </template>
    </SheetContent>
  </Sheet>

  <ConfirmDialog
    v-if="showDeleteConfirm"
    title="Delete Item"
    :message="'Delete ' + (item?.name ?? 'this item') + ' and all its stock?'"
    confirm-label="Delete"
    :destructive="true"
    @confirm="showDeleteConfirm = false; remove()"
    @cancel="showDeleteConfirm = false"
  />
</template>
