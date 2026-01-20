<template>
  <div class="trips-page">
    <header class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">我的旅行</h1>
      <NuxtLink to="/trips/new" class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline">
        <Icon icon="mdi:plus" /> 新增
      </NuxtLink>
    </header>

    <div class="text-center py-16 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
      <Icon icon="mdi:airplane-takeoff" class="text-6xl text-indigo-500 mb-4" />
      <h2 class="text-xl font-bold mb-2">即將推出</h2>
      <p class="text-neutral-400">旅行記帳功能正在開發中...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
  }
})
</script>
