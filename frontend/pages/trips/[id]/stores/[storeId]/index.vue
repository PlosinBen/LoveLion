<template>
  <div class="flex flex-col gap-6 min-h-screen bg-neutral-900 text-neutral-50 pb-4">
    
    <!-- Hero Header -->
    <ImmersiveHeader
        :image="coverImage"
        fallback-icon="mdi:store"
        class="rounded-2xl"
    >
        <template #top-left>
             <button @click="router.back()" class="w-10 h-10 rounded-full bg-black/30 backdrop-blur-md text-white flex items-center justify-center hover:bg-black/50 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:arrow-left" class="text-xl" />
             </button>
        </template>

        <template #top-right>
             <!-- Edit Button -->
             <button @click="router.push(`/trips/${route.params.id}/stores/${route.params.storeId}/edit`)" class="w-10 h-10 rounded-full bg-black/30 backdrop-blur-md text-white flex items-center justify-center hover:bg-black/50 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:pencil-outline" class="text-xl" />
             </button>
        </template>

        <template #bottom>
            <div class="pb-2">
                <div class="flex items-center gap-2 mb-1">
                    <h1 class="text-2xl font-bold text-white shadow-sm">{{ store?.name || '載入中...' }}</h1>
                    <button v-if="store?.google_map_url" @click="windowOpen(store.google_map_url)" class="w-8 h-8 rounded-full bg-neutral-800/80 backdrop-blur-md text-neutral-400 hover:text-green-500 hover:bg-neutral-800 flex items-center justify-center transition-all border-0 cursor-pointer">
                        <Icon icon="mdi:google-maps" />
                    </button>
                </div>
                <div class="flex items-center gap-2 text-sm text-neutral-300">
                    <span class="flex items-center gap-1">
                        <Icon icon="mdi:package-variant-closed" /> {{ store?.products?.length || 0 }} 個商品
                    </span>
                    <span v-if="store?.products?.length > 0" class="text-neutral-500">•</span>
                    <span v-if="store?.products?.length > 0" class="text-indigo-400">
                         {{ getPriceRange(store.products) }}
                    </span>
                </div>
            </div>
        </template>
    </ImmersiveHeader>


    <div v-if="loading" class="flex justify-center items-center text-neutral-400 py-10">
      <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <!-- Store Content (Product List) -->
    <div v-else class="">
      
      <!-- Empty State -->
      <div v-if="!store?.products || store?.products?.length === 0" class="flex flex-col items-center justify-center py-20 text-neutral-500">
        <Icon icon="mdi:package-variant" class="text-6xl mb-4 opacity-20" />
        <p>還沒有商品</p>
      </div>

      <!-- Product List (Compact) -->
      <div v-else class="flex flex-col gap-2">
        <div v-for="product in store?.products" :key="product.id" 
             class="bg-neutral-900 rounded-xl border border-neutral-800 overflow-hidden transition-all duration-300"
             :class="{'ring-2 ring-indigo-500/50': expandedProductId === product.id}">
            
            <!-- Main Row -->
            <div class="p-3 flex items-center gap-3 cursor-pointer" @click="toggleProductExpand(product.id)">
                <!-- Thumbnail -->
                <div class="w-12 h-12 rounded-lg bg-neutral-800 flex-none flex items-center justify-center overflow-hidden">
                   <Icon icon="mdi:image-outline" class="text-neutral-600 text-xl" />
                </div>

                <!-- Content -->
                <div class="flex-1 min-w-0">
                   <h3 class="text-sm font-medium truncate text-white">{{ product.name }}</h3>
                   <div class="flex items-center gap-2 mt-0.5">
                       <span class="text-indigo-400 font-bold text-base">
                           {{ product.currency }} {{ formatPrice(product.price) }}
                       </span>
                       <span v-if="product.unit" class="text-xs text-neutral-500 px-1.5 py-0.5 bg-neutral-800 rounded">
                           / {{ product.unit }}
                       </span>
                   </div>
                </div>

                <!-- Actions -->
                <div class="flex items-center">
                    <button @click.stop="router.push(`/trips/${route.params.id}/stores/${route.params.storeId}/products/${product.id}/edit`)" class="w-8 h-8 flex items-center justify-center text-neutral-500 hover:text-indigo-400 hover:bg-neutral-800 rounded-full transition-colors border-0 cursor-pointer bg-transparent">
                        <Icon icon="mdi:pencil-outline" class="text-lg" />
                    </button>
                    <button @click.stop="deleteProduct(product.id)" class="w-8 h-8 flex items-center justify-center text-neutral-500 hover:text-red-500 hover:bg-neutral-800 rounded-full transition-colors border-0 cursor-pointer bg-transparent">
                        <Icon icon="mdi:trash-can-outline" class="text-lg" />
                    </button>
                </div>
            </div>

            <!-- Expanded Details -->
            <div v-if="expandedProductId === product.id" class="px-3 pb-3 pt-0 border-t border-neutral-800/50 mt-1">
                <div class="pt-3">
                    <label class="text-xs text-neutral-500 mb-2 block">商品照片與備註</label>
                     <ImageManager 
                        :ref="(el) => setProjectImageManager(el, product.id)"
                        :entity-id="product.id" 
                        entity-type="product" 
                        :max-count="5"
                        :instant-upload="false"
                        :instant-delete="false"
                    />
                    <div v-if="product.note" class="mt-3 text-sm text-neutral-400 bg-neutral-900 p-2 rounded">
                        {{ product.note }}
                    </div>
                    
                    <div class="mt-3 flex justify-end">
                        <button 
                            @click.stop="saveProductImages(product.id)"
                            class="px-4 py-1.5 bg-indigo-500 text-white text-xs font-bold rounded-lg border-0 cursor-pointer hover:bg-indigo-600 transition-colors"
                        >
                            儲存照片
                        </button>
                    </div>
                </div>
            </div>
        </div>
      </div>
    </div>

    <!-- FAB -->
    <button 
        @click="router.push(`/trips/${route.params.id}/stores/${route.params.storeId}/products/add`)" 
        class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:plus" class="text-3xl" />
    </button>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUpdate } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'
