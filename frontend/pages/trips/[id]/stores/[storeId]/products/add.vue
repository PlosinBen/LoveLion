<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">新增商品</h1>
      <div class="w-10"></div>
    </header>

    <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
      <div class="flex flex-col gap-3">
        <BaseInput v-model="form.name" placeholder="商品名稱" :auto-focus="true" />
        <div class="grid grid-cols-2 gap-2">
          <BaseInput v-model.number="form.price" type="number" step="0.01" placeholder="價格" />
          <BaseSelect 
            v-model="form.currency" 
            :options="['TWD', 'JPY', 'USD']"
          />
        </div>
        <div class="grid grid-cols-2 gap-2">
          <BaseInput v-model.number="form.quantity" type="number" min="1" placeholder="數量" />
          <BaseInput v-model="form.unit" placeholder="單位 (如: 盒)" />
        </div>
        
        <div class="mt-2">
           <label class="text-sm text-neutral-400 font-medium mb-2 block">商品照片</label>
           <ImageManager 
             ref="imageManagerRef"
             entity-id="" 
             entity-type="product" 
             :max-count="5" 
             :instant-upload="false"
           />
        </div>
      </div>
      
      <button 
        @click="handleSubmit" 
        class="w-full mt-4 px-4 py-3 rounded-xl bg-indigo-500 text-white border-0 cursor-pointer hover:bg-indigo-600 transition-colors font-bold text-base disabled:opacity-50 disabled:cursor-not-allowed" 
        :disabled="!form.name.trim() || submitting"
      >
        {{ submitting ? '新增中...' : '新增商品' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import ImageManager from '~/components/ImageManager.vue'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const submitting = ref(false)
const imageManagerRef = ref<InstanceType<typeof ImageManager> | null>(null)
const form = ref({
  name: '',
  price: 0,
  currency: 'TWD',
  quantity: 1,
  unit: ''
})

const handleSubmit = async () => {
  if (!form.value.name.trim()) return
  
  submitting.value = true
  try {
    const product = await api.post<any>(`/api/trips/${route.params.id}/stores/${route.params.storeId}/products`, form.value)
    
    // Upload images
    if (imageManagerRef.value) {
        await imageManagerRef.value.commit(product.id)
    }

    router.push(`/trips/${route.params.id}/stores/${route.params.storeId}`)
  } catch (e: any) {
    alert(e.message || '新增失敗')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
  }
})
</script>
