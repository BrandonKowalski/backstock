<script setup lang="ts">
import { ref } from 'vue'
import type { Unit } from '@/types'
import { useUnits } from '@/composables/useUnits'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Trash2 } from 'lucide-vue-next'

defineProps<{ units: Unit[] }>()
const emit = defineEmits<{ changed: [] }>()

const { createUnit, updateUnit, deleteUnit } = useUnits()
const newName = ref('')
const newAbbr = ref('')
const editingId = ref<number | null>(null)
const editingName = ref('')
const editingAbbr = ref('')

async function add() {
  if (!newName.value.trim() || !newAbbr.value.trim()) return
  await createUnit(newName.value.trim(), newAbbr.value.trim())
  newName.value = ''
  newAbbr.value = ''
  emit('changed')
}

function startEdit(u: Unit) {
  editingId.value = u.id
  editingName.value = u.name
  editingAbbr.value = u.abbreviation
}

async function saveEdit(u: Unit) {
  editingId.value = null
  const name = editingName.value.trim()
  const abbr = editingAbbr.value.trim()
  if (!name || !abbr) return
  if (name === u.name && abbr === u.abbreviation) return
  await updateUnit(u.id, { name, abbreviation: abbr, is_food: u.is_food })
  emit('changed')
}

async function cycleFood(u: Unit) {
  // null (Both) -> true (Food) -> false (Non-food) -> null (Both)
  let next: boolean | null
  if (u.is_food === null) next = true
  else if (u.is_food === true) next = false
  else next = null
  await updateUnit(u.id, { name: u.name, abbreviation: u.abbreviation, is_food: next })
  emit('changed')
}

function foodLabel(isFood: boolean | null): string {
  if (isFood === null) return 'Both'
  return isFood ? 'Food' : 'Non-food'
}

function foodClass(isFood: boolean | null): string {
  if (isFood === null) return 'bg-tag-both-bg border-tag-both/30 text-tag-both hover:bg-tag-both/15'
  if (isFood) return 'bg-tag-food-bg border-tag-food/30 text-tag-food hover:bg-tag-food/20'
  return 'bg-tag-nonfood-bg border-tag-nonfood/30 text-tag-nonfood hover:bg-tag-nonfood/20'
}

async function remove(id: number) {
  await deleteUnit(id)
  emit('changed')
}
</script>

<template>
  <div>
    <div class="flex gap-2 items-center">
      <Input v-model="newName" @keyup.enter="add" placeholder="Name..." />
      <Input v-model="newAbbr" @keyup.enter="add" placeholder="Abbr..." class="w-20" />
      <Button @click="add">Add</Button>
    </div>
    <div class="mt-3 space-y-1">
      <div v-for="u in units" :key="u.id" class="flex items-center justify-between py-2 px-3 rounded-md bg-muted">
        <div class="flex items-center gap-2 min-w-0 flex-1">
          <div v-if="editingId === u.id" class="flex items-center gap-2">
            <input
              v-model="editingName"
              @keyup.enter="saveEdit(u)"
              @keyup.escape="editingId = null"
              class="text-sm bg-background border border-input rounded px-1.5 py-0.5 w-full max-w-[8rem] focus:outline-none focus:ring-1 focus:ring-ring"
              autofocus
            />
            <input
              v-model="editingAbbr"
              @blur="saveEdit(u)"
              @keyup.enter="saveEdit(u)"
              @keyup.escape="editingId = null"
              class="text-sm bg-background border border-input rounded px-1.5 py-0.5 w-16 focus:outline-none focus:ring-1 focus:ring-ring"
            />
          </div>
          <span
            v-else
            @click="startEdit(u)"
            class="text-sm cursor-text hover:underline decoration-dashed underline-offset-4"
          >{{ u.name }} <span class="text-muted-foreground">({{ u.abbreviation }})</span></span>
          <button
            @click="cycleFood(u)"
            class="text-xs font-medium px-2 py-0.5 rounded-full border transition-colors shrink-0 cursor-pointer"
            :class="foodClass(u.is_food)"
          >
            {{ foodLabel(u.is_food) }}
          </button>
        </div>
        <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:text-destructive shrink-0" @click="remove(u.id)">
          <Trash2 class="h-3.5 w-3.5" />
        </Button>
      </div>
      <p v-if="!units.length" class="text-sm text-muted-foreground py-2">No units yet.</p>
    </div>
  </div>
</template>
