<template>
  <div class="settings-page">
    <PageTitle
      title="個人設定"
      :show-back="true"
      :breadcrumbs="[{ label: '空間列表', to: '/' }]"
    />

    <div class="flex flex-col gap-8 pb-20">
      <!-- Profile Section -->
      <section class="flex flex-col gap-3">
        <label class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">帳戶資訊</label>
        <BaseCard padding="p-6" class="flex flex-col gap-6 shadow-sm">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
              <div class="w-14 h-14 rounded-2xl bg-neutral-800 flex items-center justify-center text-indigo-500 border border-neutral-700/50 shadow-inner">
                <Icon icon="mdi:account-outline" class="text-3xl" />
              </div>
              <div class="flex flex-col">
                <span class="text-lg font-bold text-white tracking-tight">{{ user?.display_name || user?.username || '使用者' }}</span>
                <span class="text-xs text-neutral-500 font-medium tracking-wide">@{{ user?.username }}</span>
              </div>
            </div>
            <button 
              @click="openEditProfile"
              class="w-10 h-10 rounded-full bg-neutral-800 flex items-center justify-center text-neutral-400 hover:text-white hover:bg-neutral-700 transition-all active:scale-95"
            >
              <Icon icon="mdi:pencil" class="text-xl" />
            </button>
          </div>
        </BaseCard>
      </section>

      <!-- Password Section -->
      <section class="flex flex-col gap-3">
        <label class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">安全設定</label>
        <BaseCard padding="p-6" class="flex flex-col gap-5 shadow-sm">
          <div class="flex flex-col gap-1">
            <h3 class="text-sm font-bold text-white">修改密碼</h3>
            <p class="text-xs text-neutral-500 font-medium">定期更換密碼可保護帳戶安全</p>
          </div>
          
          <div class="flex flex-col gap-4 mt-2">
            <BaseInput
              v-model="passwordForm.current"
              label="目前密碼"
              type="password"
              placeholder="請輸入目前密碼"
            />
            <BaseInput
              v-model="passwordForm.new"
              label="新密碼"
              type="password"
              placeholder="請輸入新密碼 (至少 6 個字)"
            />
            <BaseInput
              v-model="passwordForm.confirm"
              label="確認新密碼"
              type="password"
              placeholder="請再次輸入新密碼"
            />
            
            <div v-if="passwordError" class="text-xs text-rose-500 font-bold px-1 animate-pulse">
              {{ passwordError }}
            </div>
            <div v-if="passwordSuccess" class="text-xs text-emerald-500 font-bold px-1">
              密碼修改成功！
            </div>

            <BaseButton 
              @click="handleUpdatePassword" 
              class="mt-2"
              :disabled="isUpdatingPassword || !passwordForm.current || !passwordForm.new"
            >
              {{ isUpdatingPassword ? '更新中...' : '確認修改' }}
            </BaseButton>
          </div>
        </BaseCard>
      </section>

      <!-- Space Management Section -->
      <section class="flex flex-col gap-3">
        <label class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">空間管理</label>
        <div class="flex flex-col gap-3">
          <BaseCard 
            v-for="space in allSpaces" 
            :key="space.id" 
            padding="p-4" 
            class="flex items-center justify-between shadow-sm group"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-neutral-800 flex items-center justify-center text-indigo-400 border border-neutral-700/50">
                <Icon :icon="space.type === 'trip' ? 'mdi:map-marker-outline' : 'mdi:home-outline'" class="text-xl" />
              </div>
              <div class="flex flex-col">
                <span class="text-sm font-bold text-white">{{ space.name }}</span>
                <span class="text-[10px] text-neutral-500 uppercase font-bold tracking-wider">{{ space.type === 'trip' ? '旅行' : '個人' }}</span>
              </div>
            </div>
            
            <button 
              @click="confirmLeaveSpace(space)"
              class="px-3 py-1.5 rounded-lg text-xs font-bold bg-rose-500/10 text-rose-500 border border-rose-500/20 hover:bg-rose-500 hover:text-white transition-all active:scale-95"
            >
              退出
            </button>
          </BaseCard>

          <div v-if="allSpaces.length === 0" class="py-10 flex flex-col items-center justify-center text-neutral-600 bg-neutral-900/50 rounded-3xl border-2 border-dashed border-neutral-800">
            <Icon icon="mdi:folder-outline" class="text-4xl mb-2 opacity-20" />
            <span class="text-xs font-bold uppercase tracking-widest">目前沒有加入任何空間</span>
          </div>
        </div>
      </section>

      <!-- Logout Section -->
      <section class="mt-4">
        <BaseButton @click="handleLogout" variant="danger" class="w-full h-14 rounded-2xl">
          登出帳戶
        </BaseButton>
      </section>

      <!-- App Info -->
      <p class="text-center text-[10px] text-neutral-700 font-bold uppercase tracking-[0.2em] mt-6 flex flex-col gap-2">
        <span>LoveLion v1.2.0-stable</span>
        <span>Crafted with Passion © 2026 Antigravity</span>
      </p>
    </div>

    <!-- Edit Profile Modal -->
    <BaseModal
      v-model="isEditProfileOpen"
      title="編輯基本資料"
    >
      <div class="flex flex-col gap-4 p-1">
        <BaseInput
          v-model="editForm.displayName"
          label="顯示名稱"
          placeholder="您希望如何被稱呼？"
        />
        <div class="flex justify-end gap-3 mt-4">
          <BaseButton variant="secondary" @click="isEditProfileOpen = false">取消</BaseButton>
          <BaseButton @click="handleUpdateProfile" :disabled="isUpdatingProfile">
            {{ isUpdatingProfile ? '儲存中...' : '儲存修改' }}
          </BaseButton>
        </div>
      </div>
    </BaseModal>

    <!-- Confirm Leave Space Modal -->
    <BaseModal
      v-model="isLeaveModalOpen"
      title="確認退出空間"
    >
      <div class="flex flex-col gap-4 p-1 text-center">
        <div class="w-16 h-16 rounded-full bg-rose-500/10 text-rose-500 flex items-center justify-center mx-auto mb-2 border border-rose-500/20">
          <Icon icon="mdi:alert-outline" class="text-4xl" />
        </div>
        <p class="text-sm text-neutral-300">
          您確定要退出 <span class="text-white font-bold">"{{ spaceToLeave?.name }}"</span> 嗎？<br>
          退出後，您將無法再存取此空間的所有帳務資料。
        </p>
        <div class="flex flex-col gap-2 mt-4">
          <BaseButton variant="danger" @click="handleLeaveSpace" :disabled="isLeavingSpace">
            {{ isLeavingSpace ? '處理中...' : '確認退出' }}
          </BaseButton>
          <BaseButton variant="secondary" @click="isLeaveModalOpen = false">取消</BaseButton>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useSpace } from '~/composables/useSpace'
