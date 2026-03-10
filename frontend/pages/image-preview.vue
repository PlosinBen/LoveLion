<template>
  <div class="image-preview fixed inset-0 bg-black flex flex-col h-screen z-50">
    <!-- Header / Close -->
    <header class="absolute top-0 left-0 right-0 p-4 pt-6 transition-opacity duration-300 pointer-events-none bg-gradient-to-b from-black/60 to-transparent flex justify-start items-center z-10" :class="showControls ? 'opacity-100' : 'opacity-0'">
       <!-- Progress Bar -->
       <div class="absolute top-2 left-2 right-2 flex gap-1">
           <div v-for="(img, idx) in images" :key="idx" 
                class="h-1 flex-1 rounded-full transition-colors duration-300 backdrop-blur-sm shadow-sm"
                :class="idx === currentIndex ? 'bg-white/90 shadow-white/60' : 'bg-white/40'">
           </div>
       </div>

       <button @click="router.back()" class="flex justify-center items-center w-10 h-10 mt-2 rounded-full bg-black/40 text-white backdrop-blur-md pointer-events-auto border-0 cursor-pointer transition-colors hover:bg-black/60 active:scale-95">
          <Icon icon="mdi:close" class="text-2xl" />
       </button>
    </header>

    <!-- Image Container (Carousel) -->
    <div class="flex-1 flex items-center overflow-hidden relative w-full"
         @touchstart="handleTouchStart"
         @touchmove="handleTouchMove"
         @touchend="handleTouchEnd"
         @click="toggleControls">
        
        <!-- Navigation Buttons (Desktop) -->
        <button
             v-if="currentIndex > 0"
             class="flex justify-center items-center w-10 h-10 absolute left-4 top-1/2 -translate-y-1/2 z-10 rounded-full bg-black/30 text-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-black/50 border-0 cursor-pointer hidden md:flex active:scale-95"
             :class="{ 'opacity-100': showControls, 'opacity-0': !showControls }"
             @click.stop="prev">
             <Icon icon="mdi:chevron-left" class="text-4xl" />
        </button>

        <button
             v-if="currentIndex < images.length - 1"
             class="flex justify-center items-center w-10 h-10 absolute right-4 top-1/2 -translate-y-1/2 z-10 rounded-full bg-black/30 text-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-black/50 border-0 cursor-pointer hidden md:flex active:scale-95"
             :class="{ 'opacity-100': showControls, 'opacity-0': !showControls }"
             @click.stop="next">
             <Icon icon="mdi:chevron-right" class="text-4xl" />
        </button>

        <!-- Carousel Track -->
        <div class="relative w-full h-full flex">
             <div v-for="(img, idx) in images" :key="idx"
                  class="absolute inset-0 flex justify-center items-center p-2 will-change-transform"
                  :class="isDragging ? 'transition-none' : 'transition-transform duration-300 ease-out'"
                  :style="getImageStyle(idx)">
                  
                  <template v-if="Math.abs(currentIndex - idx) <= 1">
                      <img 
                          :src="img.url" 
                          class="max-w-full max-h-full object-contain shadow-2xl"
                          draggable="false"
                      />
                  </template>
             </div>
        </div>

        <!-- Loading Overlay -->
         <div v-if="loading" class="text-neutral-400 flex flex-col items-center gap-2 absolute z-10 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 pointer-events-none">
            <Icon icon="mdi:loading" class="text-4xl animate-spin" />
            <span class="text-sm font-bold">載入中...</span>
        </div>
    </div>

    <!-- Filmstrip -->
    <div v-if="images.length > 1"
         class="absolute bottom-10 left-0 right-0 h-16 flex justify-center gap-2 overflow-x-auto px-4 py-1 z-10 transition-opacity duration-300"
         :class="showControls ? 'opacity-100' : 'opacity-0'">
         <div v-for="(img, idx) in images" :key="idx"
              class="h-full aspect-square rounded-lg overflow-hidden border-2 transition-all cursor-pointer bg-neutral-900 shrink-0"
              :class="idx === currentIndex ? 'border-indigo-500 scale-110 shadow-lg' : 'border-transparent opacity-40'"
              @click.stop="currentIndex = idx">
              <img :src="img.url" class="w-full h-full object-cover" />
         </div>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && images.length === 0" class="flex-1 flex justify-center items-center text-neutral-500 font-bold">
        暫無預覽圖片
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useImages } from '~/composables/useImages'

definePageMeta({
  layout: 'empty'
})

const router = useRouter()
const route = useRoute()
const { getImages, getImageUrl } = useImages()

interface ImageItem {
    url: string
    hash?: string
}

const images = ref<ImageItem[]>([])
const currentIndex = ref(Number(route.query.index) || 0)
const loading = ref(true)
const showControls = ref(true)

const toggleControls = () => {
    showControls.value = !showControls.value
}

const next = () => {
    if (currentIndex.value < images.value.length - 1) {
        currentIndex.value++
    }
}

const prev = () => {
    if (currentIndex.value > 0) {
        currentIndex.value--
    }
}

// Swipe handling
const touchStartX = ref(0)
const currentTranslateX = ref(0)
const isDragging = ref(false)

const getImageStyle = (idx: number) => {
    const offsetPercent = (idx - currentIndex.value) * 100
    return {
        transform: `translateX(calc(${offsetPercent}% + ${currentTranslateX.value}px))`
    }
}

const handleTouchStart = (e: TouchEvent) => {
    touchStartX.value = e.touches[0]?.clientX ?? 0
    isDragging.value = true
    currentTranslateX.value = 0
}

const handleTouchMove = (e: TouchEvent) => {
    if (!isDragging.value) return
    const currentX = e.touches[0]?.clientX ?? 0
    currentTranslateX.value = currentX - touchStartX.value
}

const handleTouchEnd = () => {
    isDragging.value = false
    const threshold = window.innerWidth * 0.2
    
    if (Math.abs(currentTranslateX.value) > threshold) {
        if (currentTranslateX.value > 0) prev()
        else next()
    }
    currentTranslateX.value = 0
}

onMounted(async () => {
    loading.value = true
    
    // Attempt to get images from history state or query params
    const entityId = route.query.id as string
    const entityType = route.query.type as string
    
    if (entityId && entityType) {
        try {
            const apiImages = await getImages(entityId, entityType)
            images.value = apiImages.map(img => ({
                url: getImageUrl(img.file_path),
                hash: img.blur_hash
            }))
        } catch (e) {
            console.error("Failed to fetch preview images", e)
        }
    } else if (route.query.url) {
         images.value = [{ url: route.query.url as string }]
    }
    
    loading.value = false
})
</script>
