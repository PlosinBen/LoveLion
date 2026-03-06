<template>
  <div class="space-settings-page pb-24 px-4">
    <!-- Header -->
    <div class="px-2 pt-0 pb-8 flex items-center gap-3">
      <button @click="router.back()" class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer shrink-0">
          <Icon icon="mdi:arrow-left" class="text-xl" />
      </button>
      <h1 class="text-xl font-bold text-white tracking-tight">空間設定</h1>
    </div>

    <div v-if="loading" class="text-center py-10 text-neutral-500">
      <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-8">
      <!-- 1. Basic Info Section -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:cog-outline" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">基本設定</h2>
        </div>
        
        <div class="bg-neutral-900 rounded-3xl border border-neutral-800 p-6 flex flex-col gap-6">
          <!-- Cover Image -->
          <div>
            <label class="block text-xs text-neutral-500 uppercase tracking-widest mb-3 ml-1">空間封面</label>
            <div class="relative h-32 rounded-2xl overflow-hidden bg-neutral-800 group border border-neutral-700/50">
                <img v-if="form.cover_image" :src="getImageUrl(form.cover_image)" class="w-full h-full object-cover" />
                <div v-else class="w-full h-full flex items-center justify-center text-neutral-600">
                    <Icon icon="mdi:image-outline" class="text-4xl" />
                </div>
                <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <button @click="triggerImageUpload" class="px-4 py-2 bg-white text-black text-xs font-bold rounded-full border-0 cursor-pointer">更換照片</button>
                </div>
            </div>
          </div>

          <!-- Name -->
          <div>
            <label class="block text-xs text-neutral-500 uppercase tracking-widest mb-2 ml-1">名稱</label>
            <input 
              v-model="form.name" 
              type="text" 
              class="w-full bg-neutral-800 border-neutral-700 text-white py-3 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors"
              placeholder="空間名稱"
            />
          </div>

          <!-- Type & Currency (Readonly for now) -->
          <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs text-neutral-500 uppercase tracking-widest mb-2 ml-1">類型</label>
                <div class="bg-neutral-800/50 text-neutral-400 py-3 px-4 rounded-2xl text-sm border border-neutral-800">
                    {{ form.type === 'trip' ? '旅行專案' : '個人記帳' }}
                </div>
              </div>
              <div>
                <label class="block text-xs text-neutral-500 uppercase tracking-widest mb-2 ml-1">主幣別</label>
                <div class="bg-neutral-800/50 text-neutral-400 py-3 px-4 rounded-2xl text-sm border border-neutral-800 uppercase">
                    {{ form.base_currency }}
                </div>
              </div>
          </div>

          <!-- Dates (Only for Trip) -->
          <div v-if="form.type === 'trip'" class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs text-neutral-500 uppercase tracking-widest mb-2 ml-1">開始日期</label>
                <input type="date" v-model="startDateStr" class="w-full bg-neutral-800 border-neutral-700 text-white py-3 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors text-sm" />
              </div>
              <div>
                <label class="block text-xs text-neutral-500 uppercase tracking-widest mb-2 ml-1">結束日期</label>
                <input type="date" v-model="endDateStr" class="w-full bg-neutral-800 border-neutral-700 text-white py-3 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors text-sm" />
              </div>
          </div>

          <button 
            @click="handleUpdateSpace" 
            :disabled="updating"
            class="w-full py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all disabled:opacity-50 border-0 cursor-pointer mt-2"
          >
            {{ updating ? '儲存中...' : '儲存變更' }}
          </button>
        </div>
      </section>

      <!-- 2. Members Section -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:account-group-outline" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">成員管理</h2>
        </div>
        
        <div class="bg-neutral-900 rounded-3xl border border-neutral-800 overflow-hidden">
          <div v-for="member in members" :key="member.user_id" class="p-5 border-b border-neutral-800 last:border-0 flex items-center justify-between hover:bg-neutral-800/30 transition-colors">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-2xl bg-neutral-800 flex items-center justify-center text-indigo-400 font-black text-lg border border-neutral-700/50">
                {{ (member.alias || member.user?.display_name || '?')[0].toUpperCase() }}
              </div>
              <div class="flex flex-col">
                <div class="flex items-center gap-2">
                  <span class="font-bold text-neutral-100">{{ member.alias || member.user?.display_name }}</span>
                  <span v-if="member.role === 'owner'" class="text-xs px-1.5 py-0.5 rounded-md bg-indigo-500/20 text-indigo-400 font-bold uppercase tracking-wider border border-indigo-500/20">主辦人</span>
                </div>
                <span class="text-xs text-neutral-500 font-medium">@{{ member.user?.username }}</span>
              </div>
            </div>

            <div class="flex items-center gap-1">
              <button @click="openAliasModal(member)" class="p-2.5 rounded-xl bg-neutral-800 text-neutral-400 hover:text-white transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:pencil-outline" class="text-lg" />
              </button>
              <button v-if="member.role !== 'owner'" @click="handleRemoveMember(member)" class="p-2.5 rounded-xl bg-red-500/10 text-red-500 hover:bg-red-500/20 transition-colors border-0 cursor-pointer ml-1">
                <Icon icon="mdi:account-remove-outline" class="text-lg" />
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 3. Invites Section -->
      <section>
        <div class="flex items-center justify-between mb-4 px-1">
          <div class="flex items-center gap-2">
            <Icon icon="mdi:link-variant" class="text-indigo-500 text-xl" />
            <h2 class="text-lg font-bold">共享連結</h2>
          </div>
          <button @click="showInviteModal = true" class="text-sm font-bold text-indigo-400 hover:text-indigo-300 transition-colors bg-transparent border-0 cursor-pointer">
            + 產生連結
          </button>
        </div>

        <div v-if="invites.length === 0" class="p-12 text-center bg-neutral-900 rounded-3xl border border-neutral-800 border-dashed">
          <p class="text-neutral-500 text-sm italic">目前沒有有效的邀請連結</p>
        </div>

        <div v-else class="flex flex-col gap-3">
          <div v-for="invite in invites" :key="invite.id" class="p-5 bg-neutral-900 rounded-3xl border border-neutral-800 flex items-center justify-between hover:bg-neutral-800/20 transition-colors">
            <div class="flex flex-col gap-1">
              <div class="flex items-center gap-2">
                <span class="text-xs px-2 py-0.5 rounded-full bg-neutral-800 text-neutral-400 font-bold uppercase tracking-wider border border-neutral-700">
                  {{ invite.is_one_time ? '一次性' : '多人共用' }}
                </span>
                <span v-if="invite.expires_at" class="text-xs text-neutral-600 font-medium">
                  {{ formatExpiry(invite.expires_at) }} 到期
                </span>
              </div>
              <div class="text-xs font-mono text-neutral-500 mt-1">...{{ invite.token.slice(-12) }}</div>
            </div>
            
            <div class="flex items-center gap-2">
              <button @click="copyInviteLink(invite.token)" class="p-2.5 rounded-xl bg-indigo-500/10 text-indigo-400 hover:bg-indigo-500/20 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:content-copy" class="text-xl" />
              </button>
              <button @click="handleRevokeInvite(invite.id)" class="p-2.5 rounded-xl bg-neutral-800 text-neutral-500 hover:text-red-500 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:trash-can-outline" class="text-xl" />
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- Danger Zone -->
      <section class="mt-4 pt-8 border-t border-neutral-800">
          <button @click="handleDeleteSpace" class="w-full py-4 rounded-2xl bg-red-500/10 text-red-500 font-bold hover:bg-red-500 hover:text-white transition-all border border-red-500/20 cursor-pointer">
              刪除此空間
          </button>
      </section>
    </div>

    <!-- Hidden File Input for Image Upload -->
    <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleImageChange" />

    <!-- Modals -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
        <!-- Invite Modal -->
        <div v-if="showInviteModal" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
            <div class="w-full max-w-lg bg-neutral-900 rounded-t-3xl sm:rounded-3xl p-8 border border-neutral-800 animate-in slide-in-from-bottom duration-300">
                <div class="flex justify-between items-center mb-6">
                    <h2 class="text-xl font-bold">產生邀請連結</h2>
                    <button @click="showInviteModal = false" class="w-10 h-10 rounded-full bg-neutral-800 flex items-center justify-center text-neutral-400 border-0 cursor-pointer">
                        <Icon icon="mdi:close" class="text-xl" />
                    </button>
                </div>

                <div class="flex flex-col gap-6">
                    <div class="flex flex-col gap-3">
                        <label class="flex items-center gap-3 p-4 rounded-2xl border cursor-pointer transition-colors" :class="inviteForm.is_one_time ? 'border-indigo-500 bg-indigo-500/5' : 'border-neutral-800 bg-neutral-800/20'">
                            <input type="checkbox" v-model="inviteForm.is_one_time" class="w-5 h-5 rounded border-neutral-700 bg-neutral-800 text-indigo-500" />
                            <div class="flex flex-col">
                                <span class="font-bold text-sm">一次性連結</span>
                                <span class="text-xs text-neutral-500 mt-0.5">使用一次後立即失效（適合單人邀請）</span>
                            </div>
                        </label>
                    </div>

                    <button 
                        @click="handleCreateInvite" 
                        :disabled="updating"
                        class="w-full py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-2 border-0 cursor-pointer"
                    >
                        產生連結
                    </button>
                </div>
            </div>
        </div>

        <!-- Alias Modal -->
        <div v-else-if="showAliasModal" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
            <div class="w-full max-w-lg bg-neutral-900 rounded-t-3xl sm:rounded-3xl p-8 border border-neutral-800 animate-in slide-in-from-bottom duration-300">
                <div class="flex justify-between items-center mb-6">
                    <h2 class="text-xl font-bold">修改成員暱稱</h2>
                    <button @click="showAliasModal = false" class="w-10 h-10 rounded-full bg-neutral-800 flex items-center justify-center text-neutral-400 border-0 cursor-pointer">
                        <Icon icon="mdi:close" class="text-xl" />
                    </button>
                </div>

                <div class="flex flex-col gap-6">
                    <div>
                        <label class="block text-xs text-neutral-500 uppercase tracking-widest mb-3 ml-1">成員：{{ selectedMember?.user?.display_name }}</label>
                        <input 
                            v-model="aliasValue" 
                            type="text" 
                            class="w-full bg-neutral-800 border-neutral-700 text-white py-4 px-4 rounded-2xl outline-none focus:border-indigo-500 transition-colors"
                            placeholder="請輸入暱稱"
                        />
                    </div>

                    <button 
                        @click="handleUpdateAlias" 
                        :disabled="updating"
                        class="w-full py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-2 border-0 cursor-pointer"
                    >
                        儲存暱稱
                    </button>
                </div>
            </div>
        </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { initAuth } = useAuth()
