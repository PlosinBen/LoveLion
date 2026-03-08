<template>
  <div class="pb-4 flex items-center gap-3">
    <!-- Back Button -->
    <button
      v-if="showBack"
      @click="handleBack"
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
    <div class="shrink-0 flex items-center gap-1">
      <slot name="right" />
      <NuxtLink
        v-if="settingsTo"
        :to="settingsTo"
        class="w-10 h-10 text-neutral-400 flex items-center justify-center hover:bg-neutral-700 hover:text-white transition-colors no-underline"
      >
        <Icon icon="mdi:cog-outline" class="text-xl" />
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import SpaceSwitcher from './SpaceSwitcher.vue'

const props = defineProps({
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
  },
  backTo: {
    type: String,
    default: ''
  },
  settingsTo: {
    type: String,
    default: ''
  }
})

const router = useRouter()

const handleBack = () => {
  if (props.backTo) {
    router.push(props.backTo)
  } else {
    router.back()
  }
}
</script>
