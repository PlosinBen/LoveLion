<template>
  <div class="relative group aspect-square bg-gray-100 rounded-lg overflow-hidden border border-gray-200">
    <img 
      :src="image.isPending ? image.file_path : getImageUrl(image.file_path)" 
      class="w-full h-full object-cover"
      alt="Uploaded image"
    />
    
    <!-- Overlay Actions -->
    <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
       <button type="button" @click.stop="$emit('delete', index)" title="Delete" class="p-2 rounded bg-red-500/80 text-white border-0 cursor-pointer hover:bg-red-600 transition-colors">
         <Icon icon="mdi:trash-can-outline" class="h-4 w-4" />
       </button>

       <div v-if="allowReorder" class="flex gap-1">
         <button type="button" v-if="!isFirst" @click.stop="$emit('move', index, -1)" class="p-2 rounded bg-white text-black border-0 cursor-pointer hover:bg-neutral-200 transition-colors">
            <Icon icon="mdi:chevron-left" class="h-4 w-4" />
         </button>
         <button type="button" v-if="!isLast" @click.stop="$emit('move', index, 1)" class="p-2 rounded bg-white text-black border-0 cursor-pointer hover:bg-neutral-200 transition-colors">
            <Icon icon="mdi:chevron-right" class="h-4 w-4" />
         </button>
       </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useImages } from '~/composables/useImages'

const props = defineProps<{
  image: any
  index: number
  allowReorder?: boolean
  isFirst: boolean
  isLast: boolean
}>()

defineEmits(['delete', 'move'])

const { getImageUrl } = useImages()
</script>
