<template>
  <div 
    class="bg-neutral-900 border border-neutral-800/60 p-5 rounded-3xl flex items-center justify-between hover:bg-neutral-800/50 hover:border-indigo-500/30 transition-all cursor-pointer group active:scale-[0.98] shadow-sm"
    @click="$emit('click')"
  >
    <div class="flex items-center gap-4">
      <!-- Icon Container -->
      <div 
        class="w-14 h-14 rounded-2xl flex items-center justify-center text-2xl transition-all border border-neutral-800 group-hover:border-indigo-500/20" 
        :class="typeStyles.bg"
      >
        <Icon :icon="typeStyles.icon" :class="typeStyles.color" />
      </div>

      <!-- Info -->
      <div class="flex flex-col">
        <h3 class="font-black text-white group-hover:text-indigo-400 transition-colors tracking-tight text-lg">{{ space.name }}</h3>
        <div class="flex items-center gap-2 mt-1">
            <span class="text-[10px] font-black uppercase tracking-[0.15em] text-neutral-500">
                {{ typeStyles.label }}
            </span>
            <span class="w-1 h-1 rounded-full bg-neutral-800"></span>
            <span class="text-[10px] font-black text-neutral-600 uppercase tracking-widest">
                {{ space.base_currency }}
            </span>
        </div>
      </div>
    </div>
    
    <!-- Pin Action -->
    <button 
      @click.stop="$emit('toggle-pin')" 
      class="w-10 h-10 rounded-xl flex items-center justify-center transition-all border-0 bg-transparent cursor-pointer hover:bg-indigo-500/10"
      :class="space.is_pinned ? 'text-indigo-400' : 'text-neutral-700 hover:text-neutral-500'"
    >
      <Icon :icon="space.is_pinned ? 'mdi:pin' : 'mdi:pin-outline'" class="text-xl" />
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
  // Map legacy 'trip' to 'project' visual style, 'personal' as default
  if (props.space.type === 'trip' || props.space.type === 'project') {
    return {
      icon: 'mdi:airplane-takeoff',
      bg: 'bg-blue-500/5',
      color: 'text-blue-400',
      label: '旅遊專案'
    }
  }
  return {
    icon: 'mdi:wallet-outline',
    bg: 'bg-emerald-500/5',
    color: 'text-emerald-400',
    label: '個人空間'
  }
})
</script>
