<template>
  <div class="space-list-page">
    <PageTitle title="我的空間" :show-back="false" />

    <div v-if="localLoading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <template v-else>
      <div v-if="allSpaces.length === 0" class="flex flex-col items-center justify-center py-20 bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed text-neutral-500">
        <Icon icon="mdi:view-grid-plus-outline" class="text-5xl mb-4 opacity-20" />
        <p class="text-sm">尚未建立任何管理空間</p>
        <NuxtLink to="/spaces/add-new" class="mt-6 inline-flex justify-center items-center px-4 py-2.5 text-sm rounded bg-indigo-500 text-white font-bold hover:bg-indigo-600 no-underline shadow-lg transition-all active:scale-95">立即建立</NuxtLink>
      </div>

      <div v-else class="flex flex-col gap-4">
        <SpaceListItem
          v-for="space in sortedSpaces"
          :key="space.id"
          :space="space"
          @click="router.push(`/spaces/${space.id}/stats`)"
          @toggle-pin="handleTogglePin(space.id)"
        />
      </div>    </template>

    <NuxtLink
      to="/spaces/add-new"
      class="fixed bottom-20 right-4 z-40 flex items-center justify-center w-14 h-14 rounded-full bg-indigo-500 text-white shadow-lg active:scale-95 transition-transform no-underline"
    >
      <Icon icon="mdi:plus" class="text-2xl" />
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useSpace } from '~/composables/useSpace'
import SpaceListItem from '~/components/SpaceListItem.vue'
import PageTitle from '~/components/PageTitle.vue'


const router = useRouter()
const { allSpaces, fetchSpaces, togglePin } = useSpace()

// Use local loading state to prevent hydration mismatch
const localLoading = ref(true)

const sortedSpaces = computed(() => {
  return [...allSpaces.value].sort((a, b) => {
    // Pinned first
    if (a.is_pinned && !b.is_pinned) return -1
    if (!a.is_pinned && b.is_pinned) return 1
    // Then by name
    return a.name.localeCompare(b.name)
  })
})

const handleTogglePin = async (id: string) => {
    try {
        await togglePin(id)
    } catch (e) {
        console.error('Pin failed', e)
    }
}

onMounted(async () => {
  try {
    await fetchSpaces(true)
  } finally {
    localLoading.value = false
  }
})
</script>
