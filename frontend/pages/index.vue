<template>
  <div class="flex flex-col gap-8 pb-24">
    <!-- Status (Optional, minimized) -->
    <StatusCard
      v-if="backendStatus !== 'online'"
      :status="backendStatus"
      :message="statusMessage"
    />

    <!-- Pinned Spaces Section -->
    <section v-if="pinnedSpaces.length > 0">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-bold text-white tracking-tight flex items-center gap-2">
          <Icon icon="mdi:pin" class="text-indigo-500" />
          釘選空間
        </h2>
      </div>
      
      <div class="flex overflow-x-auto gap-4 pb-4 -mx-4 px-4 scrollbar-hide">
        <div 
          v-for="space in pinnedSpaces" 
          :key="space.id"
          class="min-w-[280px] h-44 relative rounded-3xl overflow-hidden shadow-xl cursor-pointer group shrink-0"
          @click="router.push(`/spaces/${space.id}`)"
        >
          <!-- Background -->
          <div class="absolute inset-0 bg-neutral-800">
            <img v-if="space.cover_image" :src="getImageUrl(space.cover_image)" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" />
            <div v-else class="w-full h-full bg-gradient-to-br from-indigo-600 to-purple-700 opacity-80"></div>
            <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/30 to-transparent"></div>
          </div>
          
          <!-- Content -->
          <div class="absolute inset-0 p-6 flex flex-col justify-end">
            <div class="flex flex-col">
              <span class="text-white/60 text-[10px] font-bold uppercase tracking-widest mb-1">{{ space.type === 'trip' ? '旅行專案' : '日常記帳' }}</span>
              <h3 class="text-xl font-bold text-white mb-2">{{ space.name }}</h3>
              <div class="flex items-center gap-3">
                <span class="text-white/80 text-xs flex items-center gap-1">
                  <Icon icon="mdi:currency-usd" class="text-indigo-400" />
                  {{ space.base_currency }}
                </span>
                <span v-if="space.start_date" class="text-white/80 text-xs flex items-center gap-1">
                  <Icon icon="mdi:calendar" class="text-indigo-400" />
                  {{ formatDate(space.start_date) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Unpin Toggle -->
          <button @click.stop="togglePin(space.id)" class="absolute top-4 right-4 w-8 h-8 rounded-full bg-black/20 backdrop-blur-md flex items-center justify-center text-white border-0 cursor-pointer hover:bg-white/20 transition-colors">
            <Icon icon="mdi:pin-off" />
          </button>
        </div>
      </div>
    </section>

    <!-- All Spaces Section -->
    <section>
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-bold text-white tracking-tight">所有空間</h2>
        <NuxtLink to="/spaces/add-new" class="text-sm font-bold text-indigo-400 no-underline hover:text-indigo-300">
          + 新增
        </NuxtLink>
      </div>

      <div v-if="loading" class="flex flex-col gap-3">
        <div v-for="i in 3" :key="i" class="h-24 bg-neutral-800 animate-pulse rounded-3xl"></div>
      </div>

      <div v-else-if="allSpaces.length === 0" class="py-12 flex flex-col items-center justify-center bg-neutral-900 rounded-3xl border border-neutral-800 border-dashed">
        <Icon icon="mdi:folder-open-outline" class="text-5xl text-neutral-700 mb-4" />
        <p class="text-neutral-500 text-sm">目前還沒有任何空間</p>
        <NuxtLink to="/spaces/add-new" class="mt-4 px-6 py-2 bg-indigo-500 text-white rounded-full font-bold text-sm no-underline">建立第一個空間</NuxtLink>
      </div>

      <div v-else class="flex flex-col gap-3">
        <div 
          v-for="space in otherSpaces" 
          :key="space.id"
          class="bg-neutral-900 border border-neutral-800/50 p-5 rounded-3xl flex items-center justify-between hover:bg-neutral-800 transition-colors cursor-pointer group"
          @click="router.push(`/spaces/${space.id}`)"
        >
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl flex items-center justify-center text-2xl" :class="space.type === 'trip' ? 'bg-blue-500/10 text-blue-400' : 'bg-green-500/10 text-green-400'">
              <Icon :icon="space.type === 'trip' ? 'mdi:airplane' : 'mdi:wallet-outline'" />
            </div>
            <div class="flex flex-col">
              <h3 class="font-bold text-white group-hover:text-indigo-400 transition-colors">{{ space.name }}</h3>
              <span class="text-xs text-neutral-500">{{ space.type === 'trip' ? '旅行' : '個人' }} · {{ space.base_currency }}</span>
            </div>
          </div>
          
          <button @click.stop="togglePin(space.id)" class="w-8 h-8 rounded-full flex items-center justify-center text-neutral-600 hover:text-indigo-400 transition-colors border-0 bg-transparent cursor-pointer">
            <Icon icon="mdi:pin-outline" class="text-xl" />
          </button>
        </div>
      </div>
    </section>

    <!-- App Info / Settings Link -->
    <div class="mt-4 pt-8 border-t border-neutral-800 flex flex-col gap-4">
       <button @click="router.push('/settings')" class="flex items-center justify-between p-5 rounded-3xl bg-neutral-900 border border-neutral-800 text-neutral-400 hover:text-white transition-colors border-0 cursor-pointer">
          <div class="flex items-center gap-3">
            <Icon icon="mdi:cog-outline" class="text-xl" />
            <span class="font-medium">系統設定</span>
          </div>
          <Icon icon="mdi:chevron-right" />
       </button>
       
       <button @click="handleLogout" class="flex items-center justify-center gap-2 p-4 text-red-500/60 hover:text-red-500 text-sm font-medium border-0 bg-transparent cursor-pointer transition-colors">
          <Icon icon="mdi:logout" /> 登出帳號
       </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import StatusCard from '~/components/StatusCard.vue'
import { useAuth } from '~/composables/useAuth'
import { useSpace } from '~/composables/useSpace'
import { useImages } from '~/composables/useImages'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const { isAuthenticated, initAuth, logout } = useAuth()
const { allSpaces, pinnedSpaces, otherSpaces, loading, fetchSpaces, togglePin } = useSpace()
const { getImageUrl } = useImages()

const backendStatus = ref<'online' | 'offline' | 'checking'>('checking')

const statusMessage = computed(() => {
  switch (backendStatus.value) {
    case 'online': return '後端服務正常運行'
    case 'offline': return '無法連接後端服務'
    default: return '正在檢查連線...'
  }
})

const checkBackend = async () => {
  try {
    const response = await fetch('/health')
    if (response.ok) backendStatus.value = 'online'
    else backendStatus.value = 'offline'
  } catch {
    backendStatus.value = 'offline'
  }
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-TW', { month: 'short', day: 'numeric' })
}

const handleLogout = () => {
  logout()
  router.push('/login')
}

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  checkBackend()
  fetchSpaces()
})
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
