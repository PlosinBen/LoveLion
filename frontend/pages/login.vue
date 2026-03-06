<template>
  <div class="flex flex-col justify-center py-6">
    <!-- Header is provided by layout 'default' -->
    <div class="text-center mb-8">
      <!-- Removed duplicate logo -->
      <p class="text-neutral-400">?犖閮董 & ???拇?</p>
    </div>

    <div class="w-full max-w-sm mx-auto bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
      <div v-if="!isRegister">
        <h2 class="mb-6 text-center text-xl font-semibold">?餃</h2>

        <div class="mb-4">
          <BaseInput v-model="username" label="撣唾?" placeholder="隢撓?亙董?? />
        </div>

        <div class="mb-4">
          <BaseInput v-model="password" type="password" label="撖Ⅳ" placeholder="隢撓?亙?蝣? />
        </div>

        <div v-if="error" class="text-red-500 text-sm mt-2">{{ error }}</div>

        <button @click="handleLogin" class="w-full mt-6 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="loading">
          {{ loading ? '?餃銝?..' : '?餃' }}
        </button>

        <p class="text-center mt-4 text-neutral-400">
          ???董?? <a @click="isRegister = true" class="text-indigo-500 cursor-pointer hover:underline">閮餃?</a>
        </p>
      </div>

      <div v-else>
        <h2 class="mb-6 text-center text-xl font-semibold">閮餃?</h2>

        <div class="mb-4">
          <BaseInput v-model="username" label="撣唾?" placeholder="隢撓?亙董?? />
        </div>

        <div class="mb-4">
          <BaseInput v-model="displayName" label="憿舐內?迂" placeholder="隢撓?仿＊蝷箏?蝔? />
        </div>

        <div class="mb-4">
          <BaseInput v-model="password" type="password" label="撖Ⅳ" placeholder="隢撓?亙?蝣? />
        </div>

        <div v-if="error" class="text-red-500 text-sm mt-2">{{ error }}</div>

        <button @click="handleRegister" class="w-full mt-6 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="loading">
          {{ loading ? '閮餃?銝?..' : '閮餃?' }}
        </button>

        <p class="text-center mt-4 text-neutral-400">
          撌脫?撣唾?嚗?<a @click="isRegister = false" class="text-indigo-500 cursor-pointer hover:underline">?餃</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
definePageMeta({
  layout: 'default'
})
const { login, register } = useAuth()

const isRegister = ref(false)
const username = ref('')
const password = ref('')
const displayName = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!username.value || !password.value) {
    error.value = '隢‵撖怠董??撖Ⅳ'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await login(username.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || '?餃憭望?'
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  if (!username.value || !password.value || !displayName.value) {
    error.value = '隢‵撖急???雿?
    return
  }

  loading.value = true
  error.value = ''

  try {
    await register(username.value, password.value, displayName.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || '閮餃?憭望?'
  } finally {
    loading.value = false
  }
}
</script>
