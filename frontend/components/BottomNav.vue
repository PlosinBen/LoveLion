<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-neutral-900 border-t border-neutral-800 flex justify-around items-center px-4 py-3 z-50">
    <template v-for="item in navItems" :key="item.to || item.id">
      <!-- Link Mode -->
      <NuxtLink
        v-if="item.to"
        :to="item.to"
        class="flex flex-col items-center gap-1 no-underline transition-transform active:scale-95"
        :class="isActive(item) ? 'text-indigo-400' : 'text-neutral-500'"
      >
        <Icon :icon="item.icon" class="text-2xl" />
        <span class="text-xs font-bold uppercase tracking-widest scale-90">
          {{ item.label }}
        </span>
      </NuxtLink>

      <!-- Action Mode (Button) -->
      <button
        v-else
        @click="$emit('update:modelValue', item.id)"
        class="flex flex-col items-center gap-1 bg-transparent border-0 cursor-pointer transition-transform active:scale-95"
        :class="isActive(item) ? 'text-indigo-400' : 'text-neutral-500'"
      >
        <Icon :icon="item.icon" class="text-2xl" />
        <span class="text-xs font-bold uppercase tracking-widest scale-90">
          {{ item.label }}
        </span>
      </button>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

interface NavItem {
  id?: string
  label: string
  icon: string
  to?: string
}

const props = defineProps<{
  modelValue?: string
  items?: NavItem[]
}>()

defineEmits(['update:modelValue'])

const route = useRoute()

const navItems = computed(() => {
  if (props.items) return props.items
  return [
    { label: '空間', icon: 'mdi:view-grid-outline', to: '/' },
    { label: '設定', icon: 'mdi:cog-outline', to: '/settings' }
  ]
})

const isActive = (item: NavItem & { alternateTo?: string }) => {
  if (!item.to) return item.id === props.modelValue
  
  if (item.to === '/') return route.path === '/'
  
  const checkMatch = (path: string) => {
    const isExact = route.path === path
    const isSubpath = route.path.startsWith(path + '/')
    
    if (isExact) return true
    
    if (isSubpath) {
      // Check if there is another nav item that is a better match (longer)
      const betterMatch = navItems.value.some(other => {
        if (!other.to || other.to === path) return false
        return (route.path === other.to || route.path.startsWith(other.to + '/')) && other.to.length > path.length
      })
      return !betterMatch
    }
    return false
  }

  const primaryMatch = checkMatch(item.to)
  const alternateMatch = item.alternateTo ? checkMatch(item.alternateTo) : false
  
  return primaryMatch || alternateMatch
}
</script>
