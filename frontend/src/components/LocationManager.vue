<script setup lang="ts">
import { ref } from 'vue'
import type { Location } from '@/types'
import { useLocations } from '@/composables/useLocations'
import { formatLocation } from '@/utils/format'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Trash2, Plus, ChevronRight } from 'lucide-vue-next'

defineProps<{ locations: Location[] }>()
const emit = defineEmits<{ changed: [] }>()

const { createLocation, deleteLocation, updateLocation } = useLocations()
const newName = ref('')
const addingChildOf = ref<number | null>(null)
const childName = ref('')
const editingId = ref<number | null>(null)
const editingName = ref('')

async function add() {
  const slug = newName.value.trim().toLowerCase().replace(/\s+/g, '_')
  if (!slug) return
  await createLocation(slug)
  newName.value = ''
  emit('changed')
}

async function addChild(parentId: number) {
  const slug = childName.value.trim().toLowerCase().replace(/\s+/g, '_')
  if (!slug) return
  await createLocation(slug, parentId)
  childName.value = ''
  addingChildOf.value = null
  emit('changed')
}

function startEdit(loc: Location) {
  editingId.value = loc.id
  editingName.value = formatLocation(loc.name)
}

async function saveEdit(loc: Location) {
  const slug = editingName.value.trim().toLowerCase().replace(/\s+/g, '_')
  editingId.value = null
  if (!slug || slug === loc.name) return
  await updateLocation(loc.id, { name: slug, parent_id: loc.parent_id, is_food: loc.is_food })
  emit('changed')
}

async function remove(id: number) {
  await deleteLocation(id)
  emit('changed')
}

async function toggleFood(loc: Location) {
  await updateLocation(loc.id, { name: loc.name, parent_id: loc.parent_id, is_food: !loc.is_food, exclude_default: loc.exclude_default })
  emit('changed')
}

async function toggleExclude(loc: Location) {
  await updateLocation(loc.id, { name: loc.name, parent_id: loc.parent_id, is_food: loc.is_food, exclude_default: !loc.exclude_default })
  emit('changed')
}
</script>

<template>
  <div>
    <div class="flex gap-2 items-center">
      <Input v-model="newName" @keyup.enter="add" placeholder="New location..." />
      <Button @click="add">Add</Button>
    </div>

    <div class="mt-3 space-y-1">
      <div v-for="loc in locations" :key="loc.id">
        <div class="py-2 px-3 rounded-md bg-muted">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 min-w-0 flex-1">
              <input
                v-if="editingId === loc.id"
                v-model="editingName"
                @blur="saveEdit(loc)"
                @keyup.enter="saveEdit(loc)"
                @keyup.escape="editingId = null"
                class="text-sm font-medium bg-background border border-input rounded px-1.5 py-0.5 w-full max-w-[10rem] focus:outline-none focus:ring-1 focus:ring-ring"
                autofocus
              />
              <span
                v-else
                @click="startEdit(loc)"
                class="text-sm font-medium cursor-text hover:underline decoration-dashed underline-offset-4 truncate"
              >{{ formatLocation(loc.name) }}</span>
            </div>
            <div class="flex items-center gap-0.5 shrink-0">
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="addingChildOf = addingChildOf === loc.id ? null : loc.id" aria-label="Add sub-location">
                <Plus class="h-3.5 w-3.5" />
              </Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:text-destructive" @click="remove(loc.id)" :aria-label="'Delete ' + formatLocation(loc.name)">
                <Trash2 class="h-3.5 w-3.5" />
              </Button>
            </div>
          </div>
          <div class="flex flex-wrap gap-1.5 mt-1.5">
            <button
              @click="toggleFood(loc)"
              class="text-xs font-medium px-2 py-0.5 rounded-full border transition-colors cursor-pointer"
              :class="loc.is_food
                ? 'bg-tag-food-bg border-tag-food/30 text-tag-food hover:bg-tag-food/20'
                : 'bg-tag-both-bg border-tag-both/30 text-tag-both hover:bg-tag-both/15'"
            >
              {{ loc.is_food ? 'Food' : 'Non-food' }}
            </button>
            <button
              @click="toggleExclude(loc)"
              class="text-xs font-medium px-2 py-0.5 rounded-full border transition-colors cursor-pointer"
              :class="loc.exclude_default
                ? 'bg-tag-nonfood-bg border-tag-nonfood/30 text-tag-nonfood hover:bg-tag-nonfood/20'
                : 'bg-tag-both-bg border-tag-both/30 text-tag-both hover:bg-tag-both/15'"
            >
              {{ loc.exclude_default ? 'Hidden' : 'Visible' }}
            </button>
          </div>
        </div>

        <div v-if="addingChildOf === loc.id" class="ml-6 mt-1 flex gap-2 items-center">
          <Input v-model="childName" @keyup.enter="addChild(loc.id)" placeholder="Sub-location name..." class="text-xs" />
          <Button @click="addChild(loc.id)">Add</Button>
        </div>

        <div v-if="loc.children?.length" class="ml-4 mt-1 space-y-1">
          <div v-for="child in loc.children" :key="child.id" class="py-1.5 px-3 rounded-md bg-muted/50">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-1.5 text-sm text-muted-foreground min-w-0 flex-1">
                <ChevronRight class="h-3 w-3 shrink-0" />
                <input
                  v-if="editingId === child.id"
                  v-model="editingName"
                  @blur="saveEdit(child)"
                  @keyup.enter="saveEdit(child)"
                  @keyup.escape="editingId = null"
                  class="text-sm bg-background border border-input rounded px-1.5 py-0.5 w-full max-w-[10rem] focus:outline-none focus:ring-1 focus:ring-ring"
                  autofocus
                />
                <span
                  v-else
                  @click="startEdit(child)"
                  class="cursor-text hover:underline decoration-dashed underline-offset-4 truncate"
                >{{ formatLocation(child.name) }}</span>
              </div>
              <Button variant="ghost" size="icon" class="h-6 w-6 text-destructive hover:text-destructive shrink-0" @click="remove(child.id)" :aria-label="'Delete ' + formatLocation(child.name)">
                <Trash2 class="h-3 w-3" />
              </Button>
            </div>
            <div class="flex flex-wrap gap-1.5 mt-1 ml-4.5">
              <button
                @click="toggleFood(child)"
                class="text-xs font-medium px-2 py-0.5 rounded-full border transition-colors cursor-pointer"
                :class="child.is_food
                  ? 'bg-tag-food-bg border-tag-food/30 text-tag-food hover:bg-tag-food/20'
                  : 'bg-tag-both-bg border-tag-both/30 text-tag-both hover:bg-tag-both/15'"
              >
                {{ child.is_food ? 'Food' : 'Non-food' }}
              </button>
              <button
                @click="toggleExclude(child)"
                class="text-xs font-medium px-2 py-0.5 rounded-full border transition-colors cursor-pointer"
                :class="child.exclude_default
                  ? 'bg-tag-nonfood-bg border-tag-nonfood/30 text-tag-nonfood hover:bg-tag-nonfood/20'
                  : 'bg-tag-both-bg border-tag-both/30 text-tag-both hover:bg-tag-both/15'"
              >
                {{ child.exclude_default ? 'Hidden' : 'Visible' }}
              </button>
            </div>
          </div>
        </div>
      </div>
      <p v-if="!locations.length" class="text-sm text-muted-foreground py-2">No locations yet.</p>
    </div>
  </div>
</template>