import ImageManager from '~/components/ImageManager.vue'
import ImmersiveHeader from '~/components/ImmersiveHeader.vue'

const router = useRouter()
definePageMeta({
  layout: 'main'
})
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const { getImages, getImageUrl } = useImages()

const store = ref<any>(null)
const loading = ref(true)
const expandedProductId = ref<string | null>(null)
const coverImage = ref<string | null>(null)
const imageManagers = ref<Record<string, any>>({})

const setProjectImageManager = (el: any, id: string) => {
    if (el) {
        imageManagers.value[id] = el
    }
}

const saveProductImages = async (productId: string) => {
    const manager = imageManagers.value[productId]
    if (manager) {
        try {
            await manager.commit()
            // Optionally refresh store or show success
            // alert('照片已儲存') 
        } catch (e) {
            console.error('Failed to save images', e)
            alert('儲存失敗')
        }
    }
}


const formatPrice = (price: number | string) => {
  const num = typeof price === 'string' ? parseFloat(price) : price
  return num.toLocaleString('zh-TW', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

const getPriceRange = (products: any[]) => {
    if (!products || products.length === 0) return ''
    const prices = products
        .map(p => parseFloat(p?.price || '0'))
        .filter(p => !isNaN(p))
    
    if (prices.length === 0) return ''
    
    const min = Math.min(...prices)
    const max = Math.max(...prices)
    if (min === max) return '' // Don't show range if one price
    return `${formatPrice(min)} ~ ${formatPrice(max)}`
}

const toggleProductExpand = (id: string) => {
    if (expandedProductId.value === id) {
        expandedProductId.value = null
    } else {
        expandedProductId.value = id
    }
}

const fetchStoreImages = async () => {
    try {
        const images = await getImages(route.params.storeId as string, 'store')
        if (images.length > 0) {
            coverImage.value = getImageUrl(images[0].file_path)
        } else {
            coverImage.value = null
        }
    } catch (e) {
        console.error('Failed to fetch store images', e)
    }
}

const windowOpen = (url: string) => {
    window.open(url, '_blank')
}

const fetchStore = async () => {
  try {
    store.value = await api.get<any>(`/api/trips/${route.params.id}/stores/${route.params.storeId}`)
    await fetchStoreImages()
  } catch (e) {
    console.error('Failed to fetch store:', e)
    router.push(`/trips/${route.params.id}/stores`)
  } finally {
    loading.value = false
  }
}

const deleteProduct = async (productId: string) => {
  if (!confirm('確定要刪除此商品？')) return
  try {
    await api.del(`/api/trips/${route.params.id}/stores/${route.params.storeId}/products/${productId}`)
    fetchStore()
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
  fetchStore()
})
</script>
