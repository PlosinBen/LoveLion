<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <LoadingOverlay />
</template>

<script setup lang="ts">
if (import.meta.client && 'serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js')
}
</script>
