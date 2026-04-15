<template>
  <div>
    <PageTitle
      title="編輯公告"
      back-to="/admin/announcements"
      :breadcrumbs="[{ label: '設定', to: '/settings' }, { label: '公告管理', to: '/admin/announcements' }]"
    >
      <template #right>
        <button
          @click="confirmDelete"
          title="刪除公告"
          class="w-10 h-10 rounded-full text-rose-500 flex items-center justify-center hover:bg-rose-500/10 transition-colors cursor-pointer active:scale-95"
        >
          <Icon icon="mdi:delete-outline" class="text-xl" />
        </button>
      </template>
    </PageTitle>

    <div v-if="loading" class="p-20 flex justify-center items-center text-neutral-500">
      <Icon icon="mdi:loading" class="text-4xl animate-spin" />
    </div>

    <AnnouncementForm
      v-else-if="announcement"
      :initial-data="announcement"
      @save="handleUpdate"
    />

    <div v-else class="py-20 flex flex-col items-center justify-center text-neutral-600">
      <Icon icon="mdi:alert-circle-outline" class="text-5xl mb-3 opacity-20" />
      <span class="text-sm font-bold">找不到此公告</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import PageTitle from '~/components/PageTitle.vue'
import AnnouncementForm from '~/components/AnnouncementForm.vue'
import { useToast } from '~/composables/useToast'
import { useLoading } from '~/composables/useLoading'
import type { Announcement } from '~/types'

definePageMeta({
  layout: 'default'
})

const route = useRoute()
const router = useRouter()
const api = useApi()
const toast = useToast()
const { showLoading, hideLoading } = useLoading()

const announcement = ref<Announcement | null>(null)
const loading = ref(true)

const handleUpdate = async (payload: any) => {
  showLoading()
  try {
    await api.put(`/api/admin/announcements/${route.params.id}`, payload)
    toast.success('公告已更新')
    router.push('/admin/announcements')
  } catch (e: any) {
    toast.error(e.message || '更新失敗')
  } finally {
    hideLoading()
  }
}

const confirmDelete = async () => {
  if (!confirm('確定要刪除此公告嗎？此操作無法復原。')) return

  showLoading()
  try {
    await api.del(`/api/admin/announcements/${route.params.id}`)
    toast.success('公告已刪除')
    router.push('/admin/announcements')
  } catch (e: any) {
    toast.error(e.message || '刪除失敗')
  } finally {
    hideLoading()
  }
}

onMounted(async () => {
  try {
    // Use admin endpoint to fetch (can see drafts)
    const list = await api.get<Announcement[]>('/api/admin/announcements')
    announcement.value = list.find(a => a.id === route.params.id) || null
  } catch (e) {
    console.error('Failed to fetch announcement:', e)
  } finally {
    loading.value = false
  }
})
</script>