import PageTitle from '~/components/PageTitle.vue'
import BaseButton from '~/components/BaseButton.vue'
import BaseCard from '~/components/BaseCard.vue'
import BaseInput from '~/components/BaseInput.vue'
import BaseModal from '~/components/BaseModal.vue'

definePageMeta({
  layout: 'default'
})

const router = useRouter()
const { user, logout: authLogout, initAuth, isAuthenticated, updateProfile } = useAuth()
const { allSpaces, fetchSpaces, leaveSpace } = useSpace()

// Profile Edit Logic
const isEditProfileOpen = ref(false)
const isUpdatingProfile = ref(false)
const editForm = reactive({
  displayName: ''
})

const openEditProfile = () => {
  editForm.displayName = user.value?.display_name || ''
  isEditProfileOpen.value = true
}

const handleUpdateProfile = async () => {
  if (!editForm.displayName.trim()) return
  
  isUpdatingProfile.value = true
  try {
    await updateProfile({ display_name: editForm.displayName })
    isEditProfileOpen.value = false
  } catch (e: any) {
    alert(e.error || '更新失敗')
  } finally {
    isUpdatingProfile.value = false
  }
}

// Password Edit Logic
const isUpdatingPassword = ref(false)
const passwordError = ref('')
const passwordSuccess = ref(false)
const passwordForm = reactive({
  current: '',
  new: '',
  confirm: ''
})

const handleUpdatePassword = async () => {
  passwordError.value = ''
  passwordSuccess.value = false

  if (passwordForm.new !== passwordForm.confirm) {
    passwordError.value = '新密碼與確認密碼不符'
    return
  }
  if (passwordForm.new.length < 6) {
    passwordError.value = '新密碼至少需要 6 個字元'
    return
  }

  isUpdatingPassword.value = true
  try {
    await updateProfile({
      current_password: passwordForm.current,
      new_password: passwordForm.new
    })
    passwordSuccess.value = true
    passwordForm.current = ''
    passwordForm.new = ''
    passwordForm.confirm = ''
    
    // Success message auto-hide
    setTimeout(() => {
      passwordSuccess.value = false
    }, 3000)
  } catch (e: any) {
    passwordError.value = e.error || '密碼更新失敗'
  } finally {
    isUpdatingPassword.value = false
  }
}

// Leave Space Logic
const isLeaveModalOpen = ref(false)
const isLeavingSpace = ref(false)
const spaceToLeave = ref<any>(null)

const confirmLeaveSpace = (space: any) => {
  spaceToLeave.value = space
  isLeaveModalOpen.value = true
}

const handleLeaveSpace = async () => {
  if (!spaceToLeave.value) return
  
  isLeavingSpace.value = true
  try {
    await leaveSpace(spaceToLeave.value.id)
    isLeaveModalOpen.value = false
    spaceToLeave.value = null
  } catch (e: any) {
    alert(e.error || '無法退出空間')
  } finally {
    isLeavingSpace.value = false
  }
}

const handleLogout = () => {
  authLogout()
  router.push('/login')
}

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  await fetchSpaces(true)
})
</script>

<style scoped>
.settings-page {
  @apply max-w-lg mx-auto;
}
</style>
