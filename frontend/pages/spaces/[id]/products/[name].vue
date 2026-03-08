<template>
  <div class="product-view-page">
    <PageTitle
      :title="productName"
      :show-back="true"
      :back-to="`/spaces/${route.params.id}/products`"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}` }, { label: '商品比價', to: `/spaces/${route.params.id}/products` }]"
    />

    <div v-if="detailStore.loading.products" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="matchingProducts.length === 0" class="flex flex-col items-center justify-center py-20 text-neutral-500 italic">
      找不到相關商品
    </div>

    <div v-else class="flex flex-col gap-6">
      <div class="px-1 flex items-center justify-between">
        <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest">價格對比 ({{ matchingProducts.length }} 間店)</h2>
      </div>

      <div class="flex flex-col gap-3">
        <div 
          v-for="(p, index) in sortedProducts" 
          :key="p.id"
          @click="router.push(`/spaces/${route.params.id}/stores/${p.store_id}`)"
          class="bg-neutral-900 border border-neutral-800 p-5 rounded-2xl flex flex-col gap-4 cursor-pointer hover:bg-neutral-800 transition-colors relative overflow-hidden"
        >
          <!-- Ranking Ribbon for the cheapest -->
          <div v-if="index === 0" class="absolute top-0 right-0">
             <div class="bg-emerald-500 text-white text-[10px] font-black px-3 py-1 uppercase tracking-tighter rounded-bl-xl shadow-sm">
               最便宜
             </div>
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-neutral-400 border border-neutral-700">
                <Icon icon="mdi:store-outline" class="text-xl" />
              </div>
              <div class="flex flex-col">
                <span class="font-bold text-white">{{ p.store?.name || '未知商店' }}</span>
                <span v-if="p.unit" class="text-xs text-neutral-500 font-medium">單位：{{ p.unit }}</span>
              </div>
            </div>
            
            <div class="text-right">
              <div class="text-xl font-black text-indigo-400">
                <span class="text-xs font-bold mr-0.5">{{ p.currency }}</span>
                {{ formatPrice(p.price) }}
              </div>
            </div>
          </div>

          <!-- Note if exists -->
          <div v-if="p.note" class="bg-neutral-800/50 p-3 rounded-xl text-xs text-neutral-400 italic border border-neutral-800/50">
             "{{ p.note }}"
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'

const route = useRoute()
const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const productName = computed(() => {
  const name = route.params.name
  return typeof name === 'string' ? decodeURIComponent(name) : '商品詳情'
})

const matchingProducts = computed(() => {
  return detailStore.products.filter(p => p.name === productName.value)
})

const sortedProducts = computed(() => {
  return [...matchingProducts.value].sort((a, b) => parseFloat(a.price) - parseFloat(b.price))
})

const formatPrice = (price: number | string) => {
  const num = typeof price === 'string' ? parseFloat(price) : price
  return num.toLocaleString('zh-TW', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()
  await detailStore.fetchProducts()
})
</script>
