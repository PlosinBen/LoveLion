<template>
  <OverlayPage>
    <PageTitle
      :title="`結算 ${ym}`"
      :show-back="true"
      :breadcrumbs="[{ label: '結算列表', to: '/investments/settlements' }]"
    />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <template v-else-if="detail">
      <div class="flex items-center gap-2 mb-4">
        <span
          class="text-xs font-bold px-2 py-0.5 rounded"
          :class="detail.status === 'completed' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-amber-500/20 text-amber-400'"
        >
          {{ detail.status === 'completed' ? '已完成' : '草稿' }}
        </span>
      </div>

      <!-- Futures -->
      <section class="mb-6">
        <div class="flex items-center justify-between mb-2">
          <h3 class="text-sm font-bold text-neutral-300">期貨</h3>
          <button
            v-if="detail.status === 'draft'"
            type="button"
            @click="router.push(`/investments/settlements/${ym}/futures`)"
            class="text-xs text-indigo-400 font-bold bg-transparent border-0 cursor-pointer"
          >
            編輯
          </button>
        </div>
        <div v-if="detail.futures_statement" class="grid grid-cols-2 gap-x-4 gap-y-1 text-sm bg-neutral-900 rounded-xl p-3 border border-neutral-800">
          <div class="text-neutral-500">期末權益</div><div class="text-right text-neutral-200">{{ fmt(detail.futures_statement.ending_equity) }}</div>
          <div class="text-neutral-500">浮動損益</div><div class="text-right text-neutral-200">{{ fmt(detail.futures_statement.floating_profit_loss) }}</div>
          <div class="text-neutral-500">沖銷損益</div><div class="text-right text-neutral-200">{{ fmt(detail.futures_statement.realized_profit_loss) }}</div>
          <div class="text-neutral-500">入金</div><div class="text-right text-neutral-200">{{ fmt(detail.futures_statement.deposit) }}</div>
          <div class="text-neutral-500">出金</div><div class="text-right text-neutral-200">{{ fmt(detail.futures_statement.withdrawal) }}</div>
          <div class="text-neutral-400 font-bold border-t border-neutral-700 pt-1 mt-1">損益</div>
          <div class="text-right font-bold border-t border-neutral-700 pt-1 mt-1" :class="detail.futures_statement.profit_loss >= 0 ? 'text-emerald-400' : 'text-red-400'">
            {{ formatPL(detail.futures_statement.profit_loss) }}
          </div>
        </div>
        <div v-else class="text-sm text-neutral-600 italic px-1">尚未填寫</div>
      </section>

      <!-- Stocks -->
      <section class="mb-6">
        <div class="flex items-center justify-between mb-2">
          <h3 class="text-sm font-bold text-neutral-300">股票</h3>
          <button
            v-if="detail.status === 'draft'"
            type="button"
            @click="router.push(`/investments/settlements/${ym}/stocks`)"
            class="text-xs text-indigo-400 font-bold bg-transparent border-0 cursor-pointer"
          >
            編輯
          </button>
        </div>
        <div v-if="detail.stock_statement" class="grid grid-cols-2 gap-x-4 gap-y-1 text-sm bg-neutral-900 rounded-xl p-3 border border-neutral-800">
          <div class="text-neutral-500">帳戶餘額</div><div class="text-right text-neutral-200">{{ fmt(detail.stock_statement.account_balance) }}</div>
          <div class="text-neutral-500">庫存現值</div><div class="text-right text-neutral-200">{{ fmt(detail.stock_statement.market_value) }}</div>
          <div class="text-neutral-500">入金</div><div class="text-right text-neutral-200">{{ fmt(detail.stock_statement.deposit) }}</div>
          <div class="text-neutral-500">出金</div><div class="text-right text-neutral-200">{{ fmt(detail.stock_statement.withdrawal) }}</div>
          <div class="text-neutral-400 font-bold border-t border-neutral-700 pt-1 mt-1">損益</div>
          <div class="text-right font-bold border-t border-neutral-700 pt-1 mt-1" :class="detail.stock_statement.profit_loss >= 0 ? 'text-emerald-400' : 'text-red-400'">
            {{ formatPL(detail.stock_statement.profit_loss) }}
          </div>
        </div>
        <div v-else class="text-sm text-neutral-600 italic px-1">尚未填寫</div>
      </section>

      <!-- Member deposits/withdrawals for this month (from allocations) -->
      <section class="mb-6">
        <div class="flex items-center justify-between mb-2">
          <h3 class="text-sm font-bold text-neutral-300">成員入出金</h3>
          <NuxtLink
            to="/investments/transactions"
            class="text-xs text-indigo-400 font-bold no-underline"
          >
            前往異動
          </NuxtLink>
        </div>
        <div v-if="hasDeposits" class="bg-neutral-900 rounded-xl p-3 border border-neutral-800">
          <div v-for="a in detail.allocations" :key="a.member_id" class="flex items-center justify-between text-sm py-1">
            <span class="text-neutral-300">{{ a.member_name }}</span>
            <div class="flex gap-4 text-xs">
              <span v-if="a.deposit" class="text-emerald-400">入 {{ fmt(a.deposit) }}</span>
              <span v-if="a.withdrawal" class="text-red-400">出 {{ fmt(a.withdrawal) }}</span>
              <span v-if="!a.deposit && !a.withdrawal" class="text-neutral-700">-</span>
            </div>
          </div>
        </div>
        <div v-else class="text-sm text-neutral-600 italic px-1">本月無入出金</div>
      </section>

      <!-- Allocation preview -->
      <section class="mb-6">
        <h3 class="text-sm font-bold text-neutral-300 mb-2">分配預覽</h3>
        <div class="bg-neutral-900 rounded-xl p-3 border border-neutral-800">
          <div class="flex items-center justify-between text-sm mb-2 pb-2 border-b border-neutral-800">
            <span class="text-neutral-400">總損益</span>
            <span class="font-bold" :class="detail.total_profit_loss >= 0 ? 'text-emerald-400' : 'text-red-400'">
              {{ formatPL(detail.total_profit_loss) }}
            </span>
          </div>
          <div v-for="a in detail.allocations" :key="a.member_id" class="flex items-center justify-between text-sm py-1">
            <span class="text-neutral-300">{{ a.member_name || a.member_id }}</span>
            <div class="flex items-center gap-3">
              <span class="text-xs text-neutral-500">權重 {{ a.weight }}</span>
              <span class="font-bold" :class="a.amount >= 0 ? 'text-emerald-400' : 'text-red-400'">
                {{ formatPL(a.amount) }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- Actions -->
      <div class="flex gap-3 mt-4">
        <button
          v-if="detail.status === 'draft'"
          type="button"
          @click="handleComplete"
          class="flex-1 py-3 rounded-xl font-bold text-sm bg-indigo-500 text-white border-0 cursor-pointer active:scale-[0.98] transition-transform"
        >
          完成結算
        </button>
        <button
          v-if="detail.status === 'completed'"
          type="button"
          @click="handleReopen"
          class="flex-1 py-3 rounded-xl font-bold text-sm bg-amber-500/20 text-amber-400 border border-amber-500/30 cursor-pointer active:scale-[0.98] transition-transform"
        >
          重新開啟
        </button>
        <button
          v-if="detail.status === 'draft'"
          type="button"
          @click="handleDelete"
          class="py-3 px-4 rounded-xl font-bold text-sm bg-red-500/10 text-red-400 border border-red-500/20 cursor-pointer active:scale-[0.98] transition-transform"
        >
          刪除
        </button>
      </div>
    </template>
  </OverlayPage>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useInvestment } from '~/composables/useInvestment'
