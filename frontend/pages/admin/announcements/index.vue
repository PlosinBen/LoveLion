<template>
  <div>
    <PageTitle
      title="公告管理"
      :breadcrumbs="[{ label: '設定', to: '/settings' }]"
    >
      <template #right>
        <NuxtLink
          to="/admin/announcements/new"
          class="w-10 h-10 rounded-full bg-indigo-500 text-white flex items-center justify-center hover:bg-indigo-400 transition-colors no-underline active:scale-95"
        >
          <Icon icon="mdi:plus" class="text-xl" />
        </NuxtLink>
      </template>
    </PageTitle>

    <div v-if="loading" class="p-20 flex justify-center items-center text-neutral-500">
      <Icon icon="mdi:loading" class="text-4xl animate-spin" />
    </div>

    <div v-else-if="announcements.length === 0" class="py-20 flex flex-col items-center justify-center text-neutral-600">
      <Icon icon="mdi:bullhorn-outline" class="text-5xl mb-3 opacity-20" />
      <span class="text-sm font-bold">尚未建立公告</span>
    </div>

    <div v-else class="flex flex-col gap-3 pb-20">
      <NuxtLink
        v-for="item in announcements"
        :key="item.id"
        :to="`/admin/announcements/${item.id}`"
        class="no-underline"
      >
        <BaseCard padding="p-4" class="shadow-sm hover:bg-neutral-800 transition-colors cursor-pointer active:scale-95">
          <div class="flex items-center justify-between">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <h3 class="text-sm font-bold text-neutral-100 truncate">{{ item.title }}</h3>
                <span
                  class="text-xs font-bold px-2 py-0.5 rounded-full shrink-0"
                  :class="item.status === 'published'
                    ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                    : 'bg-neutral-700/50 text-neutral-400 border border-neutral-600'"
                >
                  {{ item.status === 'published' ? '已發布' : '草稿' }}
                </span>
              </div>
              <div class="flex items-center gap-3 mt-1.5">
                <span class="text-xs text-neutral-500">{{ formatDate(item.created_at) }}</span>
                <span v-if="item.broadcast_start" class="text-xs text-indigo-400 flex items-center gap-1">
                  <Icon icon="mdi:broadcast" class="text-xs" />
                  廣播中
                </span>
              </div>
            </div>
            <Icon icon="mdi:chevron-right" class="text-neutral-600 shrink-0" />
          </div>
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
  return date.toLocaleDateString('zh-TW', { year: 'numeric', month: 'short', day: 'numeric' })
}

onMounted(async () => {
  try {
    announcements.value = await api.get<Announcement[]>('/api/admin/announcements')
  } catch (e: any) {
    if (e.message?.includes('403') || e.message?.includes('Forbidden')) {
      navigateTo('/settings')
    }
    console.error('Failed to fetch announcements:', e)
  } finally {
    loading.value = false
  }
})
</script>
