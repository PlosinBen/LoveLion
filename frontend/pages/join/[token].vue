<template>
  <div class="join-space-page flex flex-col items-center justify-center p-6 text-neutral-50">
    <div v-if="loading" class="text-neutral-400 animate-pulse font-bold">
      正在驗證邀請連結...
    </div>

    <BaseCard v-else-if="error" padding="p-8" class="text-center max-w-sm w-full shadow-2xl">
      <div class="w-20 h-20 bg-red-500/10 rounded-full flex items-center justify-center mx-auto mb-6">
        <Icon icon="mdi:alert-circle-outline" class="text-5xl text-red-500" />
      </div>
      <h1 class="text-xl font-bold mb-2">邀請無效</h1>
      <p class="text-neutral-400 mb-8 leading-relaxed">{{ error }}</p>
      <NuxtLink to="/" class="w-full inline-flex justify-center items-center px-4 py-2.5 text-sm rounded bg-neutral-800 text-neutral-400 font-bold hover:text-white hover:bg-neutral-700 no-underline transition-all active:scale-95">
        回到首頁
      </NuxtLink>
    </BaseCard>

    <BaseCard v-else-if="inviteInfo" padding="p-8" class="text-center max-w-sm w-full shadow-2xl">
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
        variant="primary"
        class="w-full shadow-lg"
      >
        接受邀請並加入
      </BaseButton>

      <NuxtLink to="/" class="mt-6 text-sm text-neutral-500 hover:text-neutral-300 no-underline font-bold transition-colors">
        暫時不要
      </NuxtLink>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useLoading } from '~/composables/useLoading'
import BaseCard from '~/components/BaseCard.vue'
import type { InviteInfo } from '~/types'

definePageMeta({
  layout: 'empty'
})

const route = useRoute()
const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const { showLoading, hideLoading } = useLoading()

const loading = ref(true)
const error = ref('')
const inviteInfo = ref<InviteInfo | null>(null)

const fetchInviteInfo = async () => {
  try {
    const data = await api.get<InviteInfo>(`/api/invites/${route.params.token}`)
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

  showLoading()
  try {
    await api.post(`/api/invites/${route.params.token}/join`, {})
    router.push('/')
  } catch (e: any) {
    alert(e.message || '加入失敗')
  } finally {
    hideLoading()
  }
}

onMounted(() => {
  initAuth()
  fetchInviteInfo()
})
</script>
