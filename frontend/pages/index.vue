<template>
  <div class="space-list-page">
    <PageTitle title="我的空間" :show-back="false">
      <template #right>
        <BaseButton
          @click="router.push('/spaces/add-new')"
          variant="secondary"
          size="sm"
          class="gap-2"
        >
          <Icon icon="mdi:plus" class="text-lg" />
          <span>新增空間</span>
        </BaseButton>
      </template>
    </PageTitle>

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="spaces.length === 0" class="flex flex-col items-center justify-center py-20 bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed text-neutral-500">
      <Icon icon="mdi:view-grid-plus-outline" class="text-5xl mb-4 opacity-20" />
      <p class="text-sm">尚未建立任何管理空間</p>
      <BaseButton @click="router.push('/spaces/add-new')" class="mt-6">立即建立</BaseButton>
    </div>

    <div v-else class="flex flex-col gap-4">
      <SpaceListItem
        v-for="space in spaces"
        :key="space.id"
        :space="space"
        @click="router.push(`/spaces/${space.id}/ledger`)"
        @toggle-pin="handleTogglePin(space.id)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useSpace } from '~/composables/useSpace'
import SpaceListItem from '~/components/SpaceListItem.vue'
import PageTitle from '~/components/PageTitle.vue'
import BaseButton from '~/components/BaseButton.vue'

const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const { togglePin } = useSpace()

const spaces = ref<any[]>([])
const loading = ref(true)

const fetchSpaces = async () => {
  try {
    spaces.value = await api.get<any[]>('/api/spaces')
  } catch (e) {
    console.error('Failed to fetch spaces:', e)
  } finally {
    loading.value = false
  }
}

const handleTogglePin = async (id: string) => {
    try {
        await togglePin(id)
        // Refresh the list to show new order
        await fetchSpaces()
    } catch (e) {
        console.error('Pin failed', e)
    }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchSpaces()
})
</script>
