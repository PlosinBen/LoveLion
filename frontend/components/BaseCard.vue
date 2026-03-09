<template>
  <component
    :is="to ? resolveComponent('NuxtLink') : tag"
    :to="to"
    class="block bg-neutral-900 border border-neutral-800 rounded-2xl"
    :class="[
      to && 'no-underline text-inherit',
      padding,
      $attrs.class
    ]"
    v-bind="{ ...$attrs, class: undefined }"
  >
    <slot />
  </component>
</template>

<script setup lang="ts">
import { resolveComponent } from 'vue'

withDefaults(defineProps<{
  to?: string
  tag?: string
  padding?: string
}>(), {
  tag: 'div',
  padding: 'p-4'
})

defineOptions({ inheritAttrs: false })
</script>
