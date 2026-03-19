<template>
  <div class="pb-4 flex flex-col gap-1">
    <!-- Breadcrumbs -->
    <nav v-if="breadcrumbs?.length" class="flex items-center gap-1 text-xs text-neutral-500 px-1">
      <template v-for="(crumb, i) in breadcrumbs" :key="i">
        <NuxtLink v-if="crumb.to" :to="crumb.to" class="hover:text-neutral-300 transition-colors no-underline text-neutral-500">
          {{ crumb.label }}
        </NuxtLink>
        <span v-else>{{ crumb.label }}</span>
        <Icon v-if="i < breadcrumbs.length - 1" icon="mdi:chevron-right" class="text-neutral-700 text-sm" />
      </template>
    </nav>

    <!-- Header Row -->
    <div class="flex items-center gap-3">
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
        <div v-else class="flex flex-col">
          <h1 class="text-xl font-bold text-white tracking-tight truncate">{{ title }}</h1>
          <slot name="subtitle" />
        </div>
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
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import SpaceSwitcher from './SpaceSwitcher.vue'

interface Breadcrumb {
  label: string
  to?: string
}

const props = withDefaults(defineProps<{
  title: string
  showBack?: boolean
  showSwitcher?: boolean
  backTo?: string
  settingsTo?: string
  breadcrumbs?: Breadcrumb[]
}>(), {
  showBack: true
})

// Sync browser tab title
const documentTitle = computed(() => {
  const crumbPath = props.breadcrumbs?.map(c => c.label).join(' › ')
  return crumbPath ? `${props.title} - ${crumbPath} | LoveLion` : `${props.title} | LoveLion`
})

useHead({ title: documentTitle })

const router = useRouter()

const handleBack = () => {
  if (props.backTo) {
    router.push(props.backTo)
  } else {
    router.back()
  }
}
</script>
