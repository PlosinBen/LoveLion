<template>
  <div>
    <PageTitle title="出入金" :show-back="false" />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="memberTransactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <Icon icon="mdi:transfer" class="text-5xl opacity-20 mb-4" />
      <p>尚無出入金紀錄</p>
    </div>

    <div v-else class="flex flex-col gap-1">
      <component
        :is="isReadOnly ? 'div' : 'button'"
        v-for="t in memberTransactions"
        :key="t.id"
        :type="isReadOnly ? undefined : 'button'"
        @click="isReadOnly ? undefined : router.push(`/investments/transactions/${t.id}/edit`)"
        class="w-full flex items-center justify-between px-4 py-3 bg-neutral-900 border border-neutral-800 rounded-xl"
        :class="isReadOnly ? '' : 'cursor-pointer transition-colors hover:bg-neutral-800/70 active:scale-[0.99]'"
      >
        <div class="flex flex-col items-start gap-0.5">
          <div class="flex items-center gap-2">
            <span
              class="text-xs font-bold px-1.5 py-0.5 rounded"
              :class="typeStyle(t.type)"
            >
              {{ typeLabel(t.type) }}
            </span>
            <span class="text-sm text-neutral-200">{{ t.member?.name || t.member_id }}</span>
          </div>
          <span class="text-xs text-neutral-500">{{ t.date.slice(0, 10) }}</span>
        </div>
        <div class="flex flex-col items-end gap-0.5">
          <span class="text-sm font-bold" :class="t.amount < 0 ? 'text-red-400' : 'text-emerald-400'">
            {{ t.amount.toLocaleString() }}
          </span>
          <span v-if="t.note" class="text-xs text-neutral-500">{{ t.note }}</span>
        </div>
      </component>
    </div>

    <BaseFab v-if="!isReadOnly" @click="router.push('/investments/transactions/add')" />

    <NuxtPage :transition="{ name: 'slide-right' }" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useInvestment } from '~/composables/useInvestment'
import PageTitle from '~/components/PageTitle.vue'
import BaseFab from '~/components/BaseFab.vue'

const route = useRoute()
const router = useRouter()
const { user } = useAuth()
const { members, memberTransactions, fetchMembers, fetchMemberTransactions } = useInvestment()

const loading = ref(true)

const memberQuery = computed(() => route.query.member as string | undefined)
const isReadOnly = computed(() => !!memberQuery.value || !user.value?.inv_is_owner)

const typeLabel = (t: string) => {
  if (t === 'deposit') return '入金'
  if (t === 'withdrawal') return '出金'
  if (t === 'fee') return '費用'
  return '損益'
}

const typeStyle = (t: string) => {
  if (t === 'deposit') return 'bg-emerald-500/20 text-emerald-400'
  if (t === 'withdrawal') return 'bg-red-500/20 text-red-400'
  if (t === 'fee') return 'bg-amber-500/20 text-amber-400'
  return 'bg-indigo-500/20 text-indigo-400'
}

onMounted(async () => {
  try {
    let memberId: string | undefined
    if (memberQuery.value) {
      await fetchMembers()
      const match = members.value.find(m => m.id === memberQuery.value || m.name === memberQuery.value)
      memberId = match?.id
    }
    await fetchMemberTransactions(memberId ? { member_id: memberId } : undefined)
  } finally {
    loading.value = false
  }
})
</script>
