<template>
  <div class="space-settings-page">
    <SpaceHeader title="空間設定" />

    <div v-if="loading" class="text-center py-10 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-8">
      <!-- 1. Basic Info Section -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:cog-outline" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">基本設定</h2>
        </div>

        <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-6 flex flex-col gap-6">
          <!-- Cover Image -->
          <div>
            <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-3 ml-1">空間封面圖</label>
            <div class="relative h-32 rounded-xl overflow-hidden bg-neutral-800 group border border-neutral-700">
                <img v-if="form.cover_image" :src="getImageUrl(form.cover_image)" class="w-full h-full object-cover" />
                <div v-else class="w-full h-full flex items-center justify-center text-neutral-600">
                    <Icon icon="mdi:image-outline" class="text-4xl" />
                </div>
                <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <button @click="triggerImageUpload" class="px-4 py-2 bg-white text-black text-xs font-bold rounded-full border-0 cursor-pointer">更換相片</button>
                </div>
            </div>
          </div>

          <!-- Name -->
          <BaseInput
            v-model="form.name"
            label="空間名稱"
            placeholder="請輸入空間名稱"
          />

          <!-- Description -->
          <div class="flex flex-col gap-1.5">
            <label class="text-xs text-neutral-400 px-1">空間描述</label>
            <textarea
              v-model="form.description"
              placeholder="簡短描述此空間的用途"
              rows="3"
              class="w-full bg-neutral-800 border border-neutral-700 text-white py-1.5 px-2 rounded outline-none focus:border-indigo-500 transition-colors placeholder-neutral-500 text-base resize-none"
            />
          </div>

          <!-- Type & Currency (Readonly) -->
          <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-2 ml-1">空間類型</label>
                <div class="bg-neutral-800 text-neutral-400 py-3 px-4 rounded-xl text-sm border border-neutral-800">
                    {{ form.type === 'trip' ? '旅遊專案' : '個人記帳' }}
                </div>
              </div>
              <div>
                <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-2 ml-1">本位幣別</label>
                <div class="bg-neutral-800 text-neutral-400 py-3 px-4 rounded-xl text-sm border border-neutral-800 uppercase">
                    {{ form.base_currency }}
                </div>
              </div>
          </div>

          <!-- Date Range -->
          <div class="grid grid-cols-2 gap-4">
            <BaseInput
              v-model="form.start_date"
              type="date"
              label="開始日期"
            />
            <BaseInput
              v-model="form.end_date"
              type="date"
              label="結束日期"
            />
          </div>

          <!-- Is Pinned -->
          <label class="flex items-center justify-between p-4 rounded-xl border cursor-pointer transition-colors" :class="form.is_pinned ? 'border-indigo-500 bg-indigo-500/5' : 'border-neutral-700 bg-neutral-800'">
            <div class="flex flex-col">
              <span class="font-bold text-sm text-white">置頂空間</span>
              <span class="text-xs text-neutral-500 mt-0.5">將此空間釘選在列表最上方</span>
            </div>
            <input type="checkbox" v-model="form.is_pinned" class="w-5 h-5 rounded border-neutral-700 bg-neutral-800 text-indigo-500" />
          </label>

          <button
            @click="handleUpdateSpace"
            :disabled="updating"
            class="w-full py-4 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer mt-2 shadow-lg"
          >
            {{ updating ? '儲存中...' : '儲存基本設定' }}
          </button>
        </div>
      </section>

      <!-- 2. Lists Section -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:format-list-bulleted" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">分類與選項</h2>
        </div>

        <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-6 flex flex-col gap-6">
          <ListEditor
            v-model="form.currencies"
            label="可用幣別"
            placeholder="輸入幣別代碼，如 USD、JPY"
          />

          <ListEditor
            v-model="form.categories"
            label="交易分類"
            placeholder="輸入分類名稱，如 餐飲、交通"
          />

          <ListEditor
            v-model="form.payment_methods"
            label="付款方式"
            placeholder="輸入付款方式，如 現金、信用卡"
          />

          <button
            @click="handleUpdateLists"
            :disabled="updating"
            class="w-full py-4 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer mt-2 shadow-lg"
          >
            {{ updating ? '儲存中...' : '儲存分類設定' }}
          </button>
        </div>
      </section>

      <!-- 3. Members Section -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:account-group-outline" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">成員與權重</h2>
        </div>

        <div class="bg-neutral-900 rounded-2xl border border-neutral-800 overflow-hidden shadow-sm">
          <div v-for="member in detailStore.members" :key="member.user_id" class="p-5 border-b border-neutral-800 last:border-0 flex items-center justify-between hover:bg-neutral-800 transition-colors">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-400 font-bold text-lg border border-neutral-700">
                {{ (member.alias || member.user?.display_name || '?')[0].toUpperCase() }}
              </div>
              <div class="flex flex-col">
                <div class="flex items-center gap-2">
                  <span class="font-bold text-neutral-100">{{ member.alias || member.user?.display_name }}</span>
                  <span v-if="member.role === 'owner'" class="text-xs px-1.5 py-0.5 rounded bg-indigo-500/20 text-indigo-400 font-bold uppercase border border-indigo-500/20">建立者</span>
                </div>
                <span class="text-xs text-neutral-500">@{{ member.user?.username }}</span>
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
            <h2 class="text-lg font-bold">邀請連結</h2>
          </div>
          <button @click="showInviteModal = true" class="text-sm font-bold text-indigo-400 hover:text-indigo-300 transition-colors bg-transparent border-0 cursor-pointer">
            + 建立邀請
          </button>
        </div>

        <div v-if="detailStore.invites.length === 0" class="p-12 text-center bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed">
          <p class="text-neutral-500 text-sm italic">目前還沒有有效的邀請連結</p>
        </div>

        <div v-else class="flex flex-col gap-3">
          <div v-for="invite in detailStore.invites" :key="invite.id" class="p-5 bg-neutral-900 rounded-2xl border border-neutral-800 flex items-center justify-between hover:bg-neutral-800 transition-colors">
            <div class="flex flex-col gap-1">
              <div class="flex items-center gap-2">
                <span class="text-xs px-2 py-0.5 rounded-full bg-neutral-800 text-neutral-400 font-bold uppercase border border-neutral-700">
                  {{ invite.is_one_time ? '單次' : '多次' }}
                </span>
                <span v-if="invite.expires_at" class="text-xs text-neutral-600">
                  {{ formatExpiry(invite.expires_at) }} 過期
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
          <button @click="handleDeleteSpace" class="w-full py-4 rounded-xl bg-red-500/10 text-red-500 font-bold hover:bg-red-500 hover:text-white transition-all border border-red-500/20 cursor-pointer">
              刪除此空間
          </button>
      </section>
    </div>

    <!-- Hidden File Input -->
    <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleImageChange" />

    <!-- Modals -->
    <BaseModal v-model="showInviteModal" title="建立邀請連結">
        <div class="p-6 flex flex-col gap-6">
            <div class="flex flex-col gap-3">
                <label class="flex items-center gap-3 p-4 rounded-xl border cursor-pointer transition-colors" :class="inviteForm.is_one_time ? 'border-indigo-500 bg-indigo-500/5' : 'border-neutral-800 bg-neutral-800'">
                    <input type="checkbox" v-model="inviteForm.is_one_time" class="w-5 h-5 rounded border-neutral-700 bg-neutral-800 text-indigo-500" />
                    <div class="flex flex-col">
                        <span class="font-bold text-sm">單次邀請連結</span>
                        <span class="text-xs text-neutral-500 mt-0.5">連結僅供一人使用，加入後立即失效</span>
                    </div>
                </label>
            </div>

            <button
                @click="handleCreateInvite"
                :disabled="updating"
                class="w-full py-4 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-2 border-0 cursor-pointer"
            >
                建立連結
            </button>
        </div>
    </BaseModal>

    <BaseModal v-model="showAliasModal" title="修改成員別名">
        <div class="p-6 flex flex-col gap-6">
            <BaseInput
                v-model="aliasValue"
                :label="`成員: ${selectedMember?.user?.display_name}`"
                placeholder="請輸入成員別名"
            />

            <button
                @click="handleUpdateAlias"
                :disabled="updating"
                class="w-full py-4 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-2 border-0 cursor-pointer"
            >
                更新別名
            </button>
        </div>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import SpaceHeader from '~/components/SpaceHeader.vue'
