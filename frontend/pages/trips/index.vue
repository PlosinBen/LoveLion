<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <h1 class="text-2xl font-bold">我的旅行</h1>
      <NuxtLink to="/trips/create" class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline">
        <Icon icon="mdi:plus" /> 新增
      </NuxtLink>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <div v-else-if="trips.length === 0" class="text-center py-16 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
      <Icon icon="mdi:airplane-takeoff" class="text-6xl text-indigo-500 mb-4" />
      <h2 class="text-xl font-bold mb-2">開始你的旅程</h2>
      <p class="text-neutral-400 mb-5">建立一個旅行來記錄費用和比價</p>
      <NuxtLink to="/trips/create" class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline">
        新增旅行
      </NuxtLink>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div
        v-for="trip in trips"
        :key="trip.id"
        class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 cursor-pointer transition-all duration-200 hover:-translate-y-0.5 hover:border-indigo-500"
        @click="router.push(`/trips/${trip.id}`)"
      >
        <div class="flex items-center gap-3">
          <div class="relative flex items-center justify-center w-12 h-12 rounded-xl bg-indigo-500/20 overflow-hidden">
             <img v-if="tripCovers[trip.id]" :src="getImageUrl(tripCovers[trip.id] || '')" class="w-full h-full object-cover" />
             <Icon v-else icon="mdi:airplane" class="text-2xl text-indigo-500" />
          </div>
          <div class="flex-1">
            <h3 class="text-base font-bold mb-1">{{ trip.name }}</h3>
            <p class="text-xs text-neutral-400">
              {{ formatDateRange(trip.start_date, trip.end_date) }}
              <span v-if="trip.members?.length" class="ml-2">• {{ trip.members.length }} 位成員</span>
            </p>
          </div>
          <Icon icon="mdi:chevron-right" class="text-xl text-neutral-500" />
        </div>
        <p v-if="trip.description" class="mt-3 pt-3 border-t border-neutral-800 text-sm text-neutral-400 line-clamp-2">
          {{ trip.description }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
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

const formatDateRange = (start: string | null, end: string | null) => {
  if (!start && !end) return '日期未設定'
  const opts: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' }
  const startStr = start ? new Date(start).toLocaleDateString('zh-TW', opts) : '?'
  const endStr = end ? new Date(end).toLocaleDateString('zh-TW', opts) : '?'
  return `${startStr} - ${endStr}`
}

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
