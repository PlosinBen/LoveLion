<template>
  <div class="flex flex-col gap-10">
    <!-- Status (Optional, minimized) -->
    <StatusCard
      v-if="backendStatus !== 'online'"
      :status="backendStatus"
      :message="statusMessage"
    />

    <!-- Pinned Spaces Section -->
    <section v-if="pinnedSpaces.length > 0">
      <div class="flex items-center justify-between mb-4 px-1">
        <h2 class="text-xl font-black text-white tracking-tight flex items-center gap-2">
          <Icon icon="mdi:pin" class="text-indigo-500" />
          ?阡?蝛粹?
        </h2>
      </div>
      
      <div class="flex flex-col gap-3">
        <SpaceListItem 
          v-for="space in pinnedSpaces" 
          :key="space.id"
          :space="space"
          @click="router.push(`/spaces/${space.id}`)"
          @toggle-pin="togglePin(space.id)"
        />
      </div>
    </section>

    <!-- All Spaces Section -->
    <section>
      <div class="flex items-center justify-between mb-4 px-1">
        <h2 class="text-xl font-black text-white tracking-tight">??征??/h2>
        <NuxtLink to="/spaces/add-new" class="text-sm font-bold text-indigo-400 no-underline hover:text-indigo-300">
          + ?啣?
        </NuxtLink>
      </div>

      <div v-if="loading" class="flex flex-col gap-3">
        <div v-for="i in 3" :key="i" class="h-20 bg-neutral-800 animate-pulse rounded-3xl"></div>
      </div>

      <div v-else-if="allSpaces.length === 0" class="py-12 flex flex-col items-center justify-center bg-neutral-900 rounded-3xl border border-neutral-800 border-dashed">
        <Icon icon="mdi:folder-open-outline" class="text-5xl text-neutral-700 mb-4" />
        <p class="text-neutral-500 text-sm font-medium">?桀????遙雿征??/p>
        <NuxtLink to="/spaces/add-new" class="mt-4 px-6 py-2 bg-indigo-500 text-white rounded-full font-bold text-sm no-underline">撱箇?蝚砌??征??/NuxtLink>
      </div>

      <div v-else class="flex flex-col gap-3">
        <SpaceListItem 
          v-for="space in otherSpaces" 
          :key="space.id"
          :space="space"
          @click="router.push(`/spaces/${space.id}`)"
          @toggle-pin="togglePin(space.id)"
        />
      </div>
    </section>

    <!-- App Info / Settings Link -->
    <div class="mt-4 pt-10 border-t border-neutral-800 flex flex-col gap-4">
       <button @click="router.push('/settings')" class="flex items-center justify-between p-6 rounded-3xl bg-neutral-900 border border-neutral-800 text-neutral-400 hover:text-white transition-all border-0 cursor-pointer active:scale-95">
          <div class="flex items-center gap-4">
            <Icon icon="mdi:cog-outline" class="text-2xl" />
            <span class="font-bold">蝟餌絞閮剖?</span>
          </div>
          <Icon icon="mdi:chevron-right" class="text-xl" />
       </button>
       
       <button @click="handleLogout" class="flex items-center justify-center gap-2 p-4 text-neutral-600 hover:text-red-500 text-sm font-bold border-0 bg-transparent cursor-pointer transition-colors">
          <Icon icon="mdi:logout" /> ?餃撣唾?
       </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import StatusCard from '~/components/StatusCard.vue'
import SpaceListItem from '~/components/SpaceListItem.vue'
import { useAuth } from '~/composables/useAuth'
import { useSpace } from '~/composables/useSpace'

definePageMeta({
  layout: 'default'
})

const router = useRouter()
const { isAuthenticated, initAuth, logout } = useAuth()
const { allSpaces, pinnedSpaces, otherSpaces, loading, fetchSpaces, togglePin } = useSpace()

const backendStatus = ref<'online' | 'offline' | 'checking'>('checking')

const statusMessage = computed(() => {
  switch (backendStatus.value) {
    case 'online': return '敺垢??甇?虜??'
    case 'offline': return '?⊥???敺垢??'
    default: return '甇?瑼Ｘ???...'
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
/* Unified simple style */
</style>