import BaseInput from '~/components/BaseInput.vue'
import BaseModal from '~/components/BaseModal.vue'
import ListEditor from '~/components/ListEditor.vue'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { initAuth } = useAuth()
const { getImageUrl, uploadImage } = useImages()
const detailStore = useSpaceDetailStore()

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
  cover_image: '',
  is_pinned: false,
  currencies: [] as string[],
  categories: [] as string[],
  payment_methods: [] as string[]
})

const fileInput = ref<HTMLInputElement | null>(null)

// Modal States
const showInviteModal = ref(false)
const showAliasModal = ref(false)
const inviteForm = ref({ is_one_time: false })
const selectedMember = ref<any>(null)
const aliasValue = ref('')

const fetchData = async () => {
  try {
    detailStore.setSpaceId(spaceId)
    await Promise.all([
      detailStore.fetchSpace(true),
      detailStore.fetchMembers(true),
      detailStore.fetchInvites(true)
    ])
    const s = detailStore.space
    const parseJSON = (v: any): string[] => {
      if (Array.isArray(v)) return v
      if (typeof v === 'string') { try { return JSON.parse(v) } catch { return [] } }
      return []
    }
    const formatDate = (d: any): string | null => {
      if (!d) return null
      return new Date(d).toISOString().split('T')[0]
    }
    form.value = {
      name: s.name || '',
      description: s.description || '',
      type: s.type || '',
      base_currency: s.base_currency || '',
      start_date: formatDate(s.start_date),
      end_date: formatDate(s.end_date),
      cover_image: s.cover_image || '',
      is_pinned: s.is_pinned || false,
      currencies: parseJSON(s.currencies),
      categories: parseJSON(s.categories),
      payment_methods: parseJSON(s.payment_methods)
    }
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
      start_date: form.value.start_date ? new Date(form.value.start_date) : null,
      end_date: form.value.end_date ? new Date(form.value.end_date) : null,
      cover_image: form.value.cover_image,
      is_pinned: form.value.is_pinned
    })
    // Refresh store so other pages see updated data
    await detailStore.fetchSpace(true)
    alert('設定已儲存')
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    updating.value = false
  }
}

