<template>
  <div class="px-2 pt-0 pb-4 flex items-center gap-3">
    <!-- Back Button -->
    <button 
      v-if="showBack" 
      @click="router.back()" 
      class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer shrink-0"
    >
      <Icon icon="mdi:arrow-left" class="text-xl" />
    </button>

    <!-- Title or Switcher -->
    <div class="flex-1 min-w-0">
      <SpaceSwitcher v-if="showSwitcher" />
      <h1 v-else class="text-xl font-bold text-white tracking-tight truncate">{{ title }}</h1>
    </div>

    <!-- Right Actions -->
    <div class="shrink-0">
      <slot name="right" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import SpaceSwitcher from './SpaceSwitcher.vue'

defineProps({
  title: {
    type: String,
    required: true
  },
  showBack: {
    type: Boolean,
    default: true
  },
  showSwitcher: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()
</script>
