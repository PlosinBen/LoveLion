<template>
  <div class="flex flex-col justify-center min-h-screen p-6 bg-neutral-950">
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold mb-2 bg-gradient-to-br from-indigo-500 to-purple-500 bg-clip-text text-transparent">LoveLion</h1>
      <p class="text-neutral-400">個人記帳 & 旅行助手</p>
    </div>

    <div class="w-full max-w-sm mx-auto bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
      <div v-if="!isRegister">
        <h2 class="mb-6 text-center text-xl font-semibold">登入</h2>

        <div class="mb-4">
          <BaseInput v-model="username" label="帳號" placeholder="請輸入帳號" />
        </div>

        <div class="mb-4">
          <BaseInput v-model="password" type="password" label="密碼" placeholder="請輸入密碼" />
        </div>

        <div v-if="error" class="text-red-500 text-sm mt-2">{{ error }}</div>

        <button @click="handleLogin" class="w-full mt-6 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="loading">
          {{ loading ? '登入中...' : '登入' }}
        </button>

        <p class="text-center mt-4 text-neutral-400">
          還沒有帳號？ <a @click="isRegister = true" class="text-indigo-500 cursor-pointer hover:underline">註冊</a>
        </p>
      </div>

      <div v-else>
        <h2 class="mb-6 text-center text-xl font-semibold">註冊</h2>

        <div class="mb-4">
          <BaseInput v-model="username" label="帳號" placeholder="請輸入帳號" />
        </div>

        <div class="mb-4">
          <BaseInput v-model="displayName" label="顯示名稱" placeholder="請輸入顯示名稱" />
        </div>

        <div class="mb-4">
          <BaseInput v-model="password" type="password" label="密碼" placeholder="請輸入密碼" />
        </div>

        <div v-if="error" class="text-red-500 text-sm mt-2">{{ error }}</div>

        <button @click="handleRegister" class="w-full mt-6 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="loading">
          {{ loading ? '註冊中...' : '註冊' }}
        </button>

        <p class="text-center mt-4 text-neutral-400">
          已有帳號？ <a @click="isRegister = false" class="text-indigo-500 cursor-pointer hover:underline">登入</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '~/composables/useAuth'

definePageMeta({
  layout: false
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
    error.value = '請填寫帳號和密碼'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await login(username.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || '登入失敗'
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
