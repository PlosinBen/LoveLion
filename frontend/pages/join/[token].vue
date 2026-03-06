<template>
  <div class="join-ledger-page min-h-screen flex flex-col items-center justify-center p-6 bg-neutral-950">
    <div v-if="loading" class="text-neutral-400 animate-pulse">
      正在讀取邀請資訊...
    </div>

    <div v-else-if="error" class="bg-neutral-900 p-8 rounded-3xl border border-neutral-800 text-center max-w-sm w-full">
      <Icon icon="mdi:alert-circle-outline" class="text-6xl text-red-500 mb-4 mx-auto" />
      <h1 class="text-xl font-bold mb-2">邀請失效</h1>
      <p class="text-neutral-400 mb-6">{{ error }}</p>
      <button @click="router.push('/')" class="w-full py-3 rounded-xl bg-neutral-800 font-semibold hover:bg-neutral-700 transition-colors">
        回首頁
      </button>
    </div>

    <div v-else-if="inviteInfo" class="bg-neutral-900 p-8 rounded-3xl border border-neutral-800 text-center max-w-sm w-full">
      <div class="w-20 h-20 bg-indigo-500/10 rounded-full flex items-center justify-center mx-auto mb-6">
        <Icon icon="mdi:account-group-outline" class="text-4xl text-indigo-500" />
      </div>
      
      <h1 class="text-xl font-bold mb-1">加入帳本邀請</h1>
      <p class="text-neutral-400 mb-6">
        <span class="text-white font-bold">{{ inviteInfo.creator_name }}</span> 邀請您加入
      </p>

      <div class="bg-neutral-800/50 p-4 rounded-2xl mb-8 border border-neutral-800">
        <div class="text-xs text-neutral-500 mb-1 uppercase tracking-wider">帳本名稱</div>
        <div class="text-lg font-bold text-indigo-400">{{ inviteInfo.ledger_name }}</div>
      </div>

      <button 
        @click="handleJoin" 
        :disabled="joining"
        class="w-full py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50"
      >
        {{ joining ? '加入中...' : '確認加入' }}
      </button>
      
      <button @click="router.push('/')" class="mt-4 text-sm text-neutral-500 hover:text-neutral-300 transition-colors">
        暫時不要
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const loading = ref(true)
const joining = ref(false)
const error = ref('')
const inviteInfo = ref<any>(null)

const fetchInviteInfo = async () => {
  try {
    const data = await api.get<any>(`/api/invites/${route.params.token}`)
    inviteInfo.value = data
  } catch (e: any) {
    error.value = e.message || '邀請連結已失效或不存在'
  } finally {
    loading.value = false
  }
}

const handleJoin = async () => {
  if (!isAuthenticated.value) {
    // Save current path to return after login
    localStorage.setItem('redirect_after_login', route.fullPath)
    router.push('/login')
    return
  }

  joining.value = true
  try {
    await api.post(`/api/invites/${route.params.token}/join`, {})
    // Redirect to the newly joined ledger
    router.push('/spaces')
  } catch (e: any) {
    alert(e.message || '加入失敗')
  } finally {
    joining.value = false
  }
}

onMounted(() => {
  initAuth()
  fetchInviteInfo()
})
</script>
