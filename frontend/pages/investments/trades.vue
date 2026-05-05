<template>
  <div>
    <PageTitle title="股票交易" :show-back="false" />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="stockTrades.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <Icon icon="mdi:swap-horizontal" class="text-5xl opacity-20 mb-4" />
      <p>尚無交易紀錄</p>
    </div>

    <div v-else class="flex flex-col gap-1">
      <button
        v-for="t in stockTrades"
        :key="t.id"
        type="button"
        @click="router.push(`/investments/trades/${t.id}/edit`)"
        class="w-full flex items-center justify-between px-4 py-3 bg-neutral-900 border border-neutral-800 rounded-xl cursor-pointer transition-colors hover:bg-neutral-800/70 active:scale-[0.99]"
      >
        <div class="flex flex-col items-start gap-0.5">
          <div class="flex items-center gap-2">
            <span class="font-mono text-sm font-bold" :class="t.shares > 0 ? 'text-red-400' : 'text-emerald-400'">
              {{ t.shares > 0 ? '買' : '賣' }}
            </span>
            <span class="text-sm text-neutral-200">{{ t.symbol }}</span>
          </div>
          <span class="text-xs text-neutral-500">{{ t.trade_date }}</span>
        </div>
        <div class="flex flex-col items-end gap-0.5">
          <span class="text-sm text-neutral-200">{{ Math.abs(t.shares) }} 股 @ {{ t.price }}</span>
          <span v-if="t.note" class="text-xs text-neutral-500">{{ t.note }}</span>
        </div>
      </button>
    </div>

    <BaseFab @click="router.push('/investments/trades/add')" />

    <Transition name="slide-right">
      <NuxtPage />
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useInvestment } from '~/composables/useInvestment'
import PageTitle from '~/components/PageTitle.vue'
import BaseFab from '~/components/BaseFab.vue'

const router = useRouter()
const { stockTrades, fetchStockTrades } = useInvestment()
const loading = ref(true)

onMounted(async () => {
  try {
    await fetchStockTrades()
  } finally {
    loading.value = false
  }
})
</script>
