<template>
  <div v-if="store.space" class="space-detail-page">
    <PageTitle
      :title="store.space.name"
      :show-back="true"
      :settings-to="`/spaces/${store.space.id}/settings`"
    />

    <div class="flex flex-col gap-6">
      <section class="flex flex-col gap-3">
        <div class="flex items-center justify-between px-1">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest">最近交易</h2>
          <button @click="showAll = !showAll" class="text-xs font-bold text-indigo-400 border-0 bg-transparent cursor-pointer hover:text-indigo-300">
             {{ showAll ? '收起' : '查看全部' }}
          </button>
        </div>

        <div v-if="store.transactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
          目前還沒有交易紀錄
        </div>
        <div v-else class="flex flex-col gap-2">
          <TransactionListItem
            v-for="txn in displayedTransactions"
            :key="txn.id"
            :transaction="txn"
            :space-id="store.space.id"
          />
        </div>
      </section>
    </div>

    <!-- FAB for adding transaction -->
    <button
      @click="router.push(`/spaces/${store.space.id}/transaction/add`)"
      class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg rounded-full flex items-center justify-center text-white transition-transform active:scale-95 z-20 cursor-pointer border-0"
    >
      <Icon icon="mdi:plus" class="text-3xl" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import TransactionListItem from '~/components/TransactionListItem.vue'

const route = useRoute()
const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()
const store = useSpaceDetailStore()

const showAll = ref(false)

const displayedTransactions = computed(() => {
  return showAll.value ? store.transactions : store.transactions.slice(0, 5)
})

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  store.setSpaceId(route.params.id as string)
  try {
    await Promise.all([store.fetchSpace(), store.fetchTransactions()])
  } catch (e) {
    router.push('/')
  }
})
</script>