import { useToast } from '~/composables/useToast'
import { useConfirm } from '~/composables/useConfirm'
import PageTitle from '~/components/PageTitle.vue'
import OverlayPage from '~/components/OverlayPage.vue'
import type { InvSettlementDetail } from '~/types'

const route = useRoute()
const router = useRouter()
const ym = route.params.ym as string

const { getSettlement, completeSettlement, reopenSettlement, deleteSettlement } = useInvestment()
const { show: showToast } = useToast()
const confirm = useConfirm()

const loading = ref(true)
const detail = ref<InvSettlementDetail | null>(null)

const fmt = (n: number) => n.toLocaleString()
const formatPL = (n: number) => (n > 0 ? '+' : '') + n.toLocaleString()

const hasDeposits = computed(() => {
  return detail.value?.allocations?.some(a => a.deposit || a.withdrawal) ?? false
})

const handleComplete = async () => {
  const ok = await confirm({ message: '完成後將寫入分配損益紀錄' })
  if (!ok) return
  try {
    await completeSettlement(ym)
    showToast('結算完成', 'success')
    detail.value = await getSettlement(ym)
  } catch (e: any) {
    showToast(e.message || '操作失敗', 'error')
  }
}

const handleReopen = async () => {
  const ok = await confirm({ message: '將回到草稿狀態' })
  if (!ok) return
  try {
    await reopenSettlement(ym)
    showToast('已重新開啟', 'success')
    detail.value = await getSettlement(ym)
  } catch (e: any) {
    showToast(e.message || '操作失敗', 'error')
  }
}

const handleDelete = async () => {
  const ok = await confirm({ message: '確認刪除？此操作無法復原', destructive: true })
  if (!ok) return
  try {
    await deleteSettlement(ym)
    showToast('已刪除', 'success')
    router.back()
  } catch (e: any) {
    showToast(e.message || '刪除失敗', 'error')
  }
}

onMounted(async () => {
  try {
    detail.value = await getSettlement(ym)
  } finally {
    loading.value = false
  }
})
</script>
