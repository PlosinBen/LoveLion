<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.push(`/trips/${route.params.id}`)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">價格比較</h1>
      <div class="w-10"></div>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <div v-else-if="products.length === 0" class="text-center py-16 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
      <Icon icon="mdi:chart-bar" class="text-6xl text-neutral-500 mb-4" />
      <h2 class="text-lg font-semibold mb-2">還沒有商品</h2>
      <p class="text-neutral-400 mb-5">請先新增商店和商品來比較價格</p>
      <NuxtLink :to="`/trips/${route.params.id}/stores`" class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline">
        前往商店
      </NuxtLink>
    </div>

    <template v-else>
      <!-- Comparison Cards -->
      <div v-for="productName in uniqueProducts" :key="productName" class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800">
        <h3 class="text-base font-semibold mb-3">{{ productName }}</h3>
        <div class="flex flex-col gap-2">
          <div
            v-for="item in getProductPrices(productName)"
            :key="item.store_id"
            class="flex items-center justify-between p-3 rounded-xl"
            :class="item.isCheapest ? 'bg-green-500/20 border border-green-500/50' : 'bg-neutral-800'"
          >
            <div class="flex items-center gap-2">
              <Icon v-if="item.isCheapest" icon="mdi:check-circle" class="text-green-500" />
              <!-- Store Thumbnail -->
              <div class="w-6 h-6 rounded bg-neutral-700 overflow-hidden flex-shrink-0" v-if="storeImages[item.store_id]">
                  <img :src="getImageUrl(storeImages[item.store_id] || '')" class="w-full h-full object-cover" />
              </div>
              <span class="text-sm">{{ item.store_name }}</span>
            </div>
            <div class="text-right">
              <span class="font-semibold" :class="item.isCheapest ? 'text-green-500' : ''">
                {{ item.currency }} {{ formatPrice(item.price) }}
              </span>
              <span v-if="item.quantity > 1" class="text-xs text-neutral-400 ml-1">
                × {{ item.quantity }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const products = ref<any[]>([])
const storeImages = ref<Record<string, string>>({})
const { getImagesBatch, getImageUrl } = useImages()
const loading = ref(true)

const uniqueProducts = computed(() => {
  const names = new Set(products.value.map(p => p.name))
  return Array.from(names).sort()
})

const getProductPrices = (productName: string) => {
  const items = products.value.filter(p => p.name === productName)
  const minPrice = Math.min(...items.map(i => parseFloat(i.price)))
  return items.map(item => ({
    ...item,
    isCheapest: parseFloat(item.price) === minPrice
  })).sort((a, b) => parseFloat(a.price) - parseFloat(b.price))
}

const formatPrice = (price: number | string) => {
  const num = typeof price === 'string' ? parseFloat(price) : price
  return num.toLocaleString('zh-TW', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

const fetchProducts = async () => {
  try {
    products.value = await api.get<any[]>(`/api/trips/${route.params.id}/products`)
    
    // Fetch Store Images
    if (products.value.length > 0) {
        const storeIds = [...new Set(products.value.map(p => p.store_id))]
        const images = await getImagesBatch(storeIds as string[], 'store')
        const map: Record<string, string> = {}
        images.forEach(img => {
            if (!map[img.entity_id]) map[img.entity_id] = img.file_path
        })
        storeImages.value = map
    }
  } catch (e) {
    console.error('Failed to fetch products:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchProducts()
})
</script>
