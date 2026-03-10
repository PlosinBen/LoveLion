<template>
  <div class="space-settings-page">
    <PageTitle title="空間設定" :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${spaceId}/ledger` }]" />

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

        <BaseCard padding="p-6" class="flex flex-col gap-6">
          <!-- Cover Image -->
          <div>
            <label class="block text-xs text-neutral-500 uppercase tracking-wider mb-3 ml-1">空間封面圖</label>
            <ImageManager
              :entity-id="spaceId"
              entity-type="space"
              :max-count="1"
              :allow-reorder="false"
            />
          </div>

          <!-- Name -->
          <BaseInput
            v-model="form.name"
            label="空間名稱"
            placeholder="請輸入空間名稱"
          />

          <!-- Pinned Status -->
          <div class="flex items-center justify-between p-4 bg-neutral-800/50 rounded-xl border border-neutral-800">
            <div class="flex flex-col">
              <span class="font-bold text-sm">置頂此空間</span>
              <span class="text-xs text-neutral-500">將此空間固定在首頁頂端</span>
            </div>
            <input type="checkbox" v-model="form.is_pinned" class="w-6 h-6 rounded border-neutral-700 bg-neutral-800 text-indigo-500 focus:ring-indigo-500 focus:ring-offset-neutral-900" />
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

          <BaseButton
            @click="handleUpdateSpace"
            variant="primary"
            class="w-full mt-2 shadow-lg"
          >
            儲存基本設定
          </BaseButton>
        </BaseCard>
      </section>

      <!-- List Configurations -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:format-list-bulleted" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">分類與選項</h2>
        </div>

        <BaseCard padding="p-6" class="flex flex-col gap-8">
          <ListEditor v-model="form.categories" label="交易分類" placeholder="新增分類..." />
          <ListEditor v-model="form.currencies" label="支援幣別" placeholder="例如: JPY, USD..." />
          <ListEditor v-model="form.payment_methods" label="付款方式" placeholder="例如: 現金, 信用卡..." />

          <BaseButton
            @click="handleUpdateLists"
            variant="primary"
            class="w-full shadow-lg"
          >
            儲存分類設定
          </BaseButton>
        </BaseCard>
      </section>

      <!-- 2. Members Section -->
      <section>
        <div class="flex items-center gap-2 mb-4 px-1">
          <Icon icon="mdi:account-group-outline" class="text-indigo-500 text-xl" />
          <h2 class="text-lg font-bold">成員與權重</h2>
        </div>

        <BaseCard padding="" class="overflow-hidden shadow-sm">
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
              <button @click="openAliasModal(member)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-800 text-neutral-400 hover:text-white hover:bg-neutral-700 border-0 cursor-pointer transition-colors active:scale-95">
                <Icon icon="mdi:pencil-outline" class="text-lg" />
              </button>
              <button v-if="member.role !== 'owner'" @click="handleRemoveMember(member)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white border border-red-500/20 cursor-pointer transition-colors active:scale-95 ml-1">
                <Icon icon="mdi:account-remove-outline" class="text-lg" />
              </button>
            </div>
          </div>
        </BaseCard>
      </section>

      <!-- 3. Invites Section -->
      <section>
        <div class="flex items-center justify-between mb-4 px-1">
          <div class="flex items-center gap-2">
            <Icon icon="mdi:link-variant" class="text-indigo-500 text-xl" />
            <h2 class="text-lg font-bold">邀請連結</h2>
          </div>
          <button @click="showInviteModal = true" class="bg-transparent text-indigo-400 font-bold text-sm border-0 cursor-pointer hover:text-indigo-300 transition-colors active:scale-95">
            + 建立邀請
          </button>
        </div>

        <div v-if="detailStore.invites.length === 0" class="p-12 text-center bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed">
          <p class="text-neutral-500 text-sm italic">目前還沒有有效的邀請連結</p>
        </div>

        <div v-else class="flex flex-col gap-3">
          <BaseCard v-for="invite in detailStore.invites" :key="invite.id" padding="p-5" class="flex items-center justify-between hover:bg-neutral-800 transition-colors">
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
              <button @click="copyInviteLink(invite.token)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-indigo-500/10 text-indigo-400 border-0 cursor-pointer hover:bg-indigo-500/20 transition-colors active:scale-95">
                <Icon icon="mdi:content-copy" class="text-xl" />
              </button>
              <button @click="handleRevokeInvite(invite.id)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-800 text-neutral-400 hover:text-white hover:bg-neutral-700 border-0 cursor-pointer transition-colors active:scale-95">
                <Icon icon="mdi:trash-can-outline" class="text-xl" />
              </button>
            </div>
          </BaseCard>
        </div>
      </section>

      <!-- Danger Zone -->
      <section class="mt-4 pt-8 border-t border-neutral-800">
          <BaseButton @click="handleDeleteSpace" variant="danger" class="w-full">
              刪除此空間
          </BaseButton>
      </section>
    </div>

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

            <BaseButton
                @click="handleCreateInvite"
                variant="primary"
                class="w-full mt-2"
            >
                建立連結
            </BaseButton>
        </div>
    </BaseModal>

    <BaseModal v-model="showAliasModal" title="修改成員別名">
        <div class="p-6 flex flex-col gap-6">
            <BaseInput
                v-model="aliasValue"
                :label="`成員: ${selectedMember?.user?.display_name}`"
                placeholder="請輸入成員別名"
            />

            <BaseButton
                @click="handleUpdateAlias"
                variant="primary"
                class="w-full mt-2"
            >
                更新別名
            </BaseButton>
        </div>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import BaseInput from '~/components/BaseInput.vue'
import ImageManager from '~/components/ImageManager.vue'
import BaseModal from '~/components/BaseModal.vue'
import ListEditor from '~/components/ListEditor.vue'
import { useLoading } from '~/composables/useLoading'
import BaseCard from '~/components/BaseCard.vue'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const { showLoading, hideLoading } = useLoading()
const loading = ref(true)
const spaceId = route.params.id as string

const form = ref({
  name: '',
  description: '',
  type: '',
  base_currency: '',
  start_date: null as string | null,
  end_date: null as string | null,
  is_pinned: false,
  currencies: [] as string[],
  categories: [] as string[],
  payment_methods: [] as string[]
})

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
      try {
        const date = new Date(d)
        if (isNaN(date.getTime())) return null
        const iso = date.toISOString()
        return iso ? (iso.split('T')[0] || null) : null
      } catch {
        return null
      }
    }

    form.value = {
      name: s.name || '',
      description: s.description || '',
      type: s.type || '',
      base_currency: s.base_currency || '',
      start_date: formatDate(s.start_date),
      end_date: formatDate(s.end_date),
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
  showLoading()
  try {
    await api.patch(`/api/spaces/${spaceId}`, {
      name: form.value.name,
      description: form.value.description,
      start_date: form.value.start_date ? new Date(form.value.start_date) : null,
      end_date: form.value.end_date ? new Date(form.value.end_date) : null,
      is_pinned: form.value.is_pinned
    })
    // Refresh store so other pages see updated data
    await detailStore.fetchSpace(true)
    alert('設定已儲存')
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    hideLoading()
  }
}

const handleUpdateLists = async () => {
  showLoading()
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
    hideLoading()
  }
}

const openAliasModal = (member: any) => {
    selectedMember.value = member
    aliasValue.value = member.alias || ''
    showAliasModal.value = true
}

const handleUpdateAlias = async () => {
    if (!selectedMember.value) return
    showLoading()
    try {
        await api.put(`/api/spaces/${spaceId}/members/${selectedMember.value.user_id}`, {
            alias: aliasValue.value
        })
        showAliasModal.value = false
        await detailStore.fetchMembers(true)
    } catch (e: any) {
        alert(e.message || '更新失敗')
    } finally {
        hideLoading()
    }
}

const handleCreateInvite = async () => {
    showLoading()
    try {
        await api.post(`/api/spaces/${spaceId}/invites`, {
            is_one_time: inviteForm.value.is_one_time
        })
        showInviteModal.value = false
        await detailStore.fetchInvites(true)
    } catch (e: any) {
        alert(e.message || '建立失敗')
    } finally {
        hideLoading()
    }
}

const handleRemoveMember = async (member: any) => {
  if (!confirm(`確定要移除成員 ${member.alias || member.user?.display_name} 嗎？`)) return
  try {
    await api.del(`/api/spaces/${spaceId}/members/${member.user_id}`)
    await detailStore.fetchMembers(true)
  } catch (e: any) {
    alert(e.message || '移除失敗')
  }
}

const handleDeleteSpace = async () => {
    if (!confirm('警告：確定要刪除此空間嗎？這將會永久刪除所有交易、比價資料與成員關聯，且無法復原。')) return
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
  alert('連結已複製')
}

const handleRevokeInvite = async (invite_id: string) => {
  if (!confirm('確定要撤銷此邀請連結嗎？')) return
  try {
    await api.del(`/api/spaces/${invite_id}`)
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
