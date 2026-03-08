<template>
  <div class="settings-page">
    <PageTitle
      title="個人設定"
      :show-back="true"
      :breadcrumbs="[{ label: '空間列表', to: '/' }]"
    />

    <div class="flex flex-col gap-6">
      <!-- Profile Section -->
      <section class="flex flex-col gap-2">
        <label class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">帳戶資訊</label>
        <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-6 flex items-center justify-between shadow-sm">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-500 border border-neutral-700">
              <Icon icon="mdi:account-outline" class="text-3xl" />
            </div>
            <div class="flex flex-col">
              <span class="text-lg font-bold text-white tracking-tight">{{ user?.display_name || user?.username || '使用者' }}</span>
              <span class="text-xs text-neutral-500 font-medium">@{{ user?.username }}</span>
            </div>
          </div>
          <BaseButton @click="handleLogout" variant="danger" size="sm">
            登出
          </BaseButton>
        </div>
      </section>

      <!-- App Info -->
      <section class="flex flex-col gap-2">
        <label class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">關於</label>
        <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-6 flex flex-col gap-4 shadow-sm">
          <div class="flex items-center justify-between text-sm">
            <span class="text-neutral-400 font-medium">軟體版本</span>
            <span class="text-white font-bold">v1.2.0-stable</span>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-neutral-400 font-medium">開發團隊</span>
            <span class="text-white font-bold">Antigravity Design</span>
          </div>
        </div>
      </section>

      <p class="text-center text-xs text-neutral-700 font-bold uppercase tracking-widest mt-10">
        LoveLion © 2026 Crafted with Passion
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import PageTitle from '~/components/PageTitle.vue'
import BaseButton from '~/components/BaseButton.vue'

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
