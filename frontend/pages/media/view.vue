<template>
  <div class="flex flex-col h-screen bg-black relative">
    <!-- Header / Close Button -->
    <div class="absolute top-0 left-0 right-0 p-4 pt-8 flex justify-between items-center z-20 bg-gradient-to-b from-black/80 to-transparent">
       <button @click="router.back()" class="w-10 h-10 flex items-center justify-center rounded-full bg-black/50 text-white backdrop-blur-md border border-white/10 cursor-pointer">
          <Icon icon="mdi:close" class="text-2xl" />
       </button>
       <div class="flex gap-4">
           <a v-if="imageUrl" :href="imageUrl" download target="_blank" class="w-10 h-10 flex items-center justify-center rounded-full bg-black/50 text-white backdrop-blur-md border border-white/10 cursor-pointer no-underline">
               <Icon icon="mdi:download" class="text-xl" />
           </a>
       </div>
    </div>

    <!-- Main Image Display -->
    <div class="flex-1 flex items-center justify-center p-4 overflow-hidden">
        <img v-if="imageUrl" :src="imageUrl" class="max-w-full max-h-full object-contain shadow-2xl" />
        <div v-else class="text-neutral-600 font-bold">
            找不到媒體資源
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

definePageMeta({
  layout: 'empty'
})

const router = useRouter()
const route = useRoute()

const imageUrl = computed(() => route.query.src as string)
</script>
