<template>
  <div>
    <PageTitle title="結算列表" :show-back="false" />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="settlements.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <Icon icon="mdi:calculator-variant" class="text-5xl opacity-20 mb-4" />
      <p>尚無結算紀錄</p>
    </div>

    <div v-else class="flex flex-col gap-1">
      <button
        v-for="s in settlements"
        :key="s.year_month"
        type="button"
        @click="router.push(`/investments/settlements/${s.year_month}`)"
        class="w-full flex items-center justify-between px-4 py-3 bg-neutral-900 border border-neutral-800 rounded-xl cursor-pointer transition-colors hover:bg-neutral-800/70 active:scale-[0.99]"
      >
        <span class="font-mono text-sm text-neutral-200">{{ s.year_month }}</span>
        <div class="flex items-center gap-3">
          <span class="text-sm font-bold" :class="s.total_profit_loss >= 0 ? 'text-emerald-400' : 'text-red-400'">
            {{ s.total_profit_loss ? formatAmount(s.total_profit_loss) : '-' }}
          </span>
          <span
            class="w-2.5 h-2.5 rounded-full"
            :class="s.status === 'completed' ? 'bg-emerald-500' : 'bg-amber-500'"
          />
        </div>
      </button>
    </div>

    <BaseFab @click="handleCreate" />

    <Transition name="slide-right">
      <NuxtPage />
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useInvestment } from '~/composables/useInvestment'
import { useToast } from '~/composables/useToast'
import { usePrompt } from '~/composables/usePrompt'
import PageTitle from '~/components/PageTitle.vue'
import BaseFab from '~/components/BaseFab.vue'

const router = useRouter()
const { settlements, fetchSettlements, createSettlement } = useInvestment()
const { show: showToast } = useToast()
const prompt = usePrompt()

const loading = ref(true)

const formatAmount = (n: number) => {
  const prefix = n > 0 ? '+' : ''
  return prefix + n.toLocaleString()
}

const handleCreate = async () => {
  const now = new Date()
  const defaultYM = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  const ym = await prompt({ title: '新增月結算', placeholder: 'YYYY-MM', defaultValue: defaultYM })
  if (!ym) return
  if (!/^\d{4}-\d{2}$/.test(ym)) {
    showToast('格式錯誤，請使用 YYYY-MM', 'error')
    return
  }
  try {
    await createSettlement(ym)
    await fetchSettlements()
    router.push(`/investments/settlements/${ym}`)
  } catch (e: any) {
    showToast(e.message || '建立失敗', 'error')
  }
}

onMounted(async () => {
  try {
    await fetchSettlements()
  } finally {
    loading.value = false
  }
})
</script>
