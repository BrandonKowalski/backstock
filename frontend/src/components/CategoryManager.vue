<script setup lang="ts">
import { ref } from 'vue'
import type { Category } from '@/types'
import { useCategories } from '@/composables/useCategories'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Trash2 } from 'lucide-vue-next'

const props = defineProps<{ categories: Category[]; foodMode: boolean }>()
const emit = defineEmits<{ changed: [] }>()

const { createCategory, updateCategory, deleteCategory } = useCategories()
const newName = ref('')
const editingId = ref<number | null>(null)
const editingName = ref('')

async function add() {
  if (!newName.value.trim()) return
  await createCategory(newName.value.trim(), props.foodMode)
  newName.value = ''
  emit('changed')
}

function startEdit(cat: Category) {
  editingId.value = cat.id
  editingName.value = cat.name
}

async function saveEdit(cat: Category) {
  editingId.value = null
  const name = editingName.value.trim()
  if (!name || name === cat.name) return
  await updateCategory(cat.id, { name, is_food: cat.is_food })
  emit('changed')
}

async function toggleFood(cat: Category) {
  await updateCategory(cat.id, { name: cat.name, is_food: !cat.is_food })
  emit('changed')
}

async function remove(id: number) {
  await deleteCategory(id)
  emit('changed')
}
</script>

<template>
  <div>
    <div class="flex gap-2 items-center">
      <Input v-model="newName" @keyup.enter="add" placeholder="New category..." />
      <Button @click="add">Add</Button>
    </div>
    <div class="mt-3 space-y-1">
      <div v-for="cat in categories" :key="cat.id" class="flex items-center justify-between py-2 px-3 rounded-md bg-muted">
        <div class="flex items-center gap-2 min-w-0 flex-1">
          <input
            v-if="editingId === cat.id"
            v-model="editingName"
            @blur="saveEdit(cat)"
            @keyup.enter="saveEdit(cat)"
            @keyup.escape="editingId = null"
            class="text-sm bg-background border border-input rounded px-1.5 py-0.5 w-full max-w-[10rem] focus:outline-none focus:ring-1 focus:ring-ring"
            autofocus
          />
          <span
            v-else
            @click="startEdit(cat)"
            class="text-sm cursor-text hover:underline decoration-dashed underline-offset-4 truncate"
          >{{ cat.name }}</span>
          <button
            @click="toggleFood(cat)"
            class="text-xs font-medium px-2 py-0.5 rounded-full border transition-colors shrink-0 cursor-pointer"
            :class="cat.is_food
              ? 'bg-tag-food-bg border-tag-food/30 text-tag-food hover:bg-tag-food/20'
              : 'bg-tag-both-bg border-tag-both/30 text-tag-both hover:bg-tag-both/15'"
          >
            {{ cat.is_food ? 'Food' : 'Non-food' }}
          </button>
        </div>
        <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive hover:text-destructive shrink-0" @click="remove(cat.id)" :aria-label="'Delete ' + cat.name">
          <Trash2 class="h-3.5 w-3.5" />
        </Button>
      </div>
      <p v-if="!categories.length" class="text-sm text-muted-foreground py-2">No categories yet.</p>
    </div>
  </div>
</template>
