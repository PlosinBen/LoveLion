<template>
  <div class="add-ledger-page">
    <header class="flex justify-between items-center mb-8">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">建立新空間</h1>
      <div class="w-10"></div>
    </header>

    <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
      <div class="bg-neutral-900 p-6 rounded-3xl border border-neutral-800 flex flex-col gap-6">
        <div>
          <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-2 px-1">空間名稱</label>
          <input 
            v-model="form.name" 
            type="text" 
            placeholder="例如：生活開銷、個人私帳" 
            class="w-full bg-neutral-800 border-neutral-700 text-white py-4 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors"
            required
          />
        </div>

        <div>
          <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-2 px-1">基礎幣別</label>
          <select v-model="form.base_currency" class="w-full bg-neutral-800 border-neutral-700 text-white py-4 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors">
            <option value="TWD">TWD - 台幣</option>
            <option value="JPY">JPY - 日圓</option>
            <option value="USD">USD - 美元</option>
            <option value="EUR">EUR - 歐元</option>
          </select>
        </div>
      </div>

      <p class="text-xs text-neutral-500 px-4 leading-relaxed">
        建立空間後，您可以透過分享連結邀請其他成員加入協作，共同紀錄收支。
      </p>

      <button 
        type="submit" 
        :disabled="submitting"
        class="w-full py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-4 shadow-lg shadow-indigo-500/20"
      >
        {{ submitting ? '建立中...' : '建立空間' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

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

const handleSubmit = async () => {
  if (!form.value.name.trim()) return
  
  submitting.value = true
  try {
    // Also ensure the currencies array contains the base currency
    if (!form.value.currencies.includes(form.value.base_currency)) {
        form.value.currencies.push(form.value.base_currency)
    }

    await api.post('/api/spaces', form.value)
    router.push('/spaces')
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
