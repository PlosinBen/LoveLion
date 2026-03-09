<template>
  <div class="products-index-page">
    <PageTitle
      :title="detailStore.space?.name || '商品比價'"
      :show-back="true"
      :settings-to="`/spaces/${route.params.id}/settings`"
    />

    <!-- View Toggle via URL -->
    <BaseSegmentControl 
      :options="[
        { label: '按商店', to: `/spaces/${route.params.id}/stores` },
        { label: '按商品', to: `/spaces/${route.params.id}/products` }
      ]"
    />

    <div v-if="detailStore.loading.products" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="groupedProducts.length === 0" class="flex flex-col items-center justify-center py-20 bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed text-neutral-500">
      <Icon icon="mdi:package-variant" class="text-5xl mb-4 opacity-20" />
      <p class="text-sm">尚未建立任何商品紀錄</p>
    </div>

    <div v-else class="flex flex-col gap-4">
      <BaseCard
        v-for="p in groupedProducts"
        :key="p.name"
        @click="router.push(`/spaces/${route.params.id}/products/${encodeURIComponent(p.name)}`)"
        class="flex items-center justify-between cursor-pointer hover:bg-neutral-800 transition-colors"
      >
        <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-400 border border-neutral-700">
              <Icon icon="mdi:tag-outline" class="text-2xl" />
            </div>
            <div class="flex flex-col">
              <span class="font-bold text-white">{{ p.name }}</span>
              <div class="flex items-center gap-2 mt-0.5">
                <span class="text-xs text-neutral-500 font-medium">{{ p.count }} 個店面</span>
                <span class="text-xs text-indigo-400 font-bold">{{ p.priceRange }}</span>
              </div>
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
import { onMounted, computed } from 'vue'
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

const formatPrice = (price: number | string) => {
  const num = typeof price === 'string' ? parseFloat(price) : price
  return num.toLocaleString('zh-TW', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

interface GroupedProduct {
  name: string
  count: number
  priceRange: string
}

const groupedProducts = computed<GroupedProduct[]>(() => {
  const groups = new Map<string, any[]>()
  
  detailStore.products.forEach(p => {
    if (!p.name) return
    if (!groups.has(p.name)) {
      groups.set(p.name, [])
    }
    groups.get(p.name)?.push(p)
  })

  const result: GroupedProduct[] = []
  
  groups.forEach((prods, name) => {
    if (!prods || prods.length === 0) return
    
    const prices = prods.map(p => parseFloat(p.price)).filter(p => !isNaN(p))
    if (prices.length === 0) return

    const min = Math.min(...prices)
    const max = Math.max(...prices)
    const currency = prods[0].currency || ''

    let priceRange = `${currency} ${formatPrice(min)}`
    if (min !== max) {
      priceRange += ` ~ ${formatPrice(max)}`
    }

    result.push({
      name,
      count: prods.length,
      priceRange
    })
  })

  return result.sort((a, b) => a.name.localeCompare(b.name))
})

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  detailStore.setSpaceId(route.params.id as string)
  await Promise.all([
    detailStore.fetchSpace(), 
    detailStore.fetchProducts()
  ])
})
</script>
