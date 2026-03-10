<template>
  <div class="transaction-detail">
    <PageTitle
      title="交易詳情"
      :back-to="`/spaces/${route.params.id}/ledger`"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}/ledger` }]"
    >
      <template #right>
        <div class="flex gap-1">
          <NuxtLink :to="`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}/edit`" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-indigo-400 hover:bg-neutral-800 transition-colors active:scale-95 no-underline">
            <Icon icon="mdi:pencil-outline" class="text-xl" />
          </NuxtLink>
          <BaseButton @click="handleDelete" variant="danger" class="!p-0 w-10 h-10">
            <Icon icon="mdi:trash-can-outline" class="text-xl" />
          </BaseButton>
        </div>
      </template>
    </PageTitle>

    <div class="content-wrapper">
      <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

      <div v-else-if="transaction" class="content">
        <BaseCard padding="py-8 px-5" class="text-center mb-6">
          <div class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-2xl mb-4">
            <Icon :icon="getCategoryIcon(transaction.category)" class="text-xl text-indigo-500" />
            <span class="font-bold text-sm">{{ transaction.category || '未分類' }}</span>
          </div>
          
          <h2 v-if="transaction.title" class="text-lg font-bold text-white mb-1">{{ transaction.title }}</h2>
          
          <div class="text-4xl font-bold mb-2 tracking-tight">
            <span class="text-neutral-500 text-lg mr-1 font-bold">{{ transaction.currency }}</span>
            {{ formatAmount(transaction.total_amount) }}
          </div>
          <div class="text-neutral-400 text-sm font-medium">{{ formatDate(transaction.date) }}</div>

          <!-- Foreign Currency Details -->
          <div v-if="transaction.currency !== 'TWD'" class="mt-6 pt-6 border-t border-neutral-800/50 flex flex-col gap-3 text-sm">
              <div class="flex justify-between items-center text-neutral-400">
                  <span>匯率</span>
                  <span class="font-mono bg-neutral-800 px-2 py-0.5 rounded text-neutral-200">{{ transaction.exchange_rate }}</span>
              </div>
              <div class="flex justify-between items-center text-neutral-400">
                  <span>折合台幣</span>
                  <span class="text-white font-bold text-base">TWD {{ formatAmount(transaction.billing_amount) }}</span>
              </div>
              <div v-if="parseFloat(transaction.handling_fee) > 0" class="flex justify-between items-center text-neutral-400">
                  <span>手續費</span>
                  <span class="text-neutral-200">TWD {{ formatAmount(transaction.handling_fee) }}</span>
              </div>
          </div>
        </BaseCard>

        <!-- Images -->
        <div v-if="transaction.images && transaction.images.length > 0" class="mb-6">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 px-1">收據 / 照片</h2>
          <div class="grid grid-cols-3 gap-2">
            <div
              v-for="(image, idx) in transaction.images"
              :key="image.id"
              class="aspect-square rounded-xl overflow-hidden bg-neutral-800 border border-neutral-700 cursor-pointer"
              @click="router.push(`/image-preview?id=${transaction.id}&type=transaction&index=${idx}`)"
            >
              <img :src="image.file_path" class="w-full h-full object-cover" />
            </div>
          </div>
        </div>

        <div class="mb-6">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 px-1">項目明細</h2>
          <div class="flex flex-col gap-2">
            <div v-for="item in transaction.items" :key="item.id" class="flex justify-between items-center p-4 bg-neutral-900 rounded-xl border border-neutral-800/60">
              <div class="flex flex-col gap-0.5">
                <span class="font-bold text-neutral-100">{{ item.name }}</span>
                <span class="text-xs text-neutral-500 font-medium">單價 {{ transaction.currency }} {{ formatAmount(item.unit_price) }} × {{ item.quantity }}</span>
              </div>
              <div class="font-bold text-white">{{ transaction.currency }} {{ formatAmount(item.amount) }}</div>
            </div>
          </div>
        </div>

        <div v-if="transaction.note" class="mb-6">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 px-1">備註</h2>
          <div class="p-4 bg-neutral-900 rounded-xl border border-neutral-800/60 text-neutral-400 text-sm leading-relaxed">{{ transaction.note }}</div>
        </div>

        <div class="px-1 text-neutral-600 text-xs font-medium space-y-1">
          <p>建立於 {{ formatDateTime(transaction.created_at) }}</p>
          <p class="font-mono">TXN: {{ transaction.id }}</p>
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
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import BaseCard from '~/components/BaseCard.vue'

definePageMeta({
  path: '/spaces/:id/ledger/transaction/:txnId'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const transaction = ref<any>(null)
const loading = ref(true)

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
    '生活': 'mdi:home',
  }
  return icons[category] || 'mdi:receipt'
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'short',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const formatDateTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', { hour12: false })
}

const fetchData = async () => {
  try {
    transaction.value = await api.get<any>(
      `/api/spaces/${route.params.id}/transactions/${route.params.txnId}`
    )
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
    router.push(`/spaces/${route.params.id}/ledger`)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除此交易嗎？')) return

  try {
    await api.del(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`)
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger`)
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
  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()
  fetchData()
})
</script>
