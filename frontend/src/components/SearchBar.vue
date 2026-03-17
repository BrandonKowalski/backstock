<script setup lang="ts">
import { Search } from 'lucide-vue-next'
import { Input } from '@/components/ui/input'

const model = defineModel<string>()
const emit = defineEmits<{ search: [] }>()

let timeout: ReturnType<typeof setTimeout>
function onInput(e: Event) {
  model.value = (e.target as HTMLInputElement).value
  clearTimeout(timeout)
  timeout = setTimeout(() => emit('search'), 300)
}
</script>

<template>
  <div class="mt-4 relative">
    <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
    <Input
      type="text"
      :model-value="model"
      @input="onInput"
      placeholder="Search items..."
      class="pl-9"
    />
  </div>
</template>
