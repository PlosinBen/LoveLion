<template>
  <BaseCard
    class="flex items-center justify-between transition-none cursor-pointer group shadow-sm"
    @click="$emit('click')"
  >
    <div class="flex items-center gap-4">
      <!-- Cover Image or No-Image Icon -->
      <div v-if="space.cover_image" class="w-12 h-12 rounded-xl overflow-hidden shrink-0">
        <img :src="space.cover_image" class="w-full h-full object-cover" />
      </div>
      <div v-else class="w-12 h-12 rounded-xl flex items-center justify-center text-xl transition-none border border-neutral-800 bg-neutral-800 shrink-0">
        <Icon icon="mdi:image-off-outline" class="text-neutral-600" />
      </div>

      <!-- Info -->
      <div class="flex flex-col">
        <div class="flex items-center gap-1.5">
          <h3 class="font-bold text-white transition-none">{{ space.name }}</h3>
          <Icon v-if="sharingIcon" :icon="sharingIcon" class="text-lg text-neutral-500" />
        </div>
        <span class="text-xs text-neutral-500 uppercase font-medium mt-0.5">
            {{ space.base_currency }}
        </span>
      </div>
    </div>
    
    <!-- Pin Action -->
    <button 
      @click.stop="$emit('toggle-pin')" 
      class="w-10 h-10 rounded-xl flex items-center justify-center transition-all active:scale-90 border-0 cursor-pointer bg-transparent"
      :class="space.is_pinned ? 'text-indigo-500' : 'text-neutral-700 hover:text-neutral-500'"
    >
      <Icon :icon="space.is_pinned ? 'mdi:pin' : 'mdi:pin-outline'" class="text-lg" />
    </button>
  </BaseCard>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import BaseCard from '~/components/BaseCard.vue'

import { computed } from 'vue'

const props = defineProps<{
  space: {
    id: string
    name: string
    type: string
    base_currency: string
    is_pinned?: boolean
    my_role?: string
    member_count?: number
    [key: string]: any
  }
}>()

const sharingIcon = computed(() => {
  const { my_role, member_count } = props.space
  if (my_role === 'owner' && member_count && member_count > 1) return 'mdi:account-multiple-outline'
  if (my_role === 'member') return 'mdi:account-arrow-left-outline'
  return null
})

defineEmits(['click', 'toggle-pin'])
</script>
