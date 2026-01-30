<template>
  <div class="image-preview fixed inset-0 bg-black z-50 flex flex-col h-[100dvh]">
    <!-- Header / Close -->
    <header class="absolute top-0 left-0 right-0 p-4 flex justify-between items-center z-50 bg-gradient-to-b from-black/50 to-transparent">
       <button @click="router.back()" class="w-10 h-10 flex justify-center items-center rounded-full bg-black/40 text-white backdrop-blur-md border-0">
          <Icon icon="mdi:close" class="text-2xl" />
       </button>
    </header>

    <!-- Image Container -->
    <div class="flex-1 flex justify-center items-center p-2 overflow-hidden relative w-full" 
         @touchstart="handleTouchStart" 
         @touchend="handleTouchEnd">
        
        <!-- Prev Tap Area -->
        <div class="absolute left-0 top-0 bottom-0 w-1/4 z-20" @click.stop="prev"></div>
        <!-- Next Tap Area -->
        <div class="absolute right-0 top-0 bottom-0 w-1/4 z-20" @click.stop="next"></div>

        <!-- Loading -->
        <div v-if="loading || isSwitching" class="text-neutral-400 flex flex-col items-center gap-2 absolute z-10">
            <Icon icon="mdi:loading" class="text-4xl animate-spin" />
            <span class="text-sm">載入中...</span>
        </div>

        <template v-if="currentUrl">
             <img 
                :key="currentUrl"
                :src="currentUrl" 
                class="max-w-full max-h-full object-contain shadow-2xl transition-opacity duration-300"
                :class="{ 'opacity-0': isSwitching, 'opacity-100': !isSwitching }"
                alt="Preview"
                @load="handleImageLoad"
            />
            <!-- Counter -->
            <div class="absolute bottom-8 left-1/2 -translate-x-1/2 bg-neutral-800/80 text-white px-3 py-1 rounded-full text-sm backdrop-blur-sm">
                {{ currentIndex + 1 }} / {{ urls.length }}
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

const router = useRouter()
const route = useRoute()
const { getImages, getImageUrl } = useImages()

const urls = ref<string[]>([])
const currentIndex = ref(0)
const loading = ref(true)

const isSwitching = ref(true)
const currentUrl = computed(() => urls.value[currentIndex.value])

const handleImageLoad = () => {
    isSwitching.value = false
}

const next = () => {
    if (currentIndex.value < urls.value.length - 1) {
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
    if (e.touches.length > 0) {
        touchStartX.value = e.touches[0].clientX
    }
}
const handleTouchEnd = (e: TouchEvent) => {
    if (e.changedTouches.length > 0) {
        const touchEndX = e.changedTouches[0].clientX
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
    const state = history.state as { urls?: string[] }
    if (state && state.urls && Array.isArray(state.urls) && state.urls.length > 0) {
        urls.value = state.urls
        loading.value = false
    } else {
        // 2. Fallback: Fetch from API using ID
        const entityId = route.query.id as string
        const entityType = route.query.type as string
        
        if (entityId && entityType) {
            try {
                const images = await getImages(entityId, entityType)
                urls.value = images.map(img => getImageUrl(img.file_path))
            } catch (e) {
                console.error("Failed to fetch preview images", e)
            }
        } else if (route.query.url) {
             // Legacy single url fallback
             urls.value = [route.query.url as string]
        }
        
        if (urls.value.length === 0) {
            isSwitching.value = false
        }
        loading.value = false
    }

})
</script>

<style scoped>
/* Ensure it covers everything including potential safe areas */
</style>
