<template>
  <div class="bottom-nav">
    <NuxtLink to="/" class="nav-item" :class="{ active: route.path === '/' }">
      <Icon icon="mdi:home" />
      <span>首頁</span>
    </NuxtLink>
    <NuxtLink to="/ledger" class="nav-item" :class="{ active: route.path.startsWith('/ledger') }">
      <Icon icon="mdi:wallet" />
      <span>記帳</span>
    </NuxtLink>
    <NuxtLink to="/trips" class="nav-item" :class="{ active: route.path.startsWith('/trips') }">
      <Icon icon="mdi:airplane" />
      <span>旅行</span>
    </NuxtLink>
    <NuxtLink to="/settings" class="nav-item" :class="{ active: route.path === '/settings' }">
      <Icon icon="mdi:cog" />
      <span>設定</span>
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

const route = useRoute()
</script>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: space-around;
  align-items: center;
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 12px;
  padding: 8px 16px;
  transition: color 0.2s;
}

.nav-item svg {
  font-size: 24px;
}

.nav-item.active {
  color: var(--primary);
}

.nav-item:hover {
  color: var(--primary);
}
</style>
