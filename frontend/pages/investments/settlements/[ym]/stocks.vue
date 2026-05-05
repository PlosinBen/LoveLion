<template>
  <OverlayPage>
    <PageTitle
      :title="`股票 ${ym}`"
      :show-back="true"
      :breadcrumbs="[{ label: '結算', to: '/investments/settlements' }, { label: ym, to: `/investments/settlements/${ym}` }]"
    />

    <form @submit.prevent="handleSave" class="flex flex-col gap-4">
      <div class="flex flex-col gap-1">
        <label class="text-xs font-bold text-neutral-400">帳戶餘額</label>
        <input
          v-model.number="form.account_balance"
          type="number"
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
        >
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-xs font-bold text-neutral-400">入金</label>
        <input
          v-model.number="form.deposit"
          type="number"
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
        >
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-xs font-bold text-neutral-400">出金</label>
        <input
          v-model.number="form.withdrawal"
          type="number"
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
        >
      </div>

      <!-- Holdings -->
      <section>
        <div class="flex items-center justify-between mb-2">
          <label class="text-xs font-bold text-neutral-400">庫存</label>
          <button
            type="button"
            @click="addHolding"
            class="text-xs text-indigo-400 font-bold bg-transparent border-0 cursor-pointer"
          >
            + 新增
          </button>
        </div>
        <div v-for="(h, i) in form.holdings" :key="i" class="flex items-center gap-2 mb-2">
          <input
            v-model="h.symbol"
            placeholder="代號"
            class="w-20 bg-neutral-900 border border-neutral-800 text-white text-sm py-2 px-2 rounded-lg focus:outline-none focus:border-indigo-500"
          >
          <input
            v-model.number="h.shares"
            type="number"
            placeholder="股數"
            class="w-20 bg-neutral-900 border border-neutral-800 text-white text-sm py-2 px-2 rounded-lg focus:outline-none focus:border-indigo-500"
          >
          <input
            v-model.number="h.closing_price"
            type="number"
            step="0.01"
            placeholder="結算價"
            class="flex-1 bg-neutral-900 border border-neutral-800 text-white text-sm py-2 px-2 rounded-lg focus:outline-none focus:border-indigo-500"
          >
          <button
            type="button"
            @click="form.holdings.splice(i, 1)"
            class="text-red-400 bg-transparent border-0 cursor-pointer text-lg"
          >
            &times;
          </button>
        </div>
      </section>

      <div class="bg-neutral-900 rounded-xl p-3 border border-neutral-800 text-sm">
        <div class="flex justify-between text-neutral-500">
          <span>庫存現值</span>
          <span class="text-neutral-200 font-bold">{{ marketValue.toLocaleString() }}</span>
        </div>
      </div>

      <button
        type="submit"
        :disabled="saving"
        class="w-full py-3 rounded-xl font-bold text-sm bg-indigo-500 text-white border-0 cursor-pointer active:scale-[0.98] transition-transform disabled:opacity-50"
      >
        {{ saving ? '儲存中...' : '儲存' }}
      </button>
    </form>
  </OverlayPage>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useInvestment } from '~/composables/useInvestment'
import { useToast } from '~/composables/useToast'
import PageTitle from '~/components/PageTitle.vue'
import OverlayPage from '~/components/OverlayPage.vue'

const route = useRoute()
const router = useRouter()
const ym = route.params.ym as string

const { getSettlement, upsertStocks } = useInvestment()
const { show: showToast } = useToast()

const saving = ref(false)

interface HoldingForm {
  symbol: string
  shares: number
  closing_price: number
}

const form = reactive({
  account_balance: 0,
  deposit: 0,
  withdrawal: 0,
  holdings: [] as HoldingForm[],
})

const marketValue = computed(() => {
  return form.holdings.reduce((sum, h) => sum + Math.round(h.shares * h.closing_price), 0)
})

const addHolding = () => {
  form.holdings.push({ symbol: '', shares: 0, closing_price: 0 })
}

const handleSave = async () => {
  saving.value = true
  try {
    await upsertStocks(ym, {
      account_balance: form.account_balance,
      deposit: form.deposit,
      withdrawal: form.withdrawal,
      holdings: form.holdings.filter(h => h.symbol),
    })
    showToast('已儲存', 'success')
    router.back()
  } catch (e: any) {
    showToast(e.message || '儲存失敗', 'error')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const detail = await getSettlement(ym)
    if (detail.stock_statement) {
      form.account_balance = detail.stock_statement.account_balance
      form.deposit = detail.stock_statement.deposit
      form.withdrawal = detail.stock_statement.withdrawal
    }
    if (detail.stock_statement?.holdings?.length) {
      form.holdings = detail.stock_statement.holdings.map(h => ({
        symbol: h.symbol,
        shares: h.shares,
        closing_price: h.closing_price,
      }))
    }
  } catch (e) {
    // Settlement not found
  }
})
</script>
