<template>
  <NuxtLink
    :to="to"
    :class="[useButtonStyle(variant), $attrs.class]"
  >
    <slot />
  </NuxtLink>
</template>

<script setup lang="ts">
import { useButtonStyle } from '~/composables/useButtonStyle'

interface Props {
  to: string
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
}

withDefaults(defineProps<Props>(), {
  variant: 'primary',
})
</script>
