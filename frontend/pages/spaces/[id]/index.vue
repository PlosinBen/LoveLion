<template>
  <div v-if="space" class="space-detail-page">
    <SpaceHeader 
      :title="space.name" 
      :show-back="true"
    />

    <div class="flex flex-col gap-6">
      <section v-if="transactions.length > 0">
        <SpaceStats 
          :transactions="transactions" 
          :members="space.members || []" 
          :base-currency="space.base_currency || 'TWD'"
        />
      </section>

      <section class="flex flex-col gap-3">
        <div class="flex items-center justify-between px-1">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest">最近交易</h2>
          <button @click="router.push(`/spaces/${space.id}/transaction/add`)" class="text-xs font-bold text-indigo-400 border-0 bg-transparent cursor-pointer hover:text-indigo-300">
             查看全部
          </button>
        </div>
        
        <div v-if="transactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
          目前還沒有交易紀錄
        </div>
        <div v-else class="flex flex-col gap-2">
          <TransactionListItem 
            v-for="txn in transactions.slice(0, 5)" 
            :key="txn.id" 
            :transaction="txn" 
            :space-id="space.id"
          />
        </div>
      </section>
    </div>

    <!-- FAB for adding transaction -->
    <button 
      @click="router.push(`/spaces/${space.id}/transaction/add`)"
      class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg rounded-full flex items-center justify-center text-white transition-transform active:scale-95 z-20 cursor-pointer border-0"
    >
      <Icon icon="mdi:plus" class="text-3xl" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'
import SpaceStats from '~/components/SpaceStats.vue'
import TransactionListItem from '~/components/TransactionListItem.vue'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const space = ref<any>(null)
const transactions = ref<any[]>([])

const fetchData = async () => {
  try {
    const spaceId = route.params.id
    // Fetch space with members
    space.value = await api.get<any>(`/api/spaces/${spaceId}`)
    // Fetch transactions
    const txnData = await api.get<any>(`/api/spaces/${spaceId}/transactions`)
    transactions.value = txnData || []
  } catch (e) {
    console.error('Failed to fetch space data:', e)
    router.push('/')
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchData()
})
</script>
