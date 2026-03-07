<template>
  <div class="add-space-page">
    <SpaceHeader 
      title="建立新空間" 
      :show-back="true"
      class="pt-0 px-2"
    />

    <div class="pt-0">
      <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">基本資訊</label>
          <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-6 flex flex-col gap-6 shadow-sm">
            <BaseInput 
              v-model="form.name" 
              label="空間名稱" 
              placeholder="例如：日本旅行、個人記帳" 
              required
              autofocus
            />

            <BaseSelect
              v-model="form.base_currency"
              label="本位幣別"
              :options="currencyOptions"
            />
          </div>
        </div>

        <p class="text-xs text-neutral-500 px-4 leading-relaxed font-medium">
          建立空間後，您可以邀請成員共同記帳、上傳收據照片，或是進行商品比價與費用分攤。
        </p>

        <button 
          type="submit" 
          :disabled="submitting"
          class="w-full py-4 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-4 shadow-lg border-0 cursor-pointer"
        >
          {{ submitting ? '建立中...' : '建立空間' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'

definePageMeta({
  layout: 'default',
  hideGlobalNav: true
})

const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const submitting = ref(false)
const form = ref({
  name: '',
  base_currency: 'TWD',
  type: 'personal',
  categories: ['餐飲', '交通', '購物', '娛樂', '生活', '其他'],
  currencies: ['TWD']
})

const currencyOptions = [
  { label: 'TWD - 新台幣', value: 'TWD' },
  { label: 'JPY - 日圓', value: 'JPY' },
  { label: 'USD - 美金', value: 'USD' },
  { label: 'EUR - 歐元', value: 'EUR' }
]

const handleSubmit = async () => {
  if (!form.value.name.trim()) return
  
  submitting.value = true
  try {
    if (!form.value.currencies.includes(form.value.base_currency)) {
        form.value.currencies.push(form.value.base_currency)
    }

    await api.post('/api/spaces', form.value)
    router.push('/')
  } catch (e: any) {
    alert(e.message || '建立失敗')
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
