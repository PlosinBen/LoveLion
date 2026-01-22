<template>
  <div class="add-transaction-page">
    <header class="flex justify-between items-center mb-6">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">{{ isEdit ? '編輯交易' : '新增記帳' }}</h1>
      <div style="width: 40px;"></div>
    </header>

    <form @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <!-- Date -->
      <BaseInput v-model="form.date" type="date" label="日期" />

      <!-- Category -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">類別</label>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="cat in categories"
            :key="cat"
            type="button"
            class="px-4 py-2.5 rounded-3xl border border-neutral-800 bg-neutral-900 text-white cursor-pointer transition-all duration-200 hover:border-indigo-500"
            :class="{ 'bg-indigo-500 border-indigo-500': form.category === cat }"
            @click="form.category = cat"
          >
            {{ cat }}
          </button>
        </div>
      </div>

      <!-- Items -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">項目明細</label>
        <div class="flex flex-col gap-3">
          <div v-for="(item, index) in form.items" :key="index" class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800">
            <div class="flex flex-col gap-2">
              <BaseInput
                v-model="item.name"
                placeholder="項目名稱"
                input-class="font-medium"
              />
              <div class="flex items-center gap-2">
                <BaseInput
                  v-model.number="item.unit_price"
                  type="number"
                  placeholder="單價"
                  input-class="flex-[2]"
                />
                <span class="text-neutral-400">×</span>
                <BaseInput
                  v-model.number="item.quantity"
                  type="number"
                  placeholder="數量"
                  min="1"
                  input-class="flex-1 w-[60px]"
                />
                <button type="button" @click="removeItem(index)" class="flex justify-center items-center w-9 h-9 rounded-lg bg-red-500 text-white border-0 cursor-pointer hover:bg-red-600 transition-colors" v-if="form.items.length > 1">
                  <Icon icon="mdi:close" />
                </button>
              </div>
            </div>
            <div class="text-right text-neutral-400 text-sm mt-3">
              {{ currency }} {{ (item.unit_price * item.quantity).toLocaleString() }}
            </div>
          </div>
        </div>
        <button type="button" @click="addItem" class="flex justify-center items-center gap-1.5 p-3 border-2 border-dashed border-neutral-800 rounded-xl bg-transparent text-neutral-400 cursor-pointer mt-3 hover:border-indigo-500 hover:text-indigo-500 transition-colors">
          <Icon icon="mdi:plus" /> 新增項目
        </button>
      </div>

      <!-- Total -->
      <div class="flex justify-between items-center text-lg bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <span>總計</span>
        <span class="text-2xl font-bold text-indigo-500">{{ currency }} {{ totalAmount.toLocaleString() }}</span>
      </div>

      <!-- Note -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">備註</label>
        <textarea v-model="form.note" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500 placeholder-neutral-400 text-base resize-none" rows="2" placeholder="選填"></textarea>
      </div>

      <!-- Submit -->
      <button type="submit" class="w-full mt-3 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '儲存中...' : '儲存' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const isEdit = computed(() => !!route.params.id)
const currency = 'TWD'
const categories = ['餐飲', '交通', '購物', '娛樂', '生活', '其他']
const ledgerId = ref('')
const submitting = ref(false)

const form = ref({
  date: new Date().toISOString().split('T')[0],
  category: '',
  note: '',
  items: [{ name: '', unit_price: 0, quantity: 1 }]
})

const totalAmount = computed(() => {
  return form.value.items.reduce((sum, item) => {
    return sum + (item.unit_price * item.quantity)
  }, 0)
})

const addItem = () => {
  form.value.items.push({ name: '', unit_price: 0, quantity: 1 })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const handleSubmit = async () => {
  if (!form.value.items.some(item => item.name && item.unit_price > 0)) {
    alert('請至少填寫一個項目')
    return
  }

  submitting.value = true

  try {
    const payload = {
      date: new Date(form.value.date).toISOString(),
      category: form.value.category,
      note: form.value.note,
      items: form.value.items.filter(item => item.name && item.unit_price > 0)
    }

    await api.post(`/api/ledgers/${ledgerId.value}/transactions`, payload)
    router.push('/ledger')
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    submitting.value = false
  }
}

const fetchLedger = async () => {
  try {
    const ledgers = await api.get<any[]>('/api/ledgers')
    if (ledgers.length === 0) {
      const newLedger = await api.post<any>('/api/ledgers', {
        name: '我的帳本',
        type: 'personal',
        currencies: ['TWD'],
        categories
      })
      ledgerId.value = newLedger.id
    } else {
      ledgerId.value = ledgers[0].id
    }
  } catch (e) {
    console.error('Failed to fetch ledger:', e)
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchLedger()
})
</script>
