<template>
  <div class="flex flex-col h-screen bg-black relative">
    <!-- Header / Close Button -->
    <div class="absolute top-0 left-0 right-0 p-4 flex justify-between items-center z-20 bg-gradient-to-b from-black/80 to-transparent">
       <button @click="router.back()" class="w-10 h-10 flex items-center justify-center rounded-full bg-black/50 text-white backdrop-blur-sm border border-white/20">
          <Icon icon="mdi:close" class="text-xl" />
       </button>
       <div class="flex gap-4">
           <!-- Optional Actions: Download, Share -->
           <a v-if="imageUrl" :href="imageUrl" download target="_blank" class="w-10 h-10 flex items-center justify-center rounded-full bg-black/50 text-white backdrop-blur-sm border border-white/20">
               <Icon icon="mdi:download" class="text-xl" />
           </a>
       </div>
    </div>

    <!-- Main Image -->
    <div class="flex-1 flex items-center justify-center p-4 overflow-hidden">
        <img v-if="imageUrl" :src="imageUrl" class="max-w-full max-h-full object-contain shadow-2xl" />
        <div v-else class="text-neutral-500">
            無法載入圖片
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const router = useRouter()
const route = useRoute()

const imageUrl = computed(() => route.query.src as string)
</script>
