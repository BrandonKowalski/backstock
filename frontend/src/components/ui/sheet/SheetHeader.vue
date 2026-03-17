<script setup lang="ts">
import { cn } from '@/lib/utils'
withDefaults(defineProps<{ class?: string }>(), {})
</script>

<template>
  <div :class="cn('flex flex-col space-y-2 p-6 pb-0', $props.class)">
    <slot />
  </div>
</template>
