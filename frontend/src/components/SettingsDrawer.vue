<script setup lang="ts">
import { ref } from 'vue'
import type { Category, Unit, Location } from '@/types'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { Separator } from '@/components/ui/separator'
import { ChevronDown, ChevronRight } from 'lucide-vue-next'
import CategoryManager from './CategoryManager.vue'
import UnitManager from './UnitManager.vue'
import LocationManager from './LocationManager.vue'

defineProps<{ foodMode: boolean; categories: Category[]; units: Unit[]; locations: Location[] }>()
defineEmits<{ close: []; 'categories-changed': []; 'units-changed': []; 'locations-changed': [] }>()

const openSections = ref<Set<string>>(new Set(['locations', 'categories', 'units']))

function toggle(section: string) {
  if (openSections.value.has(section)) {
    openSections.value.delete(section)
  } else {
    openSections.value.add(section)
  }
  openSections.value = new Set(openSections.value)
}
</script>

<template>
  <Sheet>
    <SheetContent side="right" class="sm:max-w-xl overflow-y-auto" @close="$emit('close')">
      <SheetHeader>
        <SheetTitle>Backstock Settings</SheetTitle>
      </SheetHeader>

      <div class="p-6 space-y-1">
        <!-- Locations -->
        <div>
          <button @click="toggle('locations')" class="w-full flex items-center justify-between py-3 text-sm font-semibold">
            Locations
            <component :is="openSections.has('locations') ? ChevronDown : ChevronRight" class="h-4 w-4 text-muted-foreground" />
          </button>
          <div v-if="openSections.has('locations')" class="pb-4">
            <LocationManager :locations="locations" @changed="$emit('locations-changed')" />
          </div>
        </div>

        <Separator />

        <!-- Categories -->
        <div>
          <button @click="toggle('categories')" class="w-full flex items-center justify-between py-3 text-sm font-semibold">
            Categories
            <component :is="openSections.has('categories') ? ChevronDown : ChevronRight" class="h-4 w-4 text-muted-foreground" />
          </button>
          <div v-if="openSections.has('categories')" class="pb-4">
            <CategoryManager :categories="categories" :food-mode="foodMode" @changed="$emit('categories-changed')" />
          </div>
        </div>

        <Separator />

        <!-- Units -->
        <div>
          <button @click="toggle('units')" class="w-full flex items-center justify-between py-3 text-sm font-semibold">
            Units
            <component :is="openSections.has('units') ? ChevronDown : ChevronRight" class="h-4 w-4 text-muted-foreground" />
          </button>
          <div v-if="openSections.has('units')" class="pb-4">
            <UnitManager :units="units" @changed="$emit('units-changed')" />
          </div>
        </div>
      </div>
    </SheetContent>
  </Sheet>
</template>
