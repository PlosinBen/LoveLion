<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.push(`/trips/${route.params.id}/stores`)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold truncate max-w-48">{{ store?.name || '載入中...' }}</h1>
      <button @click="showAddProduct = true" class="flex justify-center items-center w-10 h-10 rounded-xl bg-indigo-500 text-white border-0 cursor-pointer hover:bg-indigo-600 transition-colors">
        <Icon icon="mdi:plus" class="text-2xl" />
      </button>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <div v-else-if="store?.products?.length === 0" class="text-center py-16 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
      <Icon icon="mdi:package-variant" class="text-6xl text-neutral-500 mb-4" />
      <h2 class="text-lg font-semibold mb-2">還沒有商品</h2>
      <p class="text-neutral-400 mb-5">新增商品來記錄價格</p>
      <button @click="showAddProduct = true" class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm border-0 cursor-pointer">
        新增商品
      </button>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div v-for="product in store?.products" :key="product.id" class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800">
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <h3 class="text-base font-semibold">{{ product.name }}</h3>
            <p class="text-sm text-neutral-400 mt-1">
              {{ product.currency }} {{ formatPrice(product.price) }}
              <span v-if="product.quantity > 1" class="ml-1">× {{ product.quantity }} {{ product.unit }}</span>
            </p>
          </div>
          <button @click="deleteProduct(product.id)" class="flex justify-center items-center w-8 h-8 rounded-lg bg-red-500/20 text-red-500 border-0 cursor-pointer hover:bg-red-500/30 transition-colors">
            <Icon icon="mdi:delete" />
          </button>
        </div>
      </div>
    </div>

    <!-- Add Product Modal -->
    <div v-if="showAddProduct" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4" @click.self="showAddProduct = false">
      <div class="bg-neutral-900 rounded-2xl p-5 w-full max-w-sm border border-neutral-800">
        <h3 class="text-lg font-bold mb-4">新增商品</h3>
        <div class="flex flex-col gap-3">
          <BaseInput v-model="newProduct.name" placeholder="商品名稱" />
          <div class="grid grid-cols-2 gap-2">
            <BaseInput v-model.number="newProduct.price" type="number" step="0.01" placeholder="價格" />
            <select v-model="newProduct.currency" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500">
              <option value="TWD">TWD</option>
              <option value="JPY">JPY</option>
              <option value="USD">USD</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <BaseInput v-model.number="newProduct.quantity" type="number" min="1" placeholder="數量" />
            <BaseInput v-model="newProduct.unit" placeholder="單位 (如: 盒)" />
          </div>
        </div>
        <div class="flex gap-2 mt-4">
          <button @click="showAddProduct = false" class="flex-1 px-4 py-3 rounded-xl bg-neutral-800 text-white border-0 cursor-pointer hover:bg-neutral-700 transition-colors">取消</button>
          <button @click="addProduct" class="flex-1 px-4 py-3 rounded-xl bg-indigo-500 text-white border-0 cursor-pointer hover:bg-indigo-600 transition-colors" :disabled="!newProduct.name.trim()">新增</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const store = ref<any>(null)
const loading = ref(true)
const showAddProduct = ref(false)
const newProduct = ref({
  name: '',
  price: 0,
  currency: 'TWD',
  quantity: 1,
  unit: ''
})

const formatPrice = (price: number | string) => {
  const num = typeof price === 'string' ? parseFloat(price) : price
  return num.toLocaleString('zh-TW', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

const fetchStore = async () => {
  try {
    store.value = await api.get<any>(`/api/trips/${route.params.id}/stores/${route.params.storeId}`)
  } catch (e) {
    console.error('Failed to fetch store:', e)
    router.push(`/trips/${route.params.id}/stores`)
  } finally {
    loading.value = false
  }
}

const addProduct = async () => {
  if (!newProduct.value.name.trim()) return
  try {
    await api.post(`/api/trips/${route.params.id}/stores/${route.params.storeId}/products`, newProduct.value)
    newProduct.value = { name: '', price: 0, currency: 'TWD', quantity: 1, unit: '' }
    showAddProduct.value = false
    fetchStore()
  } catch (e: any) {
    alert(e.message || '新增失敗')
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
