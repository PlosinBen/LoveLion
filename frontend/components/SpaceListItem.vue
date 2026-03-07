<template>
  <div 
    class="bg-neutral-900 border border-neutral-800 p-4 rounded-2xl flex items-center justify-between hover:bg-neutral-800 transition-colors cursor-pointer group active:scale-95 shadow-sm"
    @click="$emit('click')"
  >
    <div class="flex items-center gap-4">
      <!-- Icon Container -->
      <div 
        class="w-12 h-12 rounded-xl flex items-center justify-center text-xl transition-colors border border-neutral-800" 
        :class="typeStyles.bg"
      >
        <Icon :icon="typeStyles.icon" :class="typeStyles.color" />
      </div>

      <!-- Info -->
      <div class="flex flex-col">
        <h3 class="font-bold text-white group-hover:text-indigo-400 transition-colors">{{ space.name }}</h3>
        <div class="flex items-center gap-2 mt-0.5">
            <span class="text-xs text-neutral-500">
                {{ typeStyles.label }}
            </span>
            <span class="w-1 h-1 rounded-full bg-neutral-800"></span>
            <span class="text-xs text-neutral-500 uppercase font-medium">
                {{ space.base_currency }}
            </span>
        </div>
      </div>
    </div>
    
    <!-- Pin Action -->
    <button 
      @click.stop="$emit('toggle-pin')" 
      class="w-10 h-10 rounded-full flex items-center justify-center transition-colors border-0 bg-transparent cursor-pointer"
      :class="space.is_pinned ? 'text-indigo-500' : 'text-neutral-700 hover:text-neutral-500'"
    >
      <Icon :icon="space.is_pinned ? 'mdi:pin' : 'mdi:pin-outline'" class="text-lg" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  space: {
    id: string
    name: string
    type: string
    base_currency: string
    is_pinned?: boolean
    [key: string]: any
  }
}>()

defineEmits(['click', 'toggle-pin'])

const typeStyles = computed(() => {
  if (props.space.type === 'trip' || props.space.type === 'project') {
    return {
      icon: 'mdi:airplane',
      bg: 'bg-neutral-800',
      color: 'text-indigo-400',
      label: '旅行專案'
    }
  }
  return {
    icon: 'mdi:wallet-outline',
    bg: 'bg-neutral-800',
    color: 'text-green-400',
    label: '個人空間'
  }
})
</script>
