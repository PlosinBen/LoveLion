<template>
  <div class="ledger-settings-page pb-24">
    <header class="flex justify-between items-center mb-8">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">帳本分享設定</h1>
      <div class="w-10"></div>
    </header>

    <div v-if="loading" class="text-center py-10 text-neutral-500">載入中...</div>

    <div v-else class="flex flex-col gap-8">
      <!-- Members Section -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:account-group" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">帳本成員</h2>
        </div>
        
        <div class="bg-neutral-900 rounded-3xl border border-neutral-800 overflow-hidden">
          <div v-for="member in members" :key="member.user_id" class="p-5 border-b border-neutral-800 last:border-0 flex items-center justify-between">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-2xl bg-neutral-800 flex items-center justify-center text-indigo-400 font-bold text-lg">
                {{ (member.alias || member.user?.display_name || '?')[0] }}
              </div>
              <div class="flex flex-col">
                <div class="flex items-center gap-2">
                  <span class="font-bold text-neutral-100">{{ member.alias || member.user?.display_name }}</span>
                  <span v-if="member.role === 'owner'" class="text-[10px] px-1.5 py-0.5 rounded-md bg-indigo-500/10 text-indigo-400 font-bold uppercase tracking-wider">Owner</span>
                </div>
                <span v-if="member.alias" class="text-xs text-neutral-500">原名: {{ member.user?.display_name }}</span>
                <span v-else class="text-xs text-neutral-500">@{{ member.user?.username }}</span>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <button 
                v-if="member.role !== 'owner'"
                @click="openAliasModal(member)"
                class="p-2 rounded-xl bg-neutral-800 text-neutral-400 hover:text-white transition-colors border-0 cursor-pointer"
              >
                <Icon icon="mdi:pencil-outline" class="text-xl" />
              </button>
              <button 
                v-if="member.role !== 'owner'"
                @click="handleRemoveMember(member)"
                class="p-2 rounded-xl bg-red-500/10 text-red-500 hover:bg-red-500/20 transition-colors border-0 cursor-pointer"
              >
                <Icon icon="mdi:account-remove-outline" class="text-xl" />
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- Invite Section -->
      <section>
        <div class="flex items-center justify-between mb-4 px-1">
          <div class="flex items-center gap-2">
            <Icon icon="mdi:link-variant" class="text-indigo-500 text-xl" />
            <h2 class="text-lg font-bold">邀請連結</h2>
          </div>
          <button @click="showInviteModal = true" class="text-sm font-bold text-indigo-400 hover:text-indigo-300 transition-colors bg-transparent border-0 cursor-pointer">
            + 產生連結
          </button>
        </div>

        <div v-if="activeInvites.length === 0" class="p-10 text-center bg-neutral-900 rounded-3xl border border-neutral-800 border-dashed">
          <p class="text-neutral-500 text-sm">目前沒有有效的邀請連結</p>
        </div>

        <div v-else class="flex flex-col gap-3">
          <div v-for="invite in activeInvites" :key="invite.id" class="p-5 bg-neutral-900 rounded-3xl border border-neutral-800 flex items-center justify-between">
            <div class="flex flex-col gap-1">
              <div class="flex items-center gap-2">
                <span class="text-xs px-2 py-0.5 rounded-full bg-neutral-800 text-neutral-400 font-medium">
                  {{ invite.is_one_time ? '一次性' : '多人共用' }}
                </span>
                <span v-if="invite.expires_at" class="text-[10px] text-neutral-600">
                  {{ formatExpiry(invite.expires_at) }} 到期
                </span>
              </div>
              <div class="text-sm font-mono text-neutral-500 truncate w-40">...{{ invite.token.slice(-8) }}</div>
            </div>
            
            <div class="flex items-center gap-2">
              <button @click="copyInviteLink(invite.token)" class="p-2 rounded-xl bg-indigo-500/10 text-indigo-400 hover:bg-indigo-500/20 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:content-copy" class="text-xl" />
              </button>
              <button @click="handleRevokeInvite(invite.id)" class="p-2 rounded-xl bg-neutral-800 text-neutral-500 hover:text-red-500 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:trash-can-outline" class="text-xl" />
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>

    <!-- Invite Modal -->
    <div v-if="showInviteModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm z-[100] flex items-end sm:items-center justify-center p-4">
      <div class="bg-neutral-900 w-full max-w-md rounded-t-3xl sm:rounded-3xl border border-neutral-800 p-8 animate-slide-up">
        <h3 class="text-xl font-bold mb-6">產生共享連結</h3>
        
        <div class="flex flex-col gap-6 mb-8">
          <div>
            <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-3">連結類型</label>
            <div class="grid grid-cols-2 gap-3">
              <button 
                @click="inviteForm.is_one_time = true"
                class="py-3 rounded-2xl border transition-all font-bold"
                :class="inviteForm.is_one_time ? 'bg-indigo-500 border-indigo-500 text-white' : 'bg-neutral-800 border-neutral-700 text-neutral-400'"
              >
                一次性
              </button>
              <button 
                @click="inviteForm.is_one_time = false"
                class="py-3 rounded-2xl border transition-all font-bold"
                :class="!inviteForm.is_one_time ? 'bg-indigo-500 border-indigo-500 text-white' : 'bg-neutral-800 border-neutral-700 text-neutral-400'"
              >
                多人共用
              </button>
            </div>
          </div>

          <div>
            <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-3">有效期限</label>
            <select v-model="inviteForm.expiryOption" class="w-full bg-neutral-800 border-neutral-700 text-white py-4 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors">
              <option value="1h">1 小時</option>
              <option value="24h">24 小時</option>
              <option value="7d">7 天</option>
              <option value="never">永久有效</option>
            </select>
          </div>
        </div>

        <div class="flex gap-3">
          <button @click="showInviteModal = false" class="flex-1 py-4 rounded-2xl bg-neutral-800 font-bold hover:bg-neutral-700 transition-colors border-0 text-white cursor-pointer">取消</button>
          <button @click="handleCreateInvite" class="flex-1 py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 border-0 cursor-pointer">產生並複製</button>
        </div>
      </div>
    </div>

    <!-- Alias Modal -->
    <div v-if="showAliasModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
      <div class="bg-neutral-900 w-full max-w-sm rounded-3xl border border-neutral-800 p-8 animate-scale-in">
        <h3 class="text-xl font-bold mb-2">設定別名</h3>
        <p class="text-sm text-neutral-500 mb-6">設定只有您看得見的稱呼</p>
        
        <input 
          v-model="aliasForm.alias" 
          type="text" 
          placeholder="例如：老婆、室友 A" 
          class="w-full bg-neutral-800 border-neutral-700 text-white py-4 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors mb-8"
          autofocus
        />

        <div class="flex gap-3">
          <button @click="showAliasModal = false" class="flex-1 py-3 rounded-xl bg-neutral-800 font-bold border-0 text-white cursor-pointer">取消</button>
          <button @click="handleUpdateAlias" class="flex-1 py-3 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all border-0 cursor-pointer">儲存</button>
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

