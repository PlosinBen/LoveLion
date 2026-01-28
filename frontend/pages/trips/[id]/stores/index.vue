<template>
  <div class="flex flex-col gap-4">
    <header class="flex-none pb-2 bg-neutral-900/90 backdrop-blur-md sticky top-0 z-10 -mx-4 px-4 pt-2 border-b border-white/5">
      <div class="flex justify-between items-center mb-4">
        <button @click="router.push(`/trips/${route.params.id}`)" class="flex justify-center items-center w-10 h-10 rounded-full bg-neutral-800 text-white border-0 cursor-pointer hover:bg-neutral-700 transition-colors">
          <Icon icon="mdi:arrow-left" class="text-xl" />
        </button>
        <h1 class="text-xl font-bold">比價清單</h1>
        <div class="w-10"></div> <!-- Spacer -->
      </div>

      <!-- Tabs -->
      <div class="flex p-1 bg-neutral-800 rounded-xl relative">
         <div class="absolute inset-y-1 w-[calc(50%-4px)] bg-neutral-700 rounded-lg shadow-sm transition-all duration-300 left-1"></div>
         <button class="flex-1 relative z-10 py-2 text-sm font-medium transition-colors text-center border-0 bg-transparent cursor-pointer text-white">
             商店列表
         </button>
         <button @click="router.push(`/trips/${route.params.id}/products`)" 
                 class="flex-1 relative z-10 py-2 text-sm font-medium transition-colors text-center border-0 bg-transparent cursor-pointer text-neutral-400 hover:text-neutral-200">
             商品比價
         </button>
      </div>
    </header>

    <div v-if="loading" class="flex justify-center items-center text-neutral-400 py-10">
       <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <!-- Mode: Stores -->
    <div v-else class="">
        <div v-if="stores.length === 0" class="flex flex-col items-center justify-center py-20 text-neutral-500">
           <Icon icon="mdi:store-outline" class="text-6xl mb-4 opacity-20" />
           <p>還沒有商店</p>
           <button @click="router.push(`/trips/${route.params.id}/stores/create`)" class="mt-4 px-4 py-2 bg-indigo-500 text-white rounded-lg text-sm border-0 cursor-pointer">
               新增第一間商店
           </button>
        </div>

        <div v-else class="flex flex-col gap-3">
            <div v-for="store in stores" :key="store.id"
                 class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 cursor-pointer transition-all active:scale-[0.98] hover:border-indigo-500/50"
                 @click="router.push(`/trips/${route.params.id}/stores/${store.id}`)">
                <div class="flex items-center gap-3">
                   <div class="w-12 h-12 rounded-lg bg-indigo-500/10 flex items-center justify-center text-indigo-400 overflow-hidden relative flex-none cursor-zoom-in"
                        @click.stop="openImage(storeThumbnails[store.id])">
                      <img v-if="storeThumbnails[store.id]" :src="storeThumbnails[store.id]" class="w-full h-full object-cover" />
                      <Icon v-else icon="mdi:store" class="text-2xl" />
                   </div>
                   <div class="flex-1">
                      <h3 class="font-bold text-base">{{ store.name }}</h3>
                      <p class="text-xs text-neutral-400 mt-1">{{ store.products?.length || 0 }} 個商品</p>
                   </div>
                   
                   <button v-if="store.google_map_url" @click.stop="windowOpen(store.google_map_url)" class="w-8 h-8 rounded-full hover:bg-neutral-800 text-neutral-500 hover:text-green-500 flex items-center justify-center transition-colors mr-1 cursor-pointer border-0">
                       <Icon icon="mdi:google-maps" />
                   </button>

                   <button @click.stop="deleteStore(store.id)" class="w-8 h-8 rounded-full hover:bg-red-500/20 text-neutral-500 hover:text-red-500 flex items-center justify-center transition-colors cursor-pointer border-0">
                       <Icon icon="mdi:trash-can-outline" />
                   </button>
                </div>
            </div>
        </div>
    </div>

    <!-- FAB -->
    <button 
        @click="router.push(`/trips/${route.params.id}/stores/create`)" 
        class="fixed bottom-6 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:store-plus" class="text-2xl" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const { getImagesBatch, getImageUrl } = useImages()

const stores = ref<any[]>([])
const storeThumbnails = ref<Record<string, string>>({})
const loading = ref(true)

const openImage = (url?: string) => {
    if (url) {
        router.push({
            path: '/media/view',
            query: { src: url }
        })
    }
}

const windowOpen = (url: string) => {
    window.open(url, '_blank')
}

const fetchStores = async () => {
  try {
    stores.value = await api.get<any[]>(`/api/trips/${route.params.id}/stores`)
    
    // Fetch thumbnails
    const storeIds = stores.value.map(s => s.id)
    if (storeIds.length > 0) {
        const images = await getImagesBatch(storeIds, 'store')
        images.forEach(img => {
            // Only set if not already set (to pick the first one/lowest sort order)
            if (!storeThumbnails.value[img.entity_id]) {
                storeThumbnails.value[img.entity_id] = getImageUrl(img.file_path)
            }
        })
    }
  } catch (e) {
    console.error('Failed to fetch stores:', e)
  } finally {
    loading.value = false
  }
}

const deleteStore = async (storeId: string) => {
  if (!confirm('確定要刪除此商店？')) return
  try {
    await api.del(`/api/trips/${route.params.id}/stores/${storeId}`)
    fetchStores()
  } catch (e: any) {
    alert(e.message || '刪除失敗')
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchStores()
})
</script>
