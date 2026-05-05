<template>
  <div>
    <PageTitle title="成員異動" :show-back="false" />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="memberTransactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <Icon icon="mdi:transfer" class="text-5xl opacity-20 mb-4" />
      <p>尚無異動紀錄</p>
    </div>

    <div v-else class="flex flex-col gap-1">
      <button
        v-for="t in memberTransactions"
        :key="t.id"
        type="button"
        @click="router.push(`/investments/transactions/${t.id}/edit`)"
        class="w-full flex items-center justify-between px-4 py-3 bg-neutral-900 border border-neutral-800 rounded-xl cursor-pointer transition-colors hover:bg-neutral-800/70 active:scale-[0.99]"
      >
        <div class="flex flex-col items-start gap-0.5">
          <div class="flex items-center gap-2">
            <span
              class="text-xs font-bold px-1.5 py-0.5 rounded"
              :class="typeStyle(t.type)"
            >
              {{ typeLabel(t.type) }}
            </span>
            <span class="text-sm text-neutral-200">{{ t.member_name || t.member_id }}</span>
          </div>
          <span class="text-xs text-neutral-500">{{ t.date }}</span>
        </div>
        <div class="flex flex-col items-end gap-0.5">
          <span class="text-sm font-bold" :class="t.type === 'withdrawal' ? 'text-red-400' : 'text-emerald-400'">
            {{ t.amount.toLocaleString() }}
          </span>
          <span v-if="t.note" class="text-xs text-neutral-500">{{ t.note }}</span>
        </div>
      </button>
    </div>

    <BaseFab @click="router.push('/investments/transactions/add')" />

    <Transition name="slide-right">
      <NuxtPage />
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useInvestment } from '~/composables/useInvestment'
import PageTitle from '~/components/PageTitle.vue'
import BaseFab from '~/components/BaseFab.vue'

const router = useRouter()
const { memberTransactions, fetchMemberTransactions } = useInvestment()
const loading = ref(true)

const typeLabel = (t: string) => {
  if (t === 'deposit') return '入金'
  if (t === 'withdrawal') return '出金'
  return '損益'
}

const typeStyle = (t: string) => {
  if (t === 'deposit') return 'bg-emerald-500/20 text-emerald-400'
  if (t === 'withdrawal') return 'bg-red-500/20 text-red-400'
  return 'bg-indigo-500/20 text-indigo-400'
}

onMounted(async () => {
  try {
    await fetchMemberTransactions()
  } finally {
    loading.value = false
  }
})
</script>
