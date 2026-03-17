<script setup lang="ts">
import { computed } from 'vue'
import { buttonVariants, type ButtonVariants } from '.'
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  variant?: ButtonVariants['variant']
  size?: ButtonVariants['size']
  as?: string
  class?: string
  disabled?: boolean
}>(), { as: 'button' })
const classes = computed(() => cn(buttonVariants({ variant: props.variant, size: props.size }), props.class))
</script>

<template>
  <component :is="as" :class="classes" :disabled="disabled">
    <slot />
  </component>
</template>
