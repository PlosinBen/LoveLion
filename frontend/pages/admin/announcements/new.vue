<template>
  <div class="admin-announcement-form-page">
    <PageTitle
      title="新增公告"
      back-to="/admin/announcements"
      :breadcrumbs="[{ label: '設定', to: '/settings' }, { label: '公告管理', to: '/admin/announcements' }]"
    />

    <AnnouncementForm
      :initial-data="null"
      @save="handleCreate"
    />
  </div>
</template>

<script setup lang="ts">
import PageTitle from '~/components/PageTitle.vue'
import AnnouncementForm from '~/components/AnnouncementForm.vue'
import { useToast } from '~/composables/useToast'
import { useLoading } from '~/composables/useLoading'

definePageMeta({
  layout: 'default'
})

const router = useRouter()
const api = useApi()
const toast = useToast()
const { showLoading, hideLoading } = useLoading()

const handleCreate = async (payload: any) => {
  showLoading()
  try {
    await api.post('/api/admin/announcements', payload)
    toast.success('公告已建立')
    router.push('/admin/announcements')
  } catch (e: any) {
    toast.error(e.message || '建立失敗')
  } finally {
    hideLoading()
  }
}
</script>

<style scoped>
.admin-announcement-form-page {
  @apply max-w-lg mx-auto;
}
</style>
