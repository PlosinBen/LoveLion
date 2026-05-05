<template>
  <div>
    <PageTitle
      title="投資損益"
      :show-back="true"
      back-to="/"
      :settings-to="user?.inv_is_owner ? '/investments/settings' : undefined"
    />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="grouped.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <Icon icon="mdi:chart-line" class="text-5xl opacity-20 mb-4" />
      <p>尚無損益紀錄</p>
    </div>

    <div v-else class="overflow-x-auto no-scrollbar rounded-xl border border-neutral-800">
      <table class="w-full text-sm whitespace-nowrap">
        <thead>
          <tr class="bg-neutral-900 text-neutral-400 text-xs">
            <th class="sticky left-0 z-10 bg-neutral-900 px-3 py-2 text-left font-bold">月份</th>
            <template v-for="m in members" :key="m.id">
              <th class="px-2 py-2 text-right font-bold border-l border-neutral-800" colspan="4">{{ m.name }}</th>
            </template>
          </tr>
          <tr class="bg-neutral-900/50 text-neutral-500 text-xs">
            <th class="sticky left-0 z-10 bg-neutral-900/50 px-3 py-1"></th>
            <template v-for="m in members" :key="m.id + '-sub'">
              <th class="px-2 py-1 text-right border-l border-neutral-800">入金</th>
              <th class="px-2 py-1 text-right">出金</th>
              <th class="px-2 py-1 text-right">損益</th>
              <th class="px-2 py-1 text-right">結餘</th>
            </template>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in grouped" :key="row.yearMonth" class="border-t border-neutral-800/50 hover:bg-neutral-900/30">
            <td class="sticky left-0 z-10 bg-neutral-950 px-3 py-2 font-mono text-xs text-neutral-300">{{ row.yearMonth }}</td>
            <template v-for="m in members" :key="m.id + '-' + row.yearMonth">
              <td class="px-2 py-2 text-right text-xs border-l border-neutral-800/50">
                <span v-if="row.data[m.id]?.deposit" class="text-emerald-400">{{ formatAmount(row.data[m.id]!.deposit) }}</span>
                <span v-else class="text-neutral-700">-</span>
              </td>
              <td class="px-2 py-2 text-right text-xs">
                <span v-if="row.data[m.id]?.withdrawal" class="text-red-400">{{ formatAmount(row.data[m.id]!.withdrawal) }}</span>
                <span v-else class="text-neutral-700">-</span>
              </td>
              <td class="px-2 py-2 text-right text-xs">
                <span v-if="row.data[m.id]" :class="row.data[m.id]!.amount >= 0 ? 'text-emerald-400' : 'text-red-400'">
                  {{ formatAmount(row.data[m.id]!.amount) }}
                </span>
                <span v-else class="text-neutral-700">-</span>
              </td>
              <td class="px-2 py-2 text-right text-xs">
                <span v-if="row.data[m.id]" class="text-neutral-200">{{ formatAmount(row.data[m.id]!.balance) }}</span>
                <span v-else class="text-neutral-700">-</span>
              </td>
            </template>
          </tr>
        </tbody>
      </table>
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

<style scoped>
.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>
