<template>
  <div class="flex flex-col gap-6">
    <!-- Immersive Header Area -->
    <div class="relative w-full h-64 rounded-2xl overflow-hidden shrink-0 bg-neutral-800">
       <!-- Background: Image or Fallback -->
       <img v-if="coverImage" :src="getImageUrl(coverImage)" class="absolute inset-0 w-full h-full object-cover" alt="Trip Cover" />
       
       <!-- Gradient Overlays -->
       <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-transparent to-transparent"></div> <!-- Top shadow for buttons -->
       <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div> <!-- Bottom shadow for text -->

       <!-- Header Controls (Top) -->
       <header class="relative z-10 flex justify-between items-center p-3">
          <button @click="router.push('/trips')" class="flex justify-center items-center w-10 h-10 rounded-xl bg-black/20 text-white backdrop-blur-md border-0 cursor-pointer hover:bg-black/40 transition-colors">
             <Icon icon="mdi:arrow-left" class="text-2xl" />
          </button>

          <button @click="showMenu = !showMenu" class="flex justify-center items-center w-10 h-10 rounded-xl bg-black/20 text-white backdrop-blur-md border-0 cursor-pointer hover:bg-black/40 transition-colors relative">
             <Icon icon="mdi:dots-vertical" class="text-2xl" />
             <div v-if="showMenu" class="absolute top-12 right-0 bg-neutral-800 rounded-xl border border-neutral-700 shadow-lg z-20 overflow-hidden min-w-32 py-1">
               <NuxtLink :to="`/trips/${trip?.id}/edit`" class="block w-full px-4 py-3 text-left text-white hover:bg-neutral-700 transition-colors no-underline">編輯旅行</NuxtLink>
               <button @click="handleDelete" class="w-full px-4 py-3 text-left text-red-500 hover:bg-neutral-700 transition-colors border-0 bg-transparent cursor-pointer">刪除旅行</button>
             </div>
          </button>
       </header>

       <!-- Title & Date (Bottom Overlay) -->
       <div class="absolute bottom-0 left-0 w-full p-5 z-10">
          <h1 class="text-2xl font-bold text-white mb-1 shadow-sm">{{ trip?.name || '載入中...' }}</h1>
          <div class="flex items-center gap-2 text-neutral-300 text-sm">
             <Icon icon="mdi:calendar-range" class="text-indigo-400" />
             <span>{{ formatDateRange(trip?.start_date, trip?.end_date) }}</span>
          </div>
       </div>
    </div>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <template v-else-if="trip">
      <!-- Simplified Info Bar -->
      <div class="flex flex-col gap-4 px-1">
         <p v-if="trip.description" class="text-neutral-400 text-sm leading-relaxed">{{ trip.description }}</p>
         
         <div class="flex items-center gap-2">
             <div class="px-3 py-1.5 rounded-lg bg-neutral-800 border border-neutral-700 text-xs text-neutral-300 flex items-center gap-1.5">
                 <Icon icon="mdi:currency-usd" class="text-indigo-400" />
                 <span>{{ trip.base_currency }}</span>
             </div>
             <!-- Add more stats here later if needed -->
         </div>
      </div>

      <!-- Members Section -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-base font-semibold">成員 ({{ trip.members?.length || 0 }})</h3>
          <button @click="router.push(`/trips/${trip.id}/members/add`)" class="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-indigo-500/20 text-indigo-500 text-sm border-0 cursor-pointer hover:bg-indigo-500/30 transition-colors">
            <Icon icon="mdi:plus" /> 新增
          </button>
        </div>
        <div class="flex flex-wrap gap-2">
          <div v-for="member in trip.members" :key="member.id" class="flex items-center gap-2 px-3 py-2 rounded-full bg-neutral-800 text-sm">
            <Icon icon="mdi:account" class="text-indigo-500" />
            <span>{{ member.name }}</span>
            <span v-if="member.is_owner" class="text-xs text-indigo-500">(主辦)</span>
            <button v-if="!member.is_owner" @click="removeMember(member.id)" class="ml-1 text-neutral-500 hover:text-red-500 transition-colors border-0 bg-transparent cursor-pointer p-0">
              <Icon icon="mdi:close" class="text-sm" />
            </button>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-2 gap-3">
        <NuxtLink :to="`/trips/${trip.id}/ledger`" class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500">
          <Icon icon="mdi:notebook-edit-outline" class="text-3xl text-indigo-500" />
          <span class="text-sm">公費記帳</span>
        </NuxtLink>
        <NuxtLink :to="`/trips/${trip.id}/stores`" class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500">
          <Icon icon="mdi:store" class="text-3xl text-indigo-500" />
          <span class="text-sm">比價商店</span>
        </NuxtLink>
        <NuxtLink :to="`/trips/${trip.id}/compare`" class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500">
          <Icon icon="mdi:chart-bar" class="text-3xl text-indigo-500" />
          <span class="text-sm">價格比較</span>
        </NuxtLink>
      </div>
    </template>


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const trip = ref<any>(null)
const loading = ref(true)
const showMenu = ref(false)
const coverImage = ref<string | null>(null)
const { getImages, getImageUrl } = useImages()


const formatDateRange = (start: string | null, end: string | null) => {
  if (!start && !end) return '日期未設定'
  const opts: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' }
  const startStr = start ? new Date(start).toLocaleDateString('zh-TW', opts) : '?'
  const endStr = end ? new Date(end).toLocaleDateString('zh-TW', opts) : '?'
  return `${startStr} - ${endStr}`
}

const fetchTrip = async () => {
  try {
    trip.value = await api.get<any>(`/api/trips/${route.params.id}`)
    
    // Fetch cover image
    const images = await getImages(trip.value.id, 'trip')
    if (images.length > 0) {
        coverImage.value = images[0].file_path
    }
  } catch (e) {
    console.error('Failed to fetch trip:', e)
    router.push('/trips')
  } finally {
    loading.value = false
  }
}



const removeMember = async (memberId: string) => {
  if (!confirm('確定要移除此成員？')) return
  try {
    await api.del(`/api/trips/${route.params.id}/members/${memberId}`)
    fetchTrip()
  } catch (e: any) {
    alert(e.message || '移除失敗')
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除此旅行？')) return
  try {
    await api.del(`/api/trips/${route.params.id}`)
    router.push('/trips')
  } catch (e: any) {
    alert(e.message || '刪除失敗')
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchTrip()
})
</script>
