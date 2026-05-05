<template>
  <OverlayPage>
    <PageTitle
      title="成員管理"
      :show-back="true"
      back-to="/investments"
    />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <template v-else>
      <div class="flex flex-col gap-2 mb-6">
        <div
          v-for="m in members"
          :key="m.id"
          class="flex items-center justify-between px-4 py-3 bg-neutral-900 border border-neutral-800 rounded-xl"
        >
          <div class="flex items-center gap-3">
            <span class="text-sm text-neutral-200">{{ m.name }}</span>
            <span v-if="m.is_owner" class="text-xs bg-indigo-500/20 text-indigo-400 px-1.5 py-0.5 rounded">Owner</span>
            <span v-if="!m.active" class="text-xs bg-neutral-700/50 text-neutral-500 px-1.5 py-0.5 rounded">停用</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs text-neutral-500">{{ m.net_investment.toLocaleString() }}</span>
            <button
              type="button"
              @click="toggleActive(m)"
              class="text-xs px-2 py-1 rounded border cursor-pointer transition-colors"
              :class="m.active
                ? 'text-red-400 border-red-500/30 bg-red-500/10'
                : 'text-emerald-400 border-emerald-500/30 bg-emerald-500/10'"
            >
              {{ m.active ? '停用' : '啟用' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Add member -->
      <div class="flex gap-2">
        <input
          v-model="newName"
          type="text"
          placeholder="新成員名稱"
          class="flex-1 bg-neutral-900 border border-neutral-800 text-white text-sm py-2.5 px-3 rounded-xl focus:outline-none focus:border-indigo-500"
          @keyup.enter="handleAdd"
        >
        <button
          type="button"
          @click="handleAdd"
          :disabled="!newName.trim()"
          class="px-4 py-2.5 rounded-xl font-bold text-sm bg-indigo-500 text-white border-0 cursor-pointer active:scale-[0.98] transition-transform disabled:opacity-50"
        >
          新增
        </button>
      </div>
    </template>
  </OverlayPage>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useInvestment } from '~/composables/useInvestment'
import { useToast } from '~/composables/useToast'
import PageTitle from '~/components/PageTitle.vue'
import OverlayPage from '~/components/OverlayPage.vue'

const { members, fetchMembers, createMember, updateMember } = useInvestment()
const { show: showToast } = useToast()

const loading = ref(true)
const newName = ref('')

const handleAdd = async () => {
  const name = newName.value.trim()
  if (!name) return
  try {
    await createMember({ name })
    newName.value = ''
    await fetchMembers()
    showToast('已新增', 'success')
  } catch (e: any) {
    showToast(e.message || '新增失敗', 'error')
  }
}

const toggleActive = async (m: any) => {
  try {
    await updateMember(m.id, { active: !m.active })
    await fetchMembers()
  } catch (e: any) {
    showToast(e.message || '更新失敗', 'error')
  }
}

onMounted(async () => {
  try {
    await fetchMembers()
  } finally {
    loading.value = false
  }
})
</script>
