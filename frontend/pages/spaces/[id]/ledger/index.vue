<template>
  <div v-if="store.space" class="space-detail-page">
    <PageTitle
      title="消費記帳"
      :settings-to="`/spaces/${store.space.id}/settings`"
      :breadcrumbs="[{ label: '我的空間', to: '/' }, { label: store.space.name, to: `/spaces/${store.space.id}/stats` }]"
    />

    <div class="flex flex-col gap-6">
      <section class="flex flex-col gap-3">
        <div class="flex items-center justify-between px-1">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest">交易紀錄</h2>
        </div>

        <div v-if="store.transactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
          目前還沒有交易紀錄
        </div>
        <div v-else class="flex flex-col gap-2">
          <TransactionListItem
            v-for="txn in store.transactions"
            :key="txn.id"
            :transaction="txn"
            :space-id="store.space.id"
          />
        </div>
      </section>
    </div>

    <!-- FAB for adding transaction -->
    <BaseFab @click="router.push(`/spaces/${store.space.id}/ledger/transaction/add`)" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import TransactionListItem from '~/components/TransactionListItem.vue'
import BaseFab from '~/components/BaseFab.vue'

const route = useRoute()
const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()
const store = useSpaceDetailStore()

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
