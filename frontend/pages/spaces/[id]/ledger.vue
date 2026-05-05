<template>
  <div v-if="store.space" class="space-detail-page">
    <PageTitle
      title="消費記帳"
      :settings-to="`/spaces/${store.space.id}/settings`"
      :breadcrumbs="[{ label: '我的空間', to: '/' }, { label: store.space.name, to: `/spaces/${store.space.id}/stats` }]"
    />

    <div class="flex flex-col gap-4">
      <!-- Search -->
      <div class="relative">
        <Icon icon="mdi:magnify" class="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-500 text-lg" />
        <input
          v-model="search"
          type="text"
          placeholder="搜尋交易..."
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 pl-10 pr-3 rounded-xl focus:outline-none focus:border-indigo-500"
          @input="debouncedFetch"
        >
      </div>

      <!-- Filters -->
      <div class="flex gap-2 overflow-x-auto no-scrollbar">
        <!-- Type filter -->
        <button
          v-for="opt in typeOptions"
          :key="opt.value"
          type="button"
          @click="toggleType(opt.value)"
          class="shrink-0 text-xs font-bold px-3 py-1.5 rounded-lg border cursor-pointer transition-colors"
          :class="filterType === opt.value
            ? 'bg-indigo-500/15 text-indigo-400 border-indigo-500/30'
            : 'bg-transparent text-neutral-500 border-neutral-700 hover:text-neutral-300'"
        >
          {{ opt.label }}
        </button>

        <!-- Category filter -->
        <button
          v-for="cat in categories"
          :key="cat"
          type="button"
          @click="toggleCategory(cat)"
          class="shrink-0 text-xs font-bold px-3 py-1.5 rounded-lg border cursor-pointer transition-colors"
          :class="filterCategory === cat
            ? 'bg-amber-500/15 text-amber-400 border-amber-500/30'
            : 'bg-transparent text-neutral-500 border-neutral-700 hover:text-neutral-300'"
        >
          {{ cat }}
        </button>
      </div>

      <!-- Transaction List -->
      <section class="flex flex-col gap-3">
        <div class="flex items-center justify-between px-1">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest">
            交易紀錄
            <span v-if="totalCount > 0" class="text-neutral-600 ml-1">({{ totalCount }})</span>
          </h2>
          <button
            type="button"
            @click="toggleAutoRefresh"
            class="flex items-center gap-1.5 text-xs font-bold px-2.5 py-1 rounded-lg border transition-colors"
            :class="autoRefresh
              ? 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30'
              : 'bg-transparent text-neutral-500 border-neutral-700 hover:text-neutral-300'"
          >
            <Icon :icon="autoRefresh ? 'mdi:sync' : 'mdi:sync-off'" class="text-sm" :class="{ 'animate-spin': autoRefresh && loading }" />
            自動更新
          </button>
        </div>

        <div v-if="!loading && filteredTransactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
          {{ hasFilters ? '沒有符合條件的交易' : '目前還沒有交易紀錄' }}
        </div>
        <div v-else class="flex flex-col gap-2">
          <TransactionListItem
            v-for="txn in filteredTransactions"
            :key="txn.id"
            :transaction="txn"
            :space-id="store.space.id"
          />
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-center gap-3 pt-2">
          <button
            type="button"
            :disabled="currentPage <= 1"
            @click="goToPage(currentPage - 1)"
            class="text-xs font-bold px-3 py-1.5 rounded-lg border transition-colors"
            :class="currentPage <= 1
              ? 'text-neutral-700 border-neutral-800 cursor-not-allowed'
              : 'text-neutral-400 border-neutral-700 hover:text-white hover:border-neutral-500'"
          >
            上一頁
          </button>
          <span class="text-xs text-neutral-500">{{ currentPage }} / {{ totalPages }}</span>
          <button
            type="button"
            :disabled="currentPage >= totalPages"
            @click="goToPage(currentPage + 1)"
            class="text-xs font-bold px-3 py-1.5 rounded-lg border transition-colors"
            :class="currentPage >= totalPages
              ? 'text-neutral-700 border-neutral-800 cursor-not-allowed'
              : 'text-neutral-400 border-neutral-700 hover:text-white hover:border-neutral-500'"
          >
            下一頁
          </button>
        </div>
      </section>
    </div>

    <!-- FAB for adding transaction -->
    <BaseFab @click="router.push(`/spaces/${store.space.id}/ledger/transaction/add`)" />
  </div>

  <Transition name="slide-right">
    <NuxtPage />
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import TransactionListItem from '~/components/TransactionListItem.vue'
import BaseFab from '~/components/BaseFab.vue'
import type { Transaction } from '~/types'

const PAGE_SIZE = 50
const AUTO_REFRESH_INTERVAL = 8000

const route = useRoute()
const router = useRouter()
const store = useSpaceDetailStore()

const search = ref('')
const filterType = ref('')
const filterCategory = ref('')
const filteredTransactions = ref<Transaction[]>([])
const totalCount = ref(0)
const loading = ref(false)
const currentPage = ref(1)
const autoRefresh = ref(false)
let refreshTimer: ReturnType<typeof setInterval> | null = null

const typeOptions = [
  { label: '消費', value: 'expense' },
  { label: '付款', value: 'payment' },
]

const categories = computed(() => {
  if (!store.space) return []
  return store.space.categories || []
})

const hasFilters = computed(() => search.value || filterType.value || filterCategory.value)
const totalPages = computed(() => Math.max(1, Math.ceil(totalCount.value / PAGE_SIZE)))

const toggleType = (value: string) => {
  filterType.value = filterType.value === value ? '' : value
}

const toggleCategory = (value: string) => {
  filterCategory.value = filterCategory.value === value ? '' : value
}

const buildQuery = () => {
  const params = new URLSearchParams()
  params.set('limit', String(PAGE_SIZE))
  params.set('offset', String((currentPage.value - 1) * PAGE_SIZE))
  if (search.value) params.set('search', search.value)
  if (filterType.value) params.set('type', filterType.value)
  if (filterCategory.value) params.set('category', filterCategory.value)
  return params.toString()
}

const fetchTransactions = async () => {
  if (!store.space) return
  loading.value = true

  try {
    const query = buildQuery()
    const url = `/api/spaces/${store.space.id}/transactions?${query}`

    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
    })
    const data = await response.json() as Transaction[]
    const total = parseInt(response.headers.get('X-Total-Count') || '0', 10)

    totalCount.value = total
    filteredTransactions.value = data
    store.fetched.transactions = true
  } catch (e) {
    console.error('Failed to fetch transactions:', e)
  } finally {
    loading.value = false
  }
}

const goToPage = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchTransactions()
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null
const debouncedFetch = () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    currentPage.value = 1
    fetchTransactions()
  }, 300)
}

watch([filterType, filterCategory], () => {
  currentPage.value = 1
  fetchTransactions()
})

watch(() => store.fetched.transactions, (v) => {
  if (v === false) fetchTransactions()
})

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    refreshTimer = setInterval(() => fetchTransactions(), AUTO_REFRESH_INTERVAL)
  } else {
    if (refreshTimer) clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(async () => {
  store.setSpaceId(route.params.id as string)
  try {
    await store.fetchSpace()
    await fetchTransactions()
  } catch (e) {
    router.push('/')
  }
})

onUnmounted(() => {
  if (debounceTimer) clearTimeout(debounceTimer)
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
