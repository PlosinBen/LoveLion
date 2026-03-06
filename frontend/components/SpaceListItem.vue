<template>
  <div 
    class="bg-neutral-900 border border-neutral-800/50 py-5 rounded-3xl flex items-center justify-between hover:bg-neutral-800 transition-all cursor-pointer group active:scale-95"
    @click="$emit('click')"
  >
    <div class="flex items-center gap-4">
      <!-- Icon Container -->
      <div 
        class="w-12 h-12 rounded-2xl flex items-center justify-center text-2xl transition-colors" 
        :class="space.type === 'trip' ? 'bg-blue-500/10 text-blue-400' : 'bg-green-500/10 text-green-400'"
      >
        <Icon :icon="space.type === 'trip' ? 'mdi:airplane' : 'mdi:wallet-outline'" />
      </div>

      <!-- Info -->
      <div class="flex flex-col">
        <h3 class="font-bold text-white group-hover:text-indigo-400 transition-colors">{{ space.name }}</h3>
        <div class="flex items-center gap-2 mt-0.5">
            <span class="text-xs font-bold uppercase tracking-wider text-neutral-500">
                {{ space.type === 'trip' ? '??撠?' : '?犖蝛粹?' }}
            </span>
            <span class="w-1 h-1 rounded-full bg-neutral-800"></span>
            <span class="text-xs font-bold text-neutral-500 uppercase tracking-wider">
                {{ space.base_currency }}
            </span>
        </div>
      </div>
    </div>
    
    <!-- Pin Action -->
    <button 
      @click.stop="$emit('toggle-pin')" 
      class="w-10 h-10 rounded-full flex items-center justify-center transition-all border-0 bg-transparent cursor-pointer"
      :class="space.is_pinned ? 'text-indigo-400' : 'text-neutral-700 hover:text-neutral-400'"
    >
      <Icon :icon="space.is_pinned ? 'mdi:pin' : 'mdi:pin-outline'" class="text-xl" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{
  space: {
    id: string
    name: string
    type: string
    base_currency: string
    is_pinned: boolean
    [key: string]: any
  }
}>()

defineEmits(['click', 'toggle-pin'])
</script>
