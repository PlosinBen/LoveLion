<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">編輯商品</h1>
      <div class="w-10"></div>
    </header>

    <div v-if="loading" class="flex justify-center items-center text-neutral-400 py-10">
      <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-5">
      <!-- Info Card -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-3">
        <BaseInput v-model="form.name" label="商品名稱" placeholder="商品名稱" />
        <div class="grid grid-cols-2 gap-2">
          <BaseInput v-model.number="form.price" type="number" step="0.01" label="價格" placeholder="價格" />
          <BaseSelect 
            v-model="form.currency" 
            :options="['TWD', 'JPY', 'USD']"
            label="幣別"
          />
        </div>
        <div class="grid grid-cols-2 gap-2">
          <BaseInput v-model.number="form.quantity" type="number" min="1" label="數量" placeholder="數量" />
          <BaseInput v-model="form.unit" label="單位" placeholder="單位 (如: 盒)" />
        </div>
      </div>

      <!-- Photo Card -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-2">
          <label class="text-sm text-neutral-400 font-medium flex items-center gap-2 mb-1">
              <Icon icon="mdi:image-multiple" /> 商品照片
          </label>
          <ImageManager 
            ref="imageManagerRef"
            :entity-id="route.params.productId as string" 
            entity-type="product" 
            :max-count="5" 
            :instant-upload="false"
            :allow-reorder="true"
            :instant-delete="false"
          />
      </div>
      
      <!-- Save Button -->
      <button 
        @click="handleSubmit" 
        class="w-full px-4 py-3 rounded-xl bg-indigo-500 text-white border-0 cursor-pointer hover:bg-indigo-600 transition-colors font-bold text-base disabled:opacity-50 disabled:cursor-not-allowed" 
        :disabled="!form.name.trim() || submitting"
      >
        {{ submitting ? '儲存中...' : '儲存變更' }}
      </button>

      <!-- Delete Button (Optional here since it's on list, but good for completeness) -->
       <!-- Skipped to avoid clutter, can add if requested -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import ImageManager from '~/components/ImageManager.vue'
import BaseInput from '~/components/BaseInput.vue'
import BaseSelect from '~/components/BaseSelect.vue'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const loading = ref(true)
const submitting = ref(false)
const imageManagerRef = ref<InstanceType<typeof ImageManager> | null>(null)
const form = ref({
  name: '',
  price: 0,
  currency: 'TWD',
  quantity: 1,
  unit: ''
})

const fetchProduct = async () => {
    try {
        const product = await api.get<any>(`/api/trips/${route.params.id}/stores/${route.params.storeId}/products/${route.params.productId}`)
        form.value.name = product.name
        form.value.price = product.price
        form.value.currency = product.currency
        form.value.quantity = product.quantity
        form.value.unit = product.unit || ''
    } catch (e) {
        console.error('Failed to fetch product', e)
        router.back()
    } finally {
        loading.value = false
    }
}

const handleSubmit = async () => {
  if (!form.value.name.trim()) return
  
  submitting.value = true
  try {
    await api.put<any>(`/api/trips/${route.params.id}/stores/${route.params.storeId}/products/${route.params.productId}`, form.value)
    
    // Commit images
    if (imageManagerRef.value) {
        await imageManagerRef.value.commit()
    }

    router.back()
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchProduct()
})
</script>
