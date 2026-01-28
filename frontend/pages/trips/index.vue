<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <h1 class="text-2xl font-bold">我的旅行</h1>
      <NuxtLink to="/trips/create" class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline">
        <Icon icon="mdi:plus" /> 新增
      </NuxtLink>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <EmptyState
      v-else-if="trips.length === 0"
      icon="mdi:airplane-takeoff"
      title="開始你的旅程"
      description="建立一個旅行來記錄費用和比價"
      action-label="新增旅行"
      action-link="/trips/create"
    />

    <div v-else class="flex flex-col gap-3">
      <TripListItem
        v-for="trip in trips"
        :key="trip.id"
        :trip="trip"
        :cover-image="tripCovers[trip.id]"
        @click="router.push(`/trips/${trip.id}`)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import TripListItem from '~/components/TripListItem.vue'
import EmptyState from '~/components/EmptyState.vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const trips = ref<any[]>([])
const tripCovers = ref<Record<string, string>>({})
const { getImagesBatch, getImageUrl } = useImages()
const loading = ref(true)

const fetchTrips = async () => {
  try {
    trips.value = await api.get<any[]>('/api/trips')

    // Fetch covers
    if (trips.value.length > 0) {
        const ids = trips.value.map(t => t.id)
        const images = await getImagesBatch(ids, 'trip')
        const map: Record<string, string> = {}
        // Since sorted by order, the first one found for each entity is logically the cover (or first uploaded)
        // We iterate and assign if not exists (so we keep the first one we see if list is flattened? 
        // Backend returns one list. If Image A (sort 0) and Image B (sort 1) both in list for Trip X.
        // We want A.
        images.forEach(img => {
            if (!map[img.entity_id]) {
                map[img.entity_id] = img.file_path
            }
        })
        tripCovers.value = map
    }
  } catch (e) {
    console.error('Failed to fetch trips:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchTrips()
})
</script>