const route = useRoute()
const router = useRouter()
const api = useApi()
const { initAuth } = useAuth()

const loading = ref(true)
const members = ref<any[]>([])
const activeInvites = ref<any[]>([])

const showInviteModal = ref(false)
const inviteForm = ref({
  is_one_time: true,
  expiryOption: '24h'
})

const showAliasModal = ref(false)
const aliasForm = ref({
  user_id: '',
  alias: ''
})

const fetchMembers = async () => {
  try {
    members.value = await api.get<any[]>(`/api/ledgers/${route.params.id}/members`)
  } catch (e) {
    console.error('Failed to fetch members:', e)
  } finally {
    loading.value = false
  }
}

const fetchInvites = async () => {
  try {
    activeInvites.value = await api.get<any[]>(`/api/ledgers/${route.params.id}/invites`)
  } catch (e) {
    console.error('Failed to fetch invites:', e)
  }
}

const handleCreateInvite = async () => {
  try {
    let expires_at = null
    const now = new Date()
    if (inviteForm.value.expiryOption === '1h') expires_at = new Date(now.getTime() + 3600000)
    else if (inviteForm.value.expiryOption === '24h') expires_at = new Date(now.getTime() + 86400000)
    else if (inviteForm.value.expiryOption === '7d') expires_at = new Date(now.getTime() + 604800000)

    const invite = await api.post<any>(`/api/ledgers/${route.params.id}/invites`, {
      is_one_time: inviteForm.value.is_one_time,
      expires_at: expires_at?.toISOString() || null,
      max_uses: inviteForm.value.is_one_time ? 1 : 0
    })

    const joinUrl = `${window.location.origin}/join/${invite.token}`
    await navigator.clipboard.writeText(joinUrl)
    alert('連結已產生並複製到剪貼簿！')
    
    await fetchInvites()
    showInviteModal.value = false
  } catch (e: any) {
    alert(e.message || '產生連結失敗')
  }
}

const copyInviteLink = (token: string) => {
  const joinUrl = `${window.location.origin}/join/${token}`
  navigator.clipboard.writeText(joinUrl)
  alert('連結已複製！')
}

const handleRevokeInvite = async (invite_id: string) => {
  if (!confirm('確定要撤回此邀請連結嗎？撤回後，持有連結的使用者將無法再加入此帳本。')) return
  try {
    await api.del(`/api/ledgers/${route.params.id}/invites/${invite_id}`)
    await fetchInvites()
  } catch (e: any) {
    alert(e.message || '撤回失敗')
  }
}

const openAliasModal = (member: any) => {
  aliasForm.value = {
    user_id: member.user_id,
    alias: member.alias || ''
  }
  showAliasModal.value = true
}

const handleUpdateAlias = async () => {
  try {
    await api.patch(`/api/ledgers/${route.params.id}/members/${aliasForm.value.user_id}`, {
      alias: aliasForm.value.alias
    })
    await fetchMembers()
    showAliasModal.value = false
  } catch (e: any) {
    alert(e.message || '更新失敗')
  }
}

const handleRemoveMember = async (member: any) => {
  const name = member.alias || member.user?.display_name
  if (!confirm(`確定要將 ${name} 從帳本中移除嗎？`)) return

  try {
    await api.del(`/api/ledgers/${route.params.id}/members/${member.user_id}`)
    await fetchMembers()
  } catch (e: any) {
    alert(e.message || '移除失敗')
  }
}

const formatExpiry = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  initAuth()
  fetchMembers()
  fetchInvites()
})
</script>

<style scoped>
.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}
.animate-scale-in {
  animation: scale-in 0.2s ease-out;
}
.animate-fade-in {
  animation: fade-in 0.2s ease-out;
}

@keyframes slide-up {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}
@keyframes scale-in {
  from { transform: scale(0.9); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
