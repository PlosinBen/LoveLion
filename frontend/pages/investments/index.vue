<template>
  <div>
    <PageTitle
      title="投資損益"
      :show-back="false"
      :settings-to="user?.inv_is_owner ? '/investments/settings' : undefined"
    />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="grouped.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <Icon icon="mdi:chart-line" class="text-5xl opacity-20 mb-4" />
      <p>尚無損益紀錄</p>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div
        v-for="row in grouped"
        :key="row.yearMonth"
        class="bg-neutral-900 rounded-xl border border-neutral-800 overflow-hidden"
      >
        <div class="px-4 py-2.5 border-b border-neutral-800 font-mono text-sm text-neutral-300 font-bold">
          {{ row.yearMonth }}
        </div>
        <div class="divide-y divide-neutral-800/50">
          <div
            v-for="m in members"
            :key="m.id + '-' + row.yearMonth"
            class="px-4 py-2.5 flex items-center justify-between"
          >
            <span class="text-sm text-neutral-300 font-medium">{{ m.name }}</span>
            <div v-if="row.data[m.id]" class="flex items-center gap-3 text-xs">
              <span v-if="row.data[m.id]!.deposit" class="text-emerald-400">入 {{ row.data[m.id]!.deposit.toLocaleString() }}</span>
              <span v-if="row.data[m.id]!.withdrawal" class="text-red-400">出 {{ row.data[m.id]!.withdrawal.toLocaleString() }}</span>
              <span :class="row.data[m.id]!.amount >= 0 ? 'text-emerald-400' : 'text-red-400'">
                {{ formatAmount(row.data[m.id]!.amount) }}
              </span>
              <span class="text-neutral-200 font-bold">{{ row.data[m.id]!.balance.toLocaleString() }}</span>
            </div>
            <span v-else class="text-xs text-neutral-700">-</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useInvestment } from '~/composables/useInvestment'
import PageTitle from '~/components/PageTitle.vue'
import type { InvMember, InvSettlementAllocation } from '~/types'

const { user } = useAuth()
const { allocations, members: memberList, fetchAllocations, fetchMembers } = useInvestment()

const loading = ref(true)
const members = ref<InvMember[]>([])

interface GroupedRow {
  yearMonth: string
  data: Record<string, InvSettlementAllocation>
}

const grouped = computed<GroupedRow[]>(() => {
  const map = new Map<string, Record<string, InvSettlementAllocation>>()
  for (const a of allocations.value) {
    if (!map.has(a.year_month)) map.set(a.year_month, {})
    map.get(a.year_month)![a.member_id] = a
  }
  return Array.from(map.entries())
    .sort((a, b) => b[0].localeCompare(a[0]))
    .map(([ym, data]) => ({ yearMonth: ym, data }))
})

const formatAmount = (n: number) => {
  if (n === 0) return '0'
  const prefix = n > 0 ? '+' : ''
  return prefix + n.toLocaleString()
}

onMounted(async () => {
  try {
    if (user.value?.inv_is_owner) {
      await fetchMembers()
      members.value = memberList.value
    }
    await fetchAllocations()
    if (!user.value?.inv_is_owner && allocations.value.length > 0) {
      const uniqueIds = [...new Set(allocations.value.map(a => a.member_id))]
      members.value = uniqueIds.map(id => ({
        id,
        name: allocations.value.find(a => a.member_id === id)?.member_name || id,
      } as InvMember))
    }
  } finally {
    loading.value = false
  }
})
</script>

