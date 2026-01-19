<template>
  <div class="trips-page">
    <header class="page-header">
      <h1>我的旅行</h1>
      <NuxtLink to="/trips/new" class="btn btn-primary">
        <Icon icon="mdi:plus" /> 新增
      </NuxtLink>
    </header>

    <div class="coming-soon card">
      <Icon icon="mdi:airplane-takeoff" />
      <h2>即將推出</h2>
      <p>旅行記帳功能正在開發中...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
  }
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
}

.page-header .btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
}

.coming-soon {
  text-align: center;
  padding: 60px 20px;
}

.coming-soon svg {
  font-size: 64px;
  color: var(--primary);
  margin-bottom: 16px;
}

.coming-soon h2 {
  margin-bottom: 8px;
}

.coming-soon p {
  color: var(--text-secondary);
}
</style>