const { getImageUrl, uploadImage } = useImages()

const loading = ref(true)
const updating = ref(false)
const spaceId = route.params.id as string

const form = ref({
  name: '',
  description: '',
  type: '',
  base_currency: '',
  start_date: null as string | null,
  end_date: null as string | null,
  cover_image: ''
})

const members = ref<any[]>([])
const invites = ref<any[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

// Modal States
const showInviteModal = ref(false)
const showAliasModal = ref(false)
const inviteForm = ref({ is_one_time: false })
const selectedMember = ref<any>(null)
const aliasValue = ref('')

// Date conversion helpers
const startDateStr = computed({
    get: () => form.value.start_date ? form.value.start_date.split('T')[0] : '',
    set: (val) => form.value.start_date = val ? new Date(val).toISOString() : null
})
const endDateStr = computed({
    get: () => form.value.end_date ? form.value.end_date.split('T')[0] : '',
    set: (val) => form.value.end_date = val ? new Date(val).toISOString() : null
})

const fetchData = async () => {
  try {
    const [spaceData, membersData, invitesData] = await Promise.all([
      api.get<any>(`/api/spaces/${spaceId}`),
      api.get<any[]>(`/api/spaces/${spaceId}/members`),
      api.get<any[]>(`/api/spaces/${spaceId}/invites`)
    ])
    form.value = { ...spaceData }
    members.value = membersData
    invites.value = invitesData
  } catch (e) {
    console.error('Failed to fetch data:', e)
    router.push('/')
  } finally {
    loading.value = false
  }
}

const handleUpdateSpace = async () => {
  updating.value = true
  try {
    await api.patch(`/api/spaces/${spaceId}`, {
      name: form.value.name,
      description: form.value.description,
      start_date: form.value.start_date,
      end_date: form.value.end_date,
      cover_image: form.value.cover_image
    })
    alert('空間設定已儲存！')
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    updating.value = false
  }
}

const triggerImageUpload = () => fileInput.value?.click()

const handleImageChange = async (e: Event) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    
    try {
        const result = await uploadImage(file, spaceId, 'space_cover')
        form.value.cover_image = result.file_path
        // Auto-save change
        await handleUpdateSpace()
    } catch (e) {
        alert('照片上傳失敗')
    }
}

