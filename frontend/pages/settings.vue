<template>
  <div class="settings-page">
    <header class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">設定</h1>
    </header>

    <div class="flex flex-col gap-4">
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <h2 class="text-lg font-semibold mb-3">帳號</h2>
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 rounded-full bg-indigo-500/20 flex items-center justify-center text-indigo-500 font-bold text-xl">
            {{ user?.name?.charAt(0) || 'U' }}
          </div>
          <div>
            <div class="font-medium">{{ user?.name || 'User' }}</div>
            <div class="text-sm text-neutral-400">{{ user?.email || 'user@example.com' }}</div>
          </div>
        </div>
        <button @click="logout" class="w-full mt-4 py-2.5 rounded-xl border border-red-500/20 text-red-500 hover:bg-red-500/10 transition-colors cursor-pointer bg-transparent">
          登出
        </button>
      </div>

      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <h2 class="text-lg font-semibold mb-3">關於</h2>
        <p class="text-neutral-400 text-sm">LoveLion v1.0.0</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const { user, logout: authLogout, initAuth, isAuthenticated } = useAuth()

const logout = () => {
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
