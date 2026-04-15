<template>
  <div>
    <PageTitle
      :title="announcement?.title || '公告'"
      back-to="/announcements"
      :breadcrumbs="[{ label: '設定', to: '/settings' }, { label: '公告', to: '/announcements' }]"
    />

    <div v-if="loading" class="p-20 flex justify-center items-center text-neutral-500">
      <Icon icon="mdi:loading" class="text-4xl animate-spin" />
    </div>

    <div v-else-if="!announcement" class="py-20 flex flex-col items-center justify-center text-neutral-600">
      <Icon icon="mdi:alert-circle-outline" class="text-5xl mb-3 opacity-20" />
      <span class="text-sm font-bold">找不到此公告</span>
    </div>

    <div v-else class="flex flex-col gap-6 pb-20">
      <BaseCard padding="p-6" class="shadow-sm">
        <div class="flex flex-col gap-4">
          <div>
            <h2 class="text-lg font-bold text-white">{{ announcement.title }}</h2>
            <p class="text-xs text-neutral-500 mt-1">{{ formatDate(announcement.created_at) }}</p>
          </div>
          <div class="prose prose-invert prose-sm max-w-none" v-html="renderedContent"></div>
        </div>
      </BaseCard>

      <NuxtLink
        to="/announcements"
        class="text-center text-sm text-indigo-400 hover:text-indigo-300 transition-colors no-underline font-bold"
      >
        查看所有公告
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { marked } from 'marked'
import PageTitle from '~/components/PageTitle.vue'
import BaseCard from '~/components/BaseCard.vue'
import type { Announcement } from '~/types'

definePageMeta({
  layout: 'default'
})

const route = useRoute()
const api = useApi()
const announcement = ref<Announcement | null>(null)
const loading = ref(true)

const renderedContent = computed(() => {
  if (!announcement.value?.content) return ''
  return marked(announcement.value.content)
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(async () => {
  try {
    announcement.value = await api.get<Announcement>(`/api/announcements/${route.params.id}`)
  } catch (e) {
    console.error('Failed to fetch announcement:', e)
  } finally {
    loading.value = false
  }
})
</script>
