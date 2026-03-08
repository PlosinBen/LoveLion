<template>
  <div class="flex flex-col gap-6 pb-24">
    <SpaceHeader
      :title="store?.name || '載入中...'"
      :show-back="true"
      :back-to="`/spaces/${route.params.id}/stores`"
    >
      <template #subtitle>
        <div class="flex items-center gap-1.5 text-neutral-500 text-xs font-medium mt-0.5">
          <Icon icon="mdi:package-variant-closed" class="text-indigo-500" />
          <span>{{ store?.products?.length || 0 }} 個商品</span>
          <span v-if="store?.products?.length > 0" class="text-neutral-700">|</span>
          <span v-if="store?.products?.length > 0" class="text-indigo-400">{{ getPriceRange(store.products) }}</span>
        </div>
      </template>

      <template #right>
        <div class="flex gap-2">
          <button v-if="store?.google_map_url" @click="windowOpen(store.google_map_url)" class="w-10 h-10 rounded-xl bg-neutral-900 text-green-500 flex items-center justify-center hover:bg-neutral-800 transition-colors border-0 cursor-pointer shrink-0">
            <Icon icon="mdi:google-maps" class="text-xl" />
          </button>
          <button @click="router.push(`/spaces/${route.params.id}/stores/${route.params.storeId}/edit`)" class="w-10 h-10 rounded-xl bg-neutral-900 text-white flex items-center justify-center hover:bg-neutral-800 transition-colors border-0 cursor-pointer shrink-0">
            <Icon icon="mdi:pencil-outline" class="text-xl" />
          </button>
        </div>
      </template>
    </SpaceHeader>

    <div v-if="loading" class="flex justify-center items-center text-neutral-400 py-20">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="px-2">
      <!-- Empty State -->
      <div v-if="!store?.products || store?.products?.length === 0" class="flex flex-col items-center justify-center py-20 bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed text-neutral-500">
        <Icon icon="mdi:package-variant" class="text-5xl mb-4 opacity-20" />
        <p class="text-sm">目前還沒有任何商品紀錄</p>
        <button @click="router.push(`/spaces/${route.params.id}/stores/${route.params.storeId}/products/add`)" class="mt-6 px-6 py-2 bg-indigo-500 text-white rounded-full font-bold text-sm border-0 cursor-pointer hover:bg-indigo-600 transition-colors active:scale-95">立即新增</button>
      </div>

      <!-- Product List -->
      <div v-else class="flex flex-col gap-3">
        <div v-for="product in store?.products" :key="product.id" 
             class="bg-neutral-900 rounded-2xl border border-neutral-800 transition-all duration-300">
            
            <div class="p-5 flex items-center gap-4 cursor-pointer" @click="toggleProductExpand(product.id)">
                <!-- Thumbnail/Icon -->
                <div class="w-12 h-12 rounded-xl bg-neutral-800 flex-none flex items-center justify-center text-neutral-500 border border-neutral-700">
                   <Icon icon="mdi:tag-outline" class="text-2xl" />
                </div>

                <!-- Content -->
                <div class="flex-1 min-w-0">
                   <h3 class="font-bold truncate text-neutral-100">{{ product.name }}</h3>
                   <div class="flex items-center gap-2 mt-0.5">
                       <span class="text-indigo-400 font-bold text-lg">
                           {{ product.currency }} {{ formatPrice(product.price) }}
                       </span>
                       <span v-if="product.unit" class="text-xs text-neutral-500 font-bold uppercase tracking-wider">
                           / {{ product.unit }}
                       </span>
                   </div>
                </div>

                <!-- Actions -->
                <div class="flex items-center gap-1">
                    <button @click.stop="router.push(`/spaces/${route.params.id}/stores/${route.params.storeId}/products/${product.id}/edit`)" class="w-8 h-8 flex items-center justify-center text-neutral-500 hover:text-white transition-colors border-0 cursor-pointer bg-transparent">
                        <Icon icon="mdi:pencil-outline" class="text-lg" />
                    </button>
                    <button @click.stop="deleteProduct(product.id)" class="w-8 h-8 flex items-center justify-center text-neutral-500 hover:text-red-500 transition-colors border-0 cursor-pointer bg-transparent">
                        <Icon icon="mdi:trash-can-outline" class="text-lg" />
                    </button>
                </div>
            </div>

            <!-- Expanded Details -->
            <div v-if="expandedProductId === product.id" class="px-5 pb-5 pt-0 border-t border-neutral-800/50">
                <div class="pt-4 flex flex-col gap-4">
                    <div v-if="product.note" class="text-sm text-neutral-400 bg-neutral-800/50 p-4 rounded-xl italic leading-relaxed">
                        "{{ product.note }}"
                    </div>
                    
                    <div>
                        <label class="text-xs text-neutral-500 font-bold uppercase tracking-wider mb-3 block ml-1">商品相片</label>
                        <ImageManager 
                            :ref="(el) => setProjectImageManager(el, product.id)"
                            :entity-id="product.id" 
                            entity-type="product" 
                            :max-count="5"
                            :instant-upload="false"
                            :instant-delete="false"
                        />
                        <div class="mt-4 flex justify-end">
                            <button @click.stop="saveProductImages(product.id)" class="px-5 py-2 bg-neutral-800 text-white text-xs font-bold rounded-lg border border-neutral-700 cursor-pointer hover:bg-neutral-700 transition-colors active:scale-95">
                                儲存變更
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
      </div>
    </div>

    <!-- FAB -->
    <button 
        @click="router.push(`/spaces/${route.params.id}/stores/${route.params.storeId}/products/add`)" 
        class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:plus" class="text-3xl" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import ImageManager from '~/components/ImageManager.vue'

const router = useRouter()
definePageMeta({
  layout: 'default'
})
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const { getImageUrl } = useImages()
const detailStore = useSpaceDetailStore()

const store = ref<any>(null)
const loading = ref(true)
const expandedProductId = ref<string | null>(null)
const imageManagers = ref<Record<string, any>>({})

const setProjectImageManager = (el: any, id: string) => {
    if (el) imageManagers.value[id] = el
}

const saveProductImages = async (productId: string) => {
    const manager = imageManagers.value[productId]
    if (manager) {
        try {
            await manager.commit()
            alert('相片已更新')
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
    if (min === max) return ''
    return `${formatPrice(min)} ~ ${formatPrice(max)}`
}

const toggleProductExpand = (id: string) => {
    expandedProductId.value = expandedProductId.value === id ? null : id
}

const windowOpen = (url: string) => window.open(url, '_blank')

const fetchStore = async () => {
  try {
    store.value = await api.get<any>(`/api/spaces/${route.params.id}/stores/${route.params.storeId}`)
  } catch (e) {
    console.error('Failed to fetch store:', e)
    router.push(`/spaces/${route.params.id}`)
  } finally {
    loading.value = false
  }
}

const deleteProduct = async (productId: string) => {
  if (!confirm('確定要刪除此商品紀錄嗎？')) return
  try {
    await api.delete(`/api/spaces/${route.params.id}/stores/${route.params.storeId}/products/${productId}`)
    detailStore.invalidate('stores')
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
