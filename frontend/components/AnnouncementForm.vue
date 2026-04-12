<template>
  <form @submit.prevent="handleSubmit" class="flex flex-col gap-4 pb-20">
    <!-- AI Generate -->
    <BaseCard v-if="aiAvailable" padding="p-4" class="shadow-sm border border-indigo-500/20">
      <div class="flex flex-col gap-3">
        <label class="text-xs font-bold text-indigo-400 uppercase tracking-wider flex items-center gap-1.5">
          <Icon icon="mdi:robot-outline" class="text-base" />
          AI 生成公告
        </label>
        <div class="flex gap-2">
          <BaseInput
            v-model="aiPrompt"
            placeholder="描述你想公告的事項..."
            class="flex-1"
          />
          <BaseButton
            type="button"
            @click="handleGenerate"
            :disabled="!aiPrompt.trim() || generating"
            class="shrink-0"
          >
            <Icon v-if="generating" icon="mdi:loading" class="animate-spin" />
            <span v-else>生成</span>
          </BaseButton>
        </div>
      </div>
    </BaseCard>

    <!-- Title -->
    <BaseInput
      v-model="form.title"
      label="標題"
      placeholder="公告標題"
      required
    />

    <!-- Content -->
    <BaseTextarea
      v-model="form.content"
      label="內容（Markdown）"
      placeholder="公告內容，支援 Markdown 格式"
      rows="10"
    />

    <!-- Preview -->
    <BaseCollapsible v-if="form.content" title="內容預覽">
      <div class="prose prose-invert prose-sm max-w-none p-4" v-html="renderedContent"></div>
    </BaseCollapsible>

    <!-- Status -->
    <div class="flex flex-col gap-1.5">
      <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">狀態</label>
      <BaseSegmentControl
        v-model="form.status"
        :options="statusOptions"
      />
    </div>

    <!-- Broadcast -->
    <BaseCard padding="p-4" class="shadow-sm flex flex-col gap-3">
      <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider">廣播設定</label>
      <p class="text-xs text-neutral-600">設定廣播期間後，公告會顯示在所有頁面的頂部橫幅。留空表示不廣播。</p>
      <div class="flex flex-col gap-3">
        <BaseInput
          v-model="form.broadcastStart"
          label="開始時間"
          type="datetime-local"
        />
        <BaseInput
          v-model="form.broadcastEnd"
          label="結束時間"
          type="datetime-local"
        />
      </div>
    </BaseCard>

    <!-- Submit -->
    <BaseButton type="submit" class="mt-2">
      {{ initialData ? '更新公告' : '建立公告' }}
    </BaseButton>
  </form>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { marked } from 'marked'
import BaseInput from '~/components/BaseInput.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'
import BaseCard from '~/components/BaseCard.vue'
import BaseCollapsible from '~/components/BaseCollapsible.vue'
import BaseSegmentControl from '~/components/BaseSegmentControl.vue'
import { useToast } from '~/composables/useToast'
import type { Announcement } from '~/types'

const props = defineProps<{
  initialData: Announcement | null
}>()

const emit = defineEmits<{
  save: [payload: any]
}>()

const api = useApi()
const toast = useToast()

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '已發布', value: 'published' },
]

const form = reactive({
  title: '',
  content: '',
  status: 'draft',
  broadcastStart: '',
  broadcastEnd: '',
})

const aiPrompt = ref('')
const generating = ref(false)
const aiAvailable = ref(false)

const renderedContent = computed(() => {
  return form.content ? marked(form.content) : ''
})

const toLocalDatetime = (isoStr: string | null) => {
  if (!isoStr) return ''
  const d = new Date(isoStr)
  const offset = d.getTimezoneOffset()
  const local = new Date(d.getTime() - offset * 60 * 1000)
  return local.toISOString().slice(0, 16)
}

const toISO = (localDatetime: string) => {
  if (!localDatetime) return null
  return new Date(localDatetime).toISOString()
}

const handleGenerate = async () => {
  generating.value = true
  try {
    const result = await api.post<{ title: string; content: string }>('/api/admin/announcements/generate', {
      description: aiPrompt.value,
    })
    form.title = result.title
    form.content = result.content
    toast.success('AI 已生成公告內容')
  } catch (e: any) {
    toast.error(e.message || 'AI 生成失敗')
  } finally {
    generating.value = false
  }
}

const handleSubmit = () => {
  if (!form.title.trim()) {
    toast.error('請輸入標題')
    return
  }

  emit('save', {
    title: form.title,
    content: form.content,
    status: form.status,
    broadcast_start: toISO(form.broadcastStart),
    broadcast_end: toISO(form.broadcastEnd),
  })
}

const checkAI = async () => {
  try {
    const config = await api.get<{ ai_available: boolean }>('/api/admin/announcements/config')
    aiAvailable.value = config.ai_available
  } catch {
    aiAvailable.value = false
  }
}

onMounted(() => {
  if (props.initialData) {
    form.title = props.initialData.title
    form.content = props.initialData.content
    form.status = props.initialData.status
    form.broadcastStart = toLocalDatetime(props.initialData.broadcast_start)
    form.broadcastEnd = toLocalDatetime(props.initialData.broadcast_end)
  }
  checkAI()
})
</script>