const handleUpdateLists = async () => {
  updating.value = true
  try {
    await api.patch(`/api/spaces/${spaceId}`, {
      currencies: form.value.currencies,
      categories: form.value.categories,
      payment_methods: form.value.payment_methods
    })
    await detailStore.fetchSpace(true)
    alert('分類設定已儲存')
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
        await handleUpdateSpace()
    } catch (e) {
        alert('上傳失敗')
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
        await detailStore.fetchMembers(true)
    } catch (e: any) {
        alert(e.message || '更新失敗')
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
        await detailStore.fetchInvites(true)
    } catch (e: any) {
        alert(e.message || '建立失敗')
    } finally {
        updating.value = false
    }
}

const handleRemoveMember = async (member: any) => {
  if (!confirm(`確定要移除成員 ${member.alias || member.user?.display_name} 嗎？`)) return
  try {
    await api.delete(`/api/spaces/${spaceId}/members/${member.user_id}`)
    await detailStore.fetchMembers(true)
  } catch (e: any) {
    alert(e.message || '移除失敗')
  }
}

const handleDeleteSpace = async () => {
    if (!confirm('警告：確定要刪除此空間嗎？這將會永久刪除所有交易、比價資料與成員關聯，且無法復原。')) return
    try {
        await api.delete(`/api/spaces/${spaceId}`)
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
  alert('連結已複製')
}

const handleRevokeInvite = async (invite_id: string) => {
  if (!confirm('確定要撤銷此邀請連結嗎？')) return
  try {
    await api.delete(`/api/spaces/${invite_id}`)
    await detailStore.fetchInvites(true)
  } catch (e: any) {
    alert('撤銷失敗')
  }
}

onMounted(() => {
  initAuth()
  fetchData()
})
</script>
