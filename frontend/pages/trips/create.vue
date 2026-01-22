<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">新增旅行</h1>
      <div class="w-10"></div>
    </header>

    <form @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <!-- Name -->
      <BaseInput v-model="form.name" label="旅行名稱 *" placeholder="例如：日本東京五日遊" required />

      <!-- Description -->
      <div class="flex flex-col gap-2">
        <label class="text-sm text-neutral-400">描述</label>
        <textarea v-model="form.description" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500 placeholder-neutral-400 resize-none" rows="2" placeholder="旅行簡介（選填）"></textarea>
      </div>

      <!-- Dates -->
      <div class="grid grid-cols-2 gap-3">
        <BaseInput v-model="form.start_date" type="date" label="開始日期" />
        <BaseInput v-model="form.end_date" type="date" label="結束日期" />
      </div>

      <!-- Base Currency -->
      <div class="flex flex-col gap-2">
        <label class="text-sm text-neutral-400">基準貨幣</label>
        <select v-model="form.base_currency" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500">
          <option value="TWD">TWD - 新台幣</option>
          <option value="JPY">JPY - 日圓</option>
          <option value="USD">USD - 美元</option>
          <option value="EUR">EUR - 歐元</option>
          <option value="KRW">KRW - 韓元</option>
        </select>
      </div>

      <!-- Members -->
      <div class="flex flex-col gap-2">
        <label class="text-sm text-neutral-400">旅行成員</label>
        <div class="flex flex-col gap-2">
          <div v-for="(member, idx) in form.members" :key="idx" class="flex items-center gap-2">
            <BaseInput v-model="form.members[idx]" :placeholder="`成員 ${idx + 1}`" input-class="flex-1" />
            <button type="button" @click="removeMember(idx)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-red-500/20 text-red-500 border-0 cursor-pointer hover:bg-red-500/30 transition-colors" v-if="form.members.length > 1">
              <Icon icon="mdi:close" />
            </button>
          </div>
        </div>
        <button type="button" @click="addMember" class="flex justify-center items-center gap-1.5 p-3 border-2 border-dashed border-neutral-800 rounded-xl bg-transparent text-neutral-400 cursor-pointer mt-2 hover:border-indigo-500 hover:text-indigo-500 transition-colors">
          <Icon icon="mdi:plus" /> 新增成員
        </button>
      </div>

      <!-- Submit -->
      <button type="submit" class="w-full mt-3 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '建立中...' : '建立旅行' }}
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
  description: '',
  start_date: '',
  end_date: '',
  base_currency: 'TWD',
  members: ['']
})

const addMember = () => {
  form.value.members.push('')
}

const removeMember = (idx: number) => {
  form.value.members.splice(idx, 1)
}

const handleSubmit = async () => {
  if (!form.value.name.trim()) return

  submitting.value = true
  try {
    const payload: any = {
      name: form.value.name,
      description: form.value.description,
      base_currency: form.value.base_currency,
      members: form.value.members.filter(m => m.trim())
    }
    if (form.value.start_date) {
      payload.start_date = new Date(form.value.start_date).toISOString()
    }
    if (form.value.end_date) {
      payload.end_date = new Date(form.value.end_date).toISOString()
    }

    const trip = await api.post<any>('/api/trips', payload)
    router.push(`/trips/${trip.id}`)
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
