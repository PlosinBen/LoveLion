<template>
  <div class="stores-page">
    <PageTitle
      title="商店"
      :show-back="true"
      :settings-to="`/spaces/${route.params.id}/settings`"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}/ledger` }]"
    />

    <!-- View Toggle via URL -->
    <BaseSegmentControl 
      :options="[
        { label: '按商店', to: `/spaces/${route.params.id}/stores` },
        { label: '按商品', to: `/spaces/${route.params.id}/products` }
      ]"
    />

    <div v-if="detailStore.loading.stores" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <!-- Stores View -->
    <div v-if="detailStore.stores.length === 0" class="flex flex-col items-center justify-center py-20 bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed text-neutral-500">
      <Icon icon="mdi:store-plus-outline" class="text-5xl mb-4 opacity-20" />
      <p class="text-sm">尚未建立任何商店紀錄</p>
      <BaseButton @click="router.push(`/spaces/${route.params.id}/stores/add`)" variant="primary" class="mt-6">立即新增</BaseButton>
    </div>

    <div v-else class="flex flex-col gap-4">
      <BaseCard
        v-for="s in detailStore.stores"
        :key="s.id"
        @click="router.push(`/spaces/${route.params.id}/stores/${s.id}`)"
        class="flex items-center justify-between cursor-pointer hover:bg-neutral-800 transition-colors"
      >
        <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-400 border border-neutral-700">
              <Icon icon="mdi:store-outline" class="text-2xl" />
            </div>
            <div class="flex flex-col">
              <span class="font-bold text-white">{{ s.name }}</span>
              <span class="text-xs text-neutral-500 font-medium mt-0.5">{{ s.products?.length || 0 }} 個商品</span>
            </div>
        </div>
        <Icon icon="mdi:chevron-right" class="text-neutral-700 text-xl" />
      </BaseCard>
    </div>

    <!-- FAB for adding store -->
    <BaseFab @click="router.push(`/spaces/${route.params.id}/stores/add`)" />
    </div>
    </template>

<script setup lang="ts">
import BaseButton from '~/components/BaseButton.vue'
import { onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import BaseFab from '~/components/BaseFab.vue'
import BaseSegmentControl from '~/components/BaseSegmentControl.vue'
import BaseCard from '~/components/BaseCard.vue'

definePageMeta({
  layout: 'default'
})

const route = useRoute()
const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  detailStore.setSpaceId(route.params.id as string)
  await Promise.all([
    detailStore.fetchSpace(), 
    detailStore.fetchStores()
  ])
})
</script>
