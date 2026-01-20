<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.push('/trips')" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold truncate max-w-48">{{ trip?.name || '載入中...' }}</h1>
      <button @click="showMenu = !showMenu" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors relative">
        <Icon icon="mdi:dots-vertical" class="text-2xl" />
        <div v-if="showMenu" class="absolute top-12 right-0 bg-neutral-800 rounded-xl border border-neutral-700 shadow-lg z-10 overflow-hidden min-w-32">
          <button @click="handleDelete" class="w-full px-4 py-3 text-left text-red-500 hover:bg-neutral-700 transition-colors border-0 bg-transparent cursor-pointer">刪除旅行</button>
        </div>
      </button>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <template v-else-if="trip">
      <!-- Trip Info Card -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <div class="flex items-center gap-3 mb-4">
          <div class="flex items-center justify-center w-14 h-14 rounded-xl bg-indigo-500/20">
            <Icon icon="mdi:airplane" class="text-3xl text-indigo-500" />
          </div>
          <div class="flex-1">
            <h2 class="text-lg font-bold">{{ trip.name }}</h2>
            <p class="text-sm text-neutral-400">{{ formatDateRange(trip.start_date, trip.end_date) }}</p>
          </div>
        </div>
        <p v-if="trip.description" class="text-neutral-400 text-sm mb-4">{{ trip.description }}</p>
        <div class="flex items-center gap-2 text-sm text-neutral-400">
          <Icon icon="mdi:currency-usd" class="text-lg" />
          <span>基準貨幣: {{ trip.base_currency }}</span>
        </div>
      </div>

      <!-- Members Section -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-base font-semibold">成員 ({{ trip.members?.length || 0 }})</h3>
          <button @click="showAddMember = true" class="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-indigo-500/20 text-indigo-500 text-sm border-0 cursor-pointer hover:bg-indigo-500/30 transition-colors">
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

    <!-- Add Member Modal -->
    <div v-if="showAddMember" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4" @click.self="showAddMember = false">
      <div class="bg-neutral-900 rounded-2xl p-5 w-full max-w-sm border border-neutral-800">
        <h3 class="text-lg font-bold mb-4">新增成員</h3>
        <input v-model="newMemberName" type="text" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500 placeholder-neutral-400 mb-4" placeholder="成員名稱" @keyup.enter="addMember" />
        <div class="flex gap-2">
          <button @click="showAddMember = false" class="flex-1 px-4 py-3 rounded-xl bg-neutral-800 text-white border-0 cursor-pointer hover:bg-neutral-700 transition-colors">取消</button>
          <button @click="addMember" class="flex-1 px-4 py-3 rounded-xl bg-indigo-500 text-white border-0 cursor-pointer hover:bg-indigo-600 transition-colors" :disabled="!newMemberName.trim()">新增</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const trip = ref<any>(null)
const loading = ref(true)
const showMenu = ref(false)
const showAddMember = ref(false)
const newMemberName = ref('')

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
  } catch (e) {
    console.error('Failed to fetch trip:', e)
    router.push('/trips')
  } finally {
    loading.value = false
  }
}

const addMember = async () => {
  if (!newMemberName.value.trim()) return
  try {
    await api.post(`/api/trips/${route.params.id}/members`, { name: newMemberName.value.trim() })
    newMemberName.value = ''
    showAddMember.value = false
    fetchTrip()
  } catch (e: any) {
    alert(e.message || '新增失敗')
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
