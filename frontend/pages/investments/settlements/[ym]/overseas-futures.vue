<template>
  <OverlayPage>
    <PageTitle
      :title="`群益海期 ${ym}`"
      :show-back="true"
      :breadcrumbs="[{ label: '結算', to: '/investments/settlements' }, { label: ym, to: `/investments/settlements/${ym}` }]"
    />

    <form @submit.prevent="handleSave" class="flex flex-col gap-4">
      <div class="flex flex-col gap-1">
        <label class="text-xs font-bold text-neutral-400">台幣餘額</label>
        <input
          :value="form.twd_balance"
          @input="form.twd_balance = toInt(($event.target as HTMLInputElement).value)"
          @focus="($event.target as HTMLInputElement).select()"
          type="text"
          inputmode="numeric"
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
        >
      </div>

      <div v-for="(cur, idx) in form.currencies" :key="idx" class="bg-neutral-900 rounded-xl p-3 border border-neutral-800 flex flex-col gap-3">
        <div class="flex items-center justify-between">
          <span class="text-sm font-bold text-neutral-300">{{ cur.currency }}</span>
          <button
            v-if="form.currencies.length > 1"
            type="button"
            @click="form.currencies.splice(idx, 1)"
            class="text-xs text-red-400 bg-transparent border-0 cursor-pointer"
          >
            移除
          </button>
        </div>
        <div class="grid grid-cols-3 gap-2">
          <div class="flex flex-col gap-1">
            <label class="text-xs text-neutral-500">結餘</label>
            <input
              :value="cur.balance"
              @input="cur.balance = toDecimal(($event.target as HTMLInputElement).value)"
              @focus="($event.target as HTMLInputElement).select()"
              type="text"
              inputmode="decimal"
              class="w-full bg-neutral-800 border border-neutral-700 text-white text-sm py-2 px-2 rounded-lg focus:outline-none focus:border-indigo-500"
            >
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-xs text-neutral-500">未沖銷</label>
            <input
              :value="cur.unrealized"
              @input="cur.unrealized = toDecimal(($event.target as HTMLInputElement).value)"
              @focus="($event.target as HTMLInputElement).select()"
              type="text"
              inputmode="decimal"
              class="w-full bg-neutral-800 border border-neutral-700 text-white text-sm py-2 px-2 rounded-lg focus:outline-none focus:border-indigo-500"
            >
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-xs text-neutral-500">匯率</label>
            <input
              :value="cur.exchange_rate"
              @input="cur.exchange_rate = toDecimal(($event.target as HTMLInputElement).value)"
              @focus="($event.target as HTMLInputElement).select()"
              type="text"
              inputmode="decimal"
              class="w-full bg-neutral-800 border border-neutral-700 text-white text-sm py-2 px-2 rounded-lg focus:outline-none focus:border-indigo-500"
            >
          </div>
        </div>
        <div class="text-xs text-neutral-500 flex justify-between">
          <span>淨額: {{ (cur.balance - cur.unrealized).toLocaleString() }} {{ cur.currency }}</span>
          <span>台幣: {{ ((cur.balance - cur.unrealized) * cur.exchange_rate).toLocaleString() }}</span>
        </div>
      </div>

      <button
        type="button"
        @click="addCurrency"
        class="w-full py-2 rounded-xl text-sm text-indigo-400 border border-dashed border-neutral-700 bg-transparent cursor-pointer"
      >
        + 新增幣別
      </button>

      <div class="bg-neutral-900 rounded-xl p-3 border border-neutral-800 text-sm">
        <div class="flex justify-between text-neutral-500 mb-1">
          <span>轉換後淨額</span>
          <span class="text-neutral-300 font-bold">{{ convertedNet.toLocaleString() }}</span>
        </div>
        <div class="flex justify-between text-neutral-500">
          <span>計算損益</span>
          <span class="font-bold" :class="computedPL >= 0 ? 'text-red-400' : 'text-emerald-400'">
            {{ computedPL > 0 ? '+' : '' }}{{ computedPL.toLocaleString() }}
          </span>
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

const { getSettlement, upsertOverseasFutures } = useInvestment()
const { show: showToast } = useToast()

const saving = ref(false)
const prevConvertedNet = ref(0)
const hasPrev = ref(false)

const form = reactive({
  twd_balance: 0,
  currencies: [{ currency: 'USD', balance: 0, unrealized: 0, exchange_rate: 30 }] as Array<{
    currency: string
    balance: number
    unrealized: number
    exchange_rate: number
  }>,
})

const toInt = (v: string) => {
  const n = parseInt(v.replace(/,/g, ''), 10)
  return isNaN(n) ? 0 : n
}

const toDecimal = (v: string) => {
  const n = parseFloat(v.replace(/,/g, ''))
  return isNaN(n) ? 0 : Math.round(n * 100) / 100
}

const addCurrency = () => {
  form.currencies.push({ currency: '', balance: 0, unrealized: 0, exchange_rate: 1 })
}

const convertedNet = computed(() => {
  let total = form.twd_balance
  for (const cur of form.currencies) {
    total += Math.round((cur.balance - cur.unrealized) * cur.exchange_rate)
  }
  return total
})

const computedPL = computed(() => {
  if (!hasPrev.value) return 0
  return convertedNet.value - prevConvertedNet.value
})

const handleSave = async () => {
  saving.value = true
  try {
    await upsertOverseasFutures(ym, {
      twd_balance: form.twd_balance,
      currencies: form.currencies,
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
    const [y, m] = ym.split('-').map(Number)
    const prevDate = new Date(y, m - 2, 1)
    const prevYM = `${prevDate.getFullYear()}-${String(prevDate.getMonth() + 1).padStart(2, '0')}`
    try {
      const prevDetail = await getSettlement(prevYM)
      if (prevDetail.overseas_futures_statement) {
        prevConvertedNet.value = prevDetail.overseas_futures_statement.converted_net
        hasPrev.value = true
      }
    } catch {
      // No previous settlement
    }

    const detail = await getSettlement(ym)
    if (detail.overseas_futures_statement) {
      form.twd_balance = detail.overseas_futures_statement.twd_balance
      if (detail.overseas_futures_statement.currencies?.length) {
        form.currencies = detail.overseas_futures_statement.currencies.map(c => ({
          currency: c.currency,
          balance: parseFloat(String(c.balance)) || 0,
          unrealized: parseFloat(String(c.unrealized)) || 0,
          exchange_rate: parseFloat(String(c.exchange_rate)) || 0,
        }))
      }
    }
  } catch (e) {
    // Settlement not found
  }
})
</script>
