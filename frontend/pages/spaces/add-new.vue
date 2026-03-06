<template>
  <div class="add-ledger-page">
    <header class="flex justify-between items-center mb-8">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">撱箇??啁征??/h1>
      <div class="w-10"></div>
    </header>

    <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
      <div class="bg-neutral-900 p-6 rounded-3xl border border-neutral-800 flex flex-col gap-6">
        <div>
          <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-2 px-1">蝛粹??迂</label>
          <input 
            v-model="form.name" 
            type="text" 
            placeholder="靘?嚗?瘣駁??瑯犖蝘董" 
            class="w-full bg-neutral-800 border-neutral-700 text-white py-4 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors"
            required
          />
        </div>

        <div>
          <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-2 px-1">?箇?撟?</label>
          <select v-model="form.base_currency" class="w-full bg-neutral-800 border-neutral-700 text-white py-4 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors">
            <option value="TWD">TWD - ?啣馳</option>
            <option value="JPY">JPY - ?亙?</option>
            <option value="USD">USD - 蝢?</option>
            <option value="EUR">EUR - 甇?</option>
          </select>
        </div>
      </div>

      <p class="text-xs text-neutral-500 px-4 leading-relaxed">
        撱箇?蝛粹?敺??典隞仿??澈????隢隞??∪??亙?雿??勗?蝝??胯?      </p>

      <button 
        type="submit" 
        :disabled="submitting"
        class="w-full py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-4 shadow-lg shadow-indigo-500/20"
      >
        {{ submitting ? '撱箇?銝?..' : '撱箇?蝛粹?' }}
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
  categories: ['擗ㄡ', '鈭日?, '鞈潛', '憡?', '?暑', '?嗡?'],
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
    alert(e.message || '撱箇?憭望?')
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
