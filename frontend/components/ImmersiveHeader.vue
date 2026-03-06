<template>
  <div class="relative w-full overflow-hidden shrink-0 bg-neutral-800" :class="[height]">
    <!-- Background: Image or Fallback -->
    <div class="absolute inset-0 bg-neutral-900 overflow-hidden">
      <template v-if="image">
        <img :src="imageUrl" class="w-full h-full object-cover" :alt="alt" />
        <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-transparent to-transparent"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div>
      </template>
      <div v-else class="w-full h-full bg-gradient-to-br from-indigo-900 to-neutral-900 flex items-center justify-center opacity-50">
        <Icon :icon="fallbackIcon" class="text-6xl text-white/20" />
      </div>
    </div>

    <!-- Content Overlay -->
    <div class="absolute inset-0 flex flex-col justify-between p-4 z-10">
      <!-- Top Row -->
      <header class="flex justify-between items-start">
        <div class="flex items-center gap-2">
           <slot name="top-left"></slot>
        </div>
        <div class="flex items-center gap-2">
           <slot name="top-right"></slot>
        </div>
      </header>

      <!-- Bottom Row -->
      <div class="w-full">
         <slot name="bottom"></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useImages } from '~/composables/useImages'

const props = defineProps({
  image: {
    type: String,
    default: null
  },
  fallbackIcon: {
    type: String,
    default: 'mdi:image'
  },
  height: {
    type: String,
    default: 'h-64'
  },
  alt: {
    type: String,
    default: 'Header Image'
  }
})

const { getImageUrl } = useImages()

const imageUrl = computed(() => {
    if (!props.image) return ''
    return props.image.startsWith('blob:') ? props.image : getImageUrl(props.image)
})
</script>
