<template>
  <div class="login-page">
    <div class="login-header">
      <h1>LoveLion</h1>
      <p>個人記帳 & 旅行助手</p>
    </div>

    <div class="login-form card">
      <div v-if="!isRegister">
        <h2>登入</h2>

        <div class="form-group">
          <label class="label">帳號</label>
          <input v-model="username" type="text" class="input" placeholder="請輸入帳號" />
        </div>

        <div class="form-group">
          <label class="label">密碼</label>
          <input v-model="password" type="password" class="input" placeholder="請輸入密碼" />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <button @click="handleLogin" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '登入中...' : '登入' }}
        </button>

        <p class="switch-mode">
          還沒有帳號？ <a @click="isRegister = true">註冊</a>
        </p>
      </div>

      <div v-else>
        <h2>註冊</h2>

        <div class="form-group">
          <label class="label">帳號</label>
          <input v-model="username" type="text" class="input" placeholder="請輸入帳號" />
        </div>

        <div class="form-group">
          <label class="label">顯示名稱</label>
          <input v-model="displayName" type="text" class="input" placeholder="請輸入顯示名稱" />
        </div>

        <div class="form-group">
          <label class="label">密碼</label>
          <input v-model="password" type="password" class="input" placeholder="請輸入密碼" />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <button @click="handleRegister" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '註冊中...' : '註冊' }}
        </button>

        <p class="switch-mode">
          已有帳號？ <a @click="isRegister = false">登入</a>
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

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 24px;
  background: var(--bg-primary);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 32px;
  background: linear-gradient(135deg, var(--primary), #a855f7);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 8px;
}

.login-header p {
  color: var(--text-secondary);
}

.login-form {
  max-width: 400px;
  margin: 0 auto;
  width: 100%;
}

.login-form h2 {
  margin-bottom: 24px;
  text-align: center;
}

.form-group {
  margin-bottom: 16px;
}

.btn-block {
  width: 100%;
  margin-top: 24px;
}

.switch-mode {
  text-align: center;
  margin-top: 16px;
  color: var(--text-secondary);
}

.switch-mode a {
  color: var(--primary);
  cursor: pointer;
}

.error-message {
  color: var(--danger);
  font-size: 14px;
  margin-top: 8px;
}
</style>
