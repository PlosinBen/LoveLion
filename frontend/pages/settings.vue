<template>
  <div class="settings-page min-h-screen bg-neutral-900 text-neutral-50">
    <SpaceHeader 
      title="個人設定" 
      :show-back="true"
      class="pt-0 px-2"
    />

    <div class="p-4 pt-0 flex flex-col gap-6">
      <!-- Profile Section -->
      <section class="flex flex-col gap-2">
        <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">帳戶資訊</label>
        <div class="bg-neutral-900 rounded-3xl border border-neutral-800 p-6 flex items-center justify-between shadow-sm">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 rounded-2xl bg-indigo-500/10 flex items-center justify-center text-indigo-500 border border-indigo-500/20">
              <Icon icon="mdi:account-outline" class="text-3xl" />
            </div>
            <div class="flex flex-col">
              <span class="text-lg font-black text-white tracking-tight">{{ user?.display_name || user?.username || '使用者' }}</span>
              <span class="text-xs text-neutral-500 font-medium">@{{ user?.username }}</span>
            </div>
          </div>
          <button @click="handleLogout" class="px-4 py-2 rounded-xl bg-red-500/10 text-red-500 text-xs font-black border-0 cursor-pointer hover:bg-red-500/20 transition-colors">
            登出
          </button>
        </div>
      </section>

      <!-- App Info -->
      <section class="flex flex-col gap-2">
        <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">關於</label>
        <div class="bg-neutral-900 rounded-3xl border border-neutral-800 p-6 flex flex-col gap-4 shadow-sm">
          <div class="flex items-center justify-between text-sm">
            <span class="text-neutral-400 font-medium">軟體版本</span>
            <span class="text-white font-black">v1.2.0-stable</span>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-neutral-400 font-medium">開發團隊</span>
            <span class="text-white font-black">Antigravity Design</span>
          </div>
        </div>
      </section>

      <p class="text-center text-[10px] text-neutral-700 font-black uppercase tracking-[0.2em] mt-10">
        LoveLion © 2026 Crafted with Passion
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'

definePageMeta({
  layout: 'default'
})

const router = useRouter()
const { user, logout: authLogout, initAuth, isAuthenticated } = useAuth()

const handleLogout = () => {
  authLogout()
  router.push('/login')
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
  }
})
</script>
