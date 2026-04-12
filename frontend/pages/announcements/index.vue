<template>
  <div class="announcements-page">
    <PageTitle
      title="公告"
      :breadcrumbs="[{ label: '設定', to: '/settings' }]"
    />

    <div v-if="loading" class="p-20 flex justify-center items-center text-neutral-500">
      <Icon icon="mdi:loading" class="text-4xl animate-spin" />
    </div>

    <div v-else-if="announcements.length === 0" class="py-20 flex flex-col items-center justify-center text-neutral-600">
      <Icon icon="mdi:bullhorn-outline" class="text-5xl mb-3 opacity-20" />
      <span class="text-sm font-bold">目前沒有公告</span>
    </div>

    <div v-else class="flex flex-col gap-3 pb-20">
      <NuxtLink
        v-for="item in announcements"
        :key="item.id"
        :to="`/announcements/${item.id}`"
        class="no-underline"
      >
        <BaseCard padding="p-5" class="shadow-sm hover:bg-neutral-800 transition-colors cursor-pointer active:scale-95">
          <h3 class="text-sm font-bold text-neutral-100">{{ item.title }}</h3>
          <p class="text-xs text-neutral-500 mt-1.5">{{ formatDate(item.created_at) }}</p>
        </BaseCard>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import PageTitle from '~/components/PageTitle.vue'
import BaseCard from '~/components/BaseCard.vue'
import type { Announcement } from '~/types'

definePageMeta({
  layout: 'default'
})

const api = useApi()
const announcements = ref<Announcement[]>([])
const loading = ref(true)

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(async () => {
  try {
    announcements.value = await api.get<Announcement[]>('/api/announcements')
  } catch (e) {
    console.error('Failed to fetch announcements:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.announcements-page {
  @apply max-w-lg mx-auto;
}
</style>
