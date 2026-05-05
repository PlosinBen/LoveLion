<template>
  <form @submit.prevent="handleSave" class="flex flex-col gap-4">
    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">成員</label>
      <select
        v-model="form.member_id"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
        <option value="" disabled>選擇成員</option>
        <option v-for="m in members" :key="m.id" :value="m.id">{{ m.name }}</option>
      </select>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">類型</label>
      <div class="flex gap-2">
        <button
          v-for="opt in typeOptions"
          :key="opt.value"
          type="button"
          @click="form.type = opt.value"
          class="flex-1 py-2 rounded-lg border text-sm font-bold cursor-pointer transition-colors"
          :class="form.type === opt.value
            ? 'bg-indigo-500/15 text-indigo-400 border-indigo-500/30'
            : 'bg-transparent text-neutral-500 border-neutral-700'"
        >
          {{ opt.label }}
        </button>
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">日期</label>
      <input
        v-model="form.date"
        type="date"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-xs font-bold text-neutral-400">金額</label>
      <input
        v-model.number="form.amount"
        type="number"
        class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
      >
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
      v-if="transactionId"
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
  transactionId?: string
}>()

const emit = defineEmits<{
  saved: []
  deleted: []
}>()

const { members, memberTransactions, fetchMembers, fetchMemberTransactions, createMemberTransaction, updateMemberTransaction, deleteMemberTransaction } = useInvestment()
const { show: showToast } = useToast()
const confirm = useConfirm()

const saving = ref(false)
const form = reactive({
  member_id: '',
  type: 'deposit',
  date: new Date().toISOString().slice(0, 10),
  amount: 0,
  note: '',
})

const typeOptions = [
  { value: 'deposit', label: '入金' },
  { value: 'withdrawal', label: '出金' },
]

const handleSave = async () => {
  if (!form.member_id || !form.amount) {
    showToast('請填寫成員和金額', 'error')
    return
  }
  saving.value = true
  try {
    if (props.transactionId) {
      await updateMemberTransaction(props.transactionId, form)
    } else {
      await createMemberTransaction(form)
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
  const ok = await confirm({ message: '確認刪除此異動？', destructive: true })
  if (!ok) return
  try {
    await deleteMemberTransaction(props.transactionId!)
    showToast('已刪除', 'success')
    emit('deleted')
  } catch (e: any) {
    showToast(e.message || '刪除失敗', 'error')
  }
}

onMounted(async () => {
  await fetchMembers()
  if (props.transactionId) {
    if (memberTransactions.value.length === 0) await fetchMemberTransactions()
    const txn = memberTransactions.value.find(t => t.id === props.transactionId)
    if (txn) {
      form.member_id = txn.member_id
      form.type = txn.type
      form.date = txn.date
      form.amount = txn.amount
      form.note = txn.note
    }
  }
})
</script>
