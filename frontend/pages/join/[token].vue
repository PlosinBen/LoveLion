<template>
  <div class="join-space-page flex flex-col items-center justify-center p-6 text-neutral-50">
    <div v-if="loading" class="text-neutral-400 animate-pulse font-bold">
      正在驗證邀請連結...
    </div>

    <div v-else-if="error" class="bg-neutral-900 p-8 rounded-2xl border border-neutral-800 text-center max-w-sm w-full shadow-2xl">
      <div class="w-20 h-20 bg-red-500/10 rounded-full flex items-center justify-center mx-auto mb-6">
        <Icon icon="mdi:alert-circle-outline" class="text-5xl text-red-500" />
      </div>
      <h1 class="text-xl font-bold mb-2">邀請無效</h1>
      <p class="text-neutral-400 mb-8 leading-relaxed">{{ error }}</p>
      <BaseButton @click="router.push('/')" variant="secondary" full-width>
        回到首頁
      </BaseButton>
    </div>

    <div v-else-if="inviteInfo" class="bg-neutral-900 p-8 rounded-2xl border border-neutral-800 text-center max-w-sm w-full shadow-2xl">
      <div class="w-20 h-20 bg-indigo-500/10 rounded-full flex items-center justify-center mx-auto mb-6">
        <Icon icon="mdi:account-group-outline" class="text-5xl text-indigo-500" />
      </div>
      
      <h1 class="text-xl font-bold mb-1">受邀加入空間</h1>
      <p class="text-neutral-400 mb-8 font-medium">
        <span class="text-white font-bold">{{ inviteInfo.creator_name }}</span> 邀請您一同協作
      </p>

      <div class="bg-neutral-800 p-5 rounded-xl mb-8 border border-neutral-700">
        <div class="text-xs text-neutral-500 mb-1 uppercase font-bold tracking-wider">空間名稱</div>
        <div class="text-xl font-bold text-indigo-400">{{ inviteInfo.space_name }}</div>
      </div>

      <BaseButton 
        @click="handleJoin" 
        :loading="joining"
        variant="primary"
        full-width
        class="shadow-lg"
      >
        接受邀請並加入
      </BaseButton>

      <BaseButton @click="router.push('/')" variant="ghost" class="mt-6 text-sm">
        暫時不要
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import BaseButton from '~/components/BaseButton.vue'

definePageMeta({
  layout: 'empty'
})

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
    localStorage.setItem('redirect_after_login', route.fullPath)
    router.push('/login')
    return
  }

  joining.value = true
  try {
    await api.post(`/api/invites/${route.params.token}/join`, {})
    router.push('/')
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
