<script setup lang="ts">
import type { Category } from '@/types'
import { formatLocation } from '@/utils/format'

defineProps<{ categories: Category[]; locations: string[] }>()

const location = defineModel<string>('location')
const category = defineModel<string>('category')
const sort = defineModel<string>('sort')

const emit = defineEmits<{ change: [] }>()

function update() {
  emit('change')
}
</script>

<template>
  <div class="mt-3 flex gap-2 overflow-x-auto pb-1">
    <select
      :value="location"
      @change="location = ($event.target as HTMLSelectElement).value; update()"
      aria-label="Filter by location"
      class="h-10 rounded-md border border-input bg-background px-3 text-sm ring-offset-background focus:outline-none focus:ring-1 focus:ring-ring"
    >
      <option value="">All Locations</option>
      <option v-for="loc in locations" :key="loc" :value="loc">{{ formatLocation(loc) }}</option>
    </select>

    <select
      :value="category"
      @change="category = ($event.target as HTMLSelectElement).value; update()"
      aria-label="Filter by category"
      class="h-10 rounded-md border border-input bg-background px-3 text-sm ring-offset-background focus:outline-none focus:ring-1 focus:ring-ring"
    >
      <option value="">All Categories</option>
      <option v-for="cat in categories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
    </select>

    <select
      :value="sort"
      @change="sort = ($event.target as HTMLSelectElement).value; update()"
      aria-label="Sort by"
      class="h-10 rounded-md border border-input bg-background px-3 text-sm ring-offset-background focus:outline-none focus:ring-1 focus:ring-ring"
    >
      <option value="name">Name</option>
      <option value="expiration">Expiration</option>
      <option value="recent">Recently Added</option>
    </select>
  </div>
</template>
