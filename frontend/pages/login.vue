<template>
  <div class="flex flex-col justify-center py-6">
    <div class="text-center mb-8">
      <p class="text-neutral-400 font-medium">個人記帳 & 旅遊專案助手</p>
    </div>

    <BaseCard padding="p-6" class="w-full max-w-sm mx-auto shadow-xl">
      <div v-if="!isRegister">
        <h2 class="mb-8 text-center text-2xl font-bold text-white tracking-tight">歡迎回來</h2>

        <div class="flex flex-col gap-4">
          <BaseInput 
            v-model="username" 
            label="帳號" 
            placeholder="請輸入您的帳號" 
            required
          />

          <BaseInput 
            v-model="password" 
            type="password" 
            label="密碼" 
            placeholder="請輸入您的密碼" 
            required
          />
        </div>

        <div v-if="error" class="text-red-500 text-xs mt-4 font-bold bg-red-500/10 p-3 rounded-xl border border-red-500/20">
          {{ error }}
        </div>

        <BaseButton 
          @click="handleLogin" 
          class="mt-8"
          fullWidth
          :loading="loading"
        >
          登入
        </BaseButton>

        <p class="text-center mt-6 text-neutral-500 text-sm font-medium flex items-center justify-center gap-1">
          還沒有帳號嗎？ 
          <BaseButton @click="isRegister = true" variant="ghost" class="p-0 font-bold hover:underline h-auto inline-flex">立即註冊</BaseButton>
        </p>
      </div>

      <div v-else>
        <h2 class="mb-8 text-center text-2xl font-bold text-white tracking-tight">加入 LoveLion</h2>

        <div class="flex flex-col gap-4">
          <BaseInput 
            v-model="username" 
            label="帳號" 
            placeholder="設定登入帳號" 
            required
          />

          <BaseInput 
            v-model="displayName" 
            label="顯示名稱" 
            placeholder="大家如何稱呼您" 
            required
          />

          <BaseInput 
            v-model="password" 
            type="password" 
            label="密碼" 
            placeholder="設定登入密碼" 
            required
          />
        </div>

        <div v-if="error" class="text-red-500 text-xs mt-4 font-bold bg-red-500/10 p-3 rounded-xl border border-red-500/20">
          {{ error }}
        </div>

        <BaseButton 
          @click="handleRegister" 
          class="mt-8"
          fullWidth
          :loading="loading"
        >
          註冊帳號
        </BaseButton>

        <p class="text-center mt-6 text-neutral-500 text-sm font-medium flex items-center justify-center gap-1">
          已經有帳號了？ 
          <BaseButton @click="isRegister = false" variant="ghost" class="p-0 font-bold hover:underline h-auto inline-flex">點此登入</BaseButton>
        </p>
      </div>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '~/composables/useAuth'
import BaseButton from '~/components/BaseButton.vue'
import BaseCard from '~/components/BaseCard.vue'

definePageMeta({
  layout: 'default',
  hideGlobalNav: true
})

const router = useRouter()
const { login, register } = useAuth()

const isRegister = ref(false)
const username = ref('')
const password = ref('')
const displayName = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!username.value || !password.value) {
    error.value = '請填寫帳號與密碼'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await login(username.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || '登入失敗，請檢查帳號密碼'
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  if (!username.value || !password.value || !displayName.value) {
    error.value = '請填寫所有欄位'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await register(username.value, password.value, displayName.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || '註冊失敗'
  } finally {
    loading.value = false
  }
}
</script>
