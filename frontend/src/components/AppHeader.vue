<script setup lang="ts">
import { Settings } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

defineEmits<{ 'open-settings': [] }>()
</script>

<template>
  <header class="sticky top-0 z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <div class="max-w-2xl mx-auto px-4 h-14 flex items-center justify-between">
      <h1 class="text-xl font-bold tracking-tight">Backstock</h1>
      <Button variant="ghost" size="icon" @click="$emit('open-settings')" aria-label="Settings">
        <Settings class="h-5 w-5" />
      </Button>
    </div>
  </header>
</template>
