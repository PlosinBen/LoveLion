<template>
  <form @submit.prevent="handleSave" class="flex flex-col gap-4">
    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">交易日期</label>
      <input
        v-model="form.trade_date"
        type="date"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">股票代號</label>
      <input
        v-model="form.symbol"
        type="text"
        placeholder="例: 2330"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">股數 (正=買, 負=賣)</label>
      <input
        v-model.number="form.shares"
        type="number"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">成交價</label>
      <input
        v-model.number="form.price"
        type="number"
        step="0.01"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div class="flex flex-col gap-1">
        <label class="text-xs font-bold text-neutral-400">手續費</label>
        <input
          v-model.number="form.fee"
          type="number"
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
        >
      </div>
      <div class="flex flex-col gap-1">
        <label class="text-xs font-bold text-neutral-400">交易稅</label>
        <input
          v-model.number="form.tax"
          type="number"
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
        >
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">備註</label>
      <input
        v-model="form.note"
        type="text"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
    </div>

    <button
      type="submit"
      :disabled="saving"
      class="w-full py-3 rounded-xl font-bold text-sm bg-indigo-500 text-white border-0 cursor-pointer active:scale-[0.98] transition-transform disabled:opacity-50"
    >
      {{ saving ? '儲存中...' : '儲存' }}
    </button>

    <button
      v-if="tradeId"
      type="button"
      @click="handleDelete"
      class="w-full py-3 rounded-xl font-bold text-sm bg-red-500/10 text-red-400 border border-red-500/20 cursor-pointer active:scale-[0.98] transition-transform"
    >
      刪除
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useInvestment } from '~/composables/useInvestment'
import { useToast } from '~/composables/useToast'
import { useConfirm } from '~/composables/useConfirm'

const props = defineProps<{
  tradeId?: string
}>()

const emit = defineEmits<{
  saved: []
  deleted: []
}>()

const { stockTrades, fetchStockTrades, createStockTrade, updateStockTrade, deleteStockTrade } = useInvestment()
const { show: showToast } = useToast()
const confirm = useConfirm()

const saving = ref(false)
const form = reactive({
  trade_date: new Date().toISOString().slice(0, 10),
  symbol: '',
  shares: 0,
  price: 0,
  fee: 0,
  tax: 0,
  note: '',
})

const handleSave = async () => {
  if (!form.symbol || !form.shares) {
    showToast('請填寫代號和股數', 'error')
    return
  }
  saving.value = true
  try {
    if (props.tradeId) {
      await updateStockTrade(props.tradeId, form)
    } else {
      await createStockTrade(form)
    }
    showToast('已儲存', 'success')
    emit('saved')
  } catch (e: any) {
    showToast(e.message || '儲存失敗', 'error')
  } finally {
    saving.value = false
  }
}

const handleDelete = async () => {
  const ok = await confirm({ message: '確認刪除此交易？', destructive: true })
  if (!ok) return
  try {
    await deleteStockTrade(props.tradeId!)
    showToast('已刪除', 'success')
    emit('deleted')
  } catch (e: any) {
    showToast(e.message || '刪除失敗', 'error')
  }
}

onMounted(async () => {
  if (props.tradeId) {
    if (stockTrades.value.length === 0) await fetchStockTrades()
    const trade = stockTrades.value.find(t => t.id === props.tradeId)
    if (trade) {
      form.trade_date = trade.trade_date
      form.symbol = trade.symbol
      form.shares = trade.shares
      form.price = trade.price
      form.fee = trade.fee
      form.tax = trade.tax
      form.note = trade.note
    }
  }
})
</script>
