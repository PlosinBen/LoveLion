<template>
  <div class="flex flex-col justify-center py-6">
    <div class="text-center mb-8">
      <p class="text-neutral-400 font-medium">個人記帳 & 旅遊專案助手</p>
    </div>

    <div class="w-full max-w-sm mx-auto bg-neutral-900 rounded-2xl p-6 border border-neutral-800 shadow-xl">
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

        <button 
          @click="handleLogin" 
          class="w-full mt-8 py-4 rounded-xl font-bold bg-indigo-500 text-white hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer shadow-lg" 
          :disabled="loading"
        >
          {{ loading ? '登入中...' : '登入' }}
        </button>

        <p class="text-center mt-6 text-neutral-500 text-sm font-medium">
          還沒有帳號嗎？ <button @click="isRegister = true" class="text-indigo-400 bg-transparent border-0 cursor-pointer font-bold hover:underline p-0">立即註冊</button>
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

        <button 
          @click="handleRegister" 
          class="w-full mt-8 py-4 rounded-xl font-bold bg-indigo-500 text-white hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer shadow-lg" 
          :disabled="loading"
        >
          {{ loading ? '註冊中...' : '註冊帳號' }}
        </button>

        <p class="text-center mt-6 text-neutral-500 text-sm font-medium">
          已經有帳號了？ <button @click="isRegister = false" class="text-indigo-400 bg-transparent border-0 cursor-pointer font-bold hover:underline p-0">點此登入</button>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '~/composables/useAuth'

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
