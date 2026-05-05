<template>
  <OverlayPage>
    <PageTitle
      :title="`期貨 ${ym}`"
      :show-back="true"
      :breadcrumbs="[{ label: '結算', to: '/investments/settlements' }, { label: ym, to: `/investments/settlements/${ym}` }]"
    />

    <form @submit.prevent="handleSave" class="flex flex-col gap-4">
      <div v-for="field in fields" :key="field.key" class="flex flex-col gap-1">
        <label class="text-xs font-bold text-neutral-400">{{ field.label }}</label>
        <input
          v-model.number="form[field.key]"
          type="number"
          class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
        >
      </div>

      <div class="bg-neutral-900 rounded-xl p-3 border border-neutral-800 text-sm">
        <div class="flex justify-between text-neutral-500">
          <span>計算損益</span>
          <span class="font-bold" :class="computedPL >= 0 ? 'text-emerald-400' : 'text-red-400'">
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

const { getSettlement, upsertFutures } = useInvestment()
const { show: showToast } = useToast()

const saving = ref(false)
const form = reactive({
  ending_equity: 0,
  floating_profit_loss: 0,
  realized_profit_loss: 0,
  deposit: 0,
  withdrawal: 0,
})

const fields = [
  { key: 'ending_equity' as const, label: '期末權益' },
  { key: 'floating_profit_loss' as const, label: '浮動損益' },
  { key: 'realized_profit_loss' as const, label: '沖銷損益' },
  { key: 'deposit' as const, label: '入金' },
  { key: 'withdrawal' as const, label: '出金' },
]

const computedPL = computed(() => {
  const realEquity = form.ending_equity - form.floating_profit_loss
  // We don't have prev real equity on the frontend, so just show realized for now
  // The backend computes the actual PL
  return form.realized_profit_loss
})

const handleSave = async () => {
  saving.value = true
  try {
    await upsertFutures(ym, form)
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
    if (detail.futures_statement) {
      form.ending_equity = detail.futures_statement.ending_equity
      form.floating_profit_loss = detail.futures_statement.floating_profit_loss
      form.realized_profit_loss = detail.futures_statement.realized_profit_loss
      form.deposit = detail.futures_statement.deposit
      form.withdrawal = detail.futures_statement.withdrawal
    }
  } catch (e) {
    // Settlement not found
  }
})
</script>
