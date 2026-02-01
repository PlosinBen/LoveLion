<template>
  <div class="image-preview fixed inset-0 bg-black z-50 flex flex-col h-[100dvh]">
    <!-- Header / Close -->
    <header class="absolute top-0 left-0 right-0 p-4 pt-6 transition-opacity duration-300 pointer-events-none bg-gradient-to-b from-black/60 to-transparent flex justify-start items-center z-10" :class="showControls ? 'opacity-100' : 'opacity-0'">
       <!-- Progress Bar -->
       <div class="absolute top-2 left-2 right-2 flex gap-1">
           <div v-for="(img, idx) in images" :key="idx" 
                class="h-1 flex-1 rounded-full transition-colors duration-300 backdrop-blur-sm shadow-sm"
                :class="idx === currentIndex ? 'bg-white/90 shadow-[0_0_8px_rgba(255,255,255,0.6)]' : 'bg-white/40'">
           </div>
       </div>

       <button @click="router.back()" class="w-10 h-10 flex justify-center items-center rounded-full bg-black/40 text-white backdrop-blur-md border-0 pointer-events-auto cursor-pointer mt-2">
          <Icon icon="mdi:close" class="text-2xl" />
       </button>
    </header>

    <!-- Image Container -->
    <div class="flex-1 flex justify-center items-center p-2 overflow-hidden relative w-full" 
         @touchstart="handleTouchStart" 
         @touchend="handleTouchEnd"
         @click="toggleControls">
        
        <!-- Prev Button -->
        <button 
             class="absolute left-4 top-1/2 -translate-y-1/2 z-10 rounded-full p-2 text-white/80 bg-black/30 backdrop-blur-sm transition-all duration-300 hover:bg-black/50 border-0 cursor-pointer"
             :class="{ 'opacity-100 pointer-events-auto': currentIndex > 0 && showControls, 'opacity-0 pointer-events-none': currentIndex === 0 || !showControls }"
             @click.stop="prev">
             <Icon icon="mdi:chevron-left" class="text-4xl" />
        </button>
        
        <!-- Next Button -->
        <button
             class="absolute right-4 top-1/2 -translate-y-1/2 z-10 rounded-full p-2 text-white/80 bg-black/30 backdrop-blur-sm transition-all duration-300 hover:bg-black/50 border-0 cursor-pointer"
             :class="{ 'opacity-100 pointer-events-auto': currentIndex < images.length - 1 && showControls, 'opacity-0 pointer-events-none': currentIndex === images.length - 1 || !showControls }"
             @click.stop="next">
             <Icon icon="mdi:chevron-right" class="text-4xl" />
        </button>

        <!-- Loading -->
        <div v-if="loading || isSwitching" class="text-neutral-400 flex flex-col items-center gap-2 absolute z-10 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
            <Icon icon="mdi:loading" class="text-4xl animate-spin" />
            <span class="text-sm">載入中...</span>
        </div>

        <template v-if="currentImage">
            <!-- BlurHash Placeholder -->
            <div v-if="currentImage.hash && isSwitching" class="absolute inset-0 flex justify-center items-center p-2">
                 <BlurHashCanvas :hash="currentImage.hash" :width="32" :height="32" class="max-w-full max-h-full object-contain" />
            </div>

             <img 
                :key="currentImage.url"
                :src="currentImage.url" 
                class="max-w-full max-h-full object-contain shadow-2xl transition-opacity duration-500"
                :class="{ 'opacity-0': isSwitching, 'opacity-100': !isSwitching }"
                alt="Preview"
                @load="handleImageLoad"
            />
            <!-- Filmstrip -->
            <div class="absolute bottom-6 left-0 right-0 h-16 flex justify-center gap-2 overflow-x-auto px-4 py-1 z-10 transition-opacity duration-300 scrollbar-hide"
                 :class="showControls ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'">
                 <div v-for="(img, idx) in images" :key="idx"
                      class="h-full aspect-square rounded-md overflow-hidden border-2 transition-all cursor-pointer box-content bg-neutral-900"
                      :class="idx === currentIndex ? 'border-white scale-110' : 'border-transparent opacity-60 hover:opacity-100'"
                      @click.stop="currentIndex = idx; isSwitching = true">
                      <BlurHashCanvas v-if="img.hash" :hash="img.hash" class="w-full h-full object-cover" />
                      <img v-else :src="img.url" class="w-full h-full object-cover" loading="lazy" />
                 </div>
            </div>
        </template>

        <div v-else-if="!loading" class="text-neutral-500">
            無效的圖片連結
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { computed, onMounted, ref } from 'vue'
import { useImages } from '~/composables/useImages'
import BlurHashCanvas from '~/components/BlurHashCanvas.vue'

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
const currentIndex = ref(0)
const loading = ref(true)

const isSwitching = ref(true)
const showControls = ref(true)
const currentImage = computed(() => images.value[currentIndex.value])

const toggleControls = () => {
    showControls.value = !showControls.value
}

const handleImageLoad = () => {
    isSwitching.value = false
}

const next = () => {
    if (currentIndex.value < images.value.length - 1) {
        isSwitching.value = true
        currentIndex.value++
    }
}

const prev = () => {
    if (currentIndex.value > 0) {
        isSwitching.value = true
        currentIndex.value--
    }
}

// Simple swipe detection
const touchStartX = ref(0)
const handleTouchStart = (e: TouchEvent) => {
    if (e.touches && e.touches.length > 0) {
        touchStartX.value = e.touches[0]?.clientX ?? 0
    }
}
const handleTouchEnd = (e: TouchEvent) => {
    if (e.changedTouches && e.changedTouches.length > 0) {
        const touchEndX = e.changedTouches[0]?.clientX ?? 0
        const diff = touchStartX.value - touchEndX
        if (Math.abs(diff) > 50) { // Threshold
            if (diff > 0) next() // Swipe Left -> Next
            else prev() // Swipe Right -> Prev
        }
    }
}

onMounted(async () => {
    loading.value = true
    const indexParam = route.query.index
    
    // 1. Try to get from History State (Instant)
    const state = history.state as { images?: ImageItem[], urls?: string[] }
    
    if (state && state.images && Array.isArray(state.images) && state.images.length > 0) {
        images.value = state.images
        loading.value = false
    } else if (state && state.urls && Array.isArray(state.urls) && state.urls.length > 0) {
        // Legacy state support
        images.value = state.urls.map(url => ({ url }))
        loading.value = false
    } else {
        // 2. Fallback: Fetch from API using ID
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
             // Legacy single url fallback
             images.value = [{ url: route.query.url as string }]
        }
        
        if (images.value.length === 0) {
            isSwitching.value = false
        }
        loading.value = false
    }

})
</script>

<style scoped>
/* Ensure it covers everything including potential safe areas */
</style>
