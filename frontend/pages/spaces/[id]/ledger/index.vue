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

          <!-- Load more -->
          <div v-if="hasMore" ref="loadMoreRef" class="flex justify-center py-4">
            <span v-if="loadingMore" class="text-xs text-neutral-500">載入中...</span>
          </div>
        </div>
      </section>
    </div>

    <!-- FAB for adding transaction -->
    <BaseFab @click="router.push(`/spaces/${store.space.id}/ledger/transaction/add`)" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import TransactionListItem from '~/components/TransactionListItem.vue'
import BaseFab from '~/components/BaseFab.vue'
import type { Transaction } from '~/types'

const PAGE_SIZE = 20

const route = useRoute()
const router = useRouter()
const api = useApi()
const store = useSpaceDetailStore()

const search = ref('')
const filterType = ref('')
const filterCategory = ref('')
const filteredTransactions = ref<Transaction[]>([])
const totalCount = ref(0)
const loading = ref(false)
const loadingMore = ref(false)
const offset = ref(0)
const loadMoreRef = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

const typeOptions = [
  { label: '消費', value: 'expense' },
  { label: '付款', value: 'payment' },
]

const categories = computed(() => {
  if (!store.space) return []
  return store.space.categories || []
})

const hasFilters = computed(() => search.value || filterType.value || filterCategory.value)
const hasMore = computed(() => filteredTransactions.value.length < totalCount.value)

const toggleType = (value: string) => {
  filterType.value = filterType.value === value ? '' : value
}

const toggleCategory = (value: string) => {
  filterCategory.value = filterCategory.value === value ? '' : value
}

const buildQuery = (currentOffset: number) => {
  const params = new URLSearchParams()
  params.set('limit', String(PAGE_SIZE))
  params.set('offset', String(currentOffset))
  if (search.value) params.set('search', search.value)
  if (filterType.value) params.set('type', filterType.value)
  if (filterCategory.value) params.set('category', filterCategory.value)
  return params.toString()
}

const fetchTransactions = async (append = false) => {
  if (!store.space) return
  if (append) {
    loadingMore.value = true
  } else {
    loading.value = true
    offset.value = 0
    filteredTransactions.value = []
  }

  try {
    const query = buildQuery(offset.value)
    const url = `/api/spaces/${store.space.id}/transactions?${query}`

    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
    })
    const data = await response.json() as Transaction[]
    const total = parseInt(response.headers.get('X-Total-Count') || '0', 10)

    totalCount.value = total
    if (append) {
      filteredTransactions.value = [...filteredTransactions.value, ...data]
    } else {
      filteredTransactions.value = data
    }
  } catch (e) {
    console.error('Failed to fetch transactions:', e)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const loadMore = () => {
  if (loadingMore.value || !hasMore.value) return
  offset.value += PAGE_SIZE
  fetchTransactions(true)
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null
const debouncedFetch = () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchTransactions(), 300)
}

// Watch filter changes (not search, which has its own debounce)
watch([filterType, filterCategory], () => {
  fetchTransactions()
})

// Intersection observer for infinite scroll
const setupObserver = () => {
  observer = new IntersectionObserver((entries) => {
    if (entries[0]?.isIntersecting) {
      loadMore()
    }
  }, { rootMargin: '200px' })
}

watch(loadMoreRef, (el) => {
  if (el && observer) {
    observer.observe(el)
  }
})

onMounted(async () => {
  store.setSpaceId(route.params.id as string)
  setupObserver()
  try {
    await store.fetchSpace()
    await fetchTransactions()
  } catch (e) {
    router.push('/')
  }
})

onUnmounted(() => {
  observer?.disconnect()
  if (debounceTimer) clearTimeout(debounceTimer)
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