const openAliasModal = (member: any) => {
    selectedMember.value = member
    aliasValue.value = member.alias || ''
    showAliasModal.value = true
}

const handleUpdateAlias = async () => {
    if (!selectedMember.value) return
    updating.value = true
    try {
        await api.put(`/api/spaces/${spaceId}/members/${selectedMember.value.user_id}`, {
            alias: aliasValue.value
        })
        showAliasModal.value = false
        await fetchData()
    } catch (e: any) {
        alert(e.message || '更新暱稱失敗')
    } finally {
        updating.value = false
    }
}

const handleCreateInvite = async () => {
    updating.value = true
    try {
        await api.post(`/api/spaces/${spaceId}/invites`, {
            is_one_time: inviteForm.value.is_one_time
        })
        showInviteModal.value = false
        await fetchData()
    } catch (e: any) {
        alert(e.message || '產生連結失敗')
    } finally {
        updating.value = false
    }
}

const handleRemoveMember = async (member: any) => {
  if (!confirm(`確定要將 ${member.alias || member.user?.display_name} 移出空間？`)) return
  try {
    await api.del(`/api/spaces/${spaceId}/members/${member.user_id}`)
    await fetchData()
  } catch (e: any) {
    alert(e.message || '移除失敗')
  }
}

const handleDeleteSpace = async () => {
    if (!confirm('🚨 警告：確定要永久刪除此空間嗎？所有交易紀錄與比價資料將無法復原。')) return
    try {
        await api.del(`/api/spaces/${spaceId}`)
        router.push('/')
    } catch (e: any) {
        alert(e.message || '刪除失敗')
    }
}

const formatExpiry = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-TW', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

const copyInviteLink = (token: string) => {
  const joinUrl = `${window.location.origin}/join/${token}`
  navigator.clipboard.writeText(joinUrl)
  alert('連結已複製！')
}

const handleRevokeInvite = async (invite_id: string) => {
  if (!confirm('確定要撤回此連結？')) return
  try {
    await api.del(`/api/spaces/${spaceId}/invites/${invite_id}`)
    await fetchData()
  } catch (e: any) {
    alert('撤回失敗')
  }
}

onMounted(() => {
  initAuth()
  fetchData()
})
</script>
