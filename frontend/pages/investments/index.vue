<template>
  <div>
    <PageTitle
      title="投資損益"
      :show-back="false"
      :settings-to="user?.inv_is_owner && !viewMember ? '/investments/settings' : undefined"
    />

    <div v-if="viewMember" class="flex items-center gap-2 mb-3 px-1">
      <span class="text-xs text-amber-400 bg-amber-500/15 px-2 py-1 rounded font-bold">模擬：{{ viewMember.name }}</span>
      <NuxtLink to="/investments" class="text-xs text-neutral-500 no-underline">清除</NuxtLink>
    </div>

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="grouped.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <Icon icon="mdi:chart-line" class="text-5xl opacity-20 mb-4" />
      <p>尚無損益紀錄</p>
    </div>

    <div v-else-if="viewMember" class="flex flex-col gap-3">
      <div
        v-for="row in grouped"
        :key="row.yearMonth"
        class="bg-neutral-900 rounded-xl border border-neutral-800 overflow-hidden"
      >
        <div class="px-4 py-2.5 border-b border-neutral-800 font-mono text-sm text-neutral-300 font-bold">
          {{ row.yearMonth }}
        </div>
        <div v-if="row.data[viewMember.id]" class="grid grid-cols-2 gap-x-4 gap-y-1 text-sm px-4 py-3">
          <template v-if="row.data[viewMember.id]!.deposit">
            <span class="text-neutral-500">入金</span>
            <span class="text-right text-emerald-400">{{ row.data[viewMember.id]!.deposit.toLocaleString() }}</span>
          </template>
          <template v-if="row.data[viewMember.id]!.withdrawal">
            <span class="text-neutral-500">出金</span>
            <span class="text-right text-red-400">{{ row.data[viewMember.id]!.withdrawal.toLocaleString() }}</span>
          </template>
          <template v-if="row.data[viewMember.id]!.fee">
            <span class="text-neutral-500">費用</span>
            <span class="text-right text-amber-400">{{ row.data[viewMember.id]!.fee.toLocaleString() }}</span>
          </template>
          <span class="text-neutral-500">損益</span>
          <span class="text-right" :class="row.data[viewMember.id]!.amount >= 0 ? 'text-red-400' : 'text-emerald-400'">
            {{ formatAmount(row.data[viewMember.id]!.amount) }}
          </span>
          <span class="text-neutral-400 font-bold border-t border-neutral-800 pt-1 mt-1">結餘</span>
          <span class="text-right text-neutral-200 font-bold border-t border-neutral-800 pt-1 mt-1">
            {{ row.data[viewMember.id]!.balance.toLocaleString() }}
          </span>
        </div>
      </div>
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
            v-for="m in displayMembers"
            :key="m.id + '-' + row.yearMonth"
            class="px-4 py-2.5"
          >
            <template v-if="row.data[m.id]">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-sm text-neutral-300 font-medium">{{ m.name }}</span>
                  <span v-if="row.data[m.id]!.deposit" class="text-xs text-emerald-400">入 {{ row.data[m.id]!.deposit.toLocaleString() }}</span>
                  <span v-if="row.data[m.id]!.withdrawal" class="text-xs text-red-400">出 {{ row.data[m.id]!.withdrawal.toLocaleString() }}</span>
                  <span v-if="row.data[m.id]!.fee" class="text-xs text-amber-400">費 {{ row.data[m.id]!.fee.toLocaleString() }}</span>
                </div>
              </div>
              <div class="flex items-center justify-end gap-4 mt-0.5">
                <span class="text-xs" :class="row.data[m.id]!.amount >= 0 ? 'text-red-400' : 'text-emerald-400'">
                  損益 {{ formatAmount(row.data[m.id]!.amount) }}
                </span>
                <span class="text-xs text-neutral-200 font-bold">
                  結餘 {{ row.data[m.id]!.balance.toLocaleString() }}
                </span>
              </div>
            </template>
            <div v-else class="flex items-center justify-between">
              <span class="text-sm text-neutral-300 font-medium">{{ m.name }}</span>
              <span class="text-xs text-neutral-700">-</span>
            </div>
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

const route = useRoute()
const { user } = useAuth()
const { allocations, members: memberList, fetchAllocations, fetchMembers } = useInvestment()

const loading = ref(true)
const members = ref<InvMember[]>([])

const viewMemberQuery = computed(() => route.query.member as string | undefined)
const viewMember = computed(() => {
  if (!viewMemberQuery.value) return undefined
  const q = viewMemberQuery.value
  return members.value.find(m => m.id === q || m.name === q)
})
const displayMembers = computed(() => viewMember.value ? [viewMember.value] : members.value)

interface GroupedRow {
  yearMonth: string
  data: Record<string, InvSettlementAllocation>
}

const grouped = computed<GroupedRow[]>(() => {
  const map = new Map<string, Record<string, InvSettlementAllocation>>()
  for (const a of allocations.value) {
    if (viewMember.value && a.member_id !== viewMember.value.id) continue
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

