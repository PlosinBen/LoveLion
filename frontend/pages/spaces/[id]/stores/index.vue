<template>
  <div class="stores-page">
    <SpaceHeader
      :title="detailStore.space?.name || '比價清單'"
      :show-back="true"
      :settings-to="`/spaces/${route.params.id}/settings`"
    />

    <!-- View Toggle -->
    <div class="flex p-1 bg-neutral-800 rounded-xl mb-6">
      <button 
        @click="viewMode = 'stores'"
        class="flex-1 py-2 text-sm font-bold rounded-lg transition-all border-0 cursor-pointer"
        :class="viewMode === 'stores' ? 'bg-indigo-500 text-white shadow-sm' : 'text-neutral-500 hover:text-neutral-300 bg-transparent'"
      >
        按商店
      </button>
      <button 
        @click="viewMode = 'products'"
        class="flex-1 py-2 text-sm font-bold rounded-lg transition-all border-0 cursor-pointer"
        :class="viewMode === 'products' ? 'bg-indigo-500 text-white shadow-sm' : 'text-neutral-500 hover:text-neutral-300 bg-transparent'"
      >
        按商品
      </button>
    </div>

    <div v-if="detailStore.loading.stores || detailStore.loading.products" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <!-- Stores View -->
    <template v-else-if="viewMode === 'stores'">
      <div v-if="detailStore.stores.length === 0" class="flex flex-col items-center justify-center py-20 bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed text-neutral-500">
        <Icon icon="mdi:store-plus-outline" class="text-5xl mb-4 opacity-20" />
        <p class="text-sm">尚未建立任何商店紀錄</p>
        <button @click="router.push(`/spaces/${route.params.id}/stores/add`)" class="mt-6 px-6 py-2 bg-indigo-500 text-white rounded-full font-bold text-sm border-0 cursor-pointer">立即新增</button>
      </div>

      <div v-else class="flex flex-col gap-4">
        <div
          v-for="s in detailStore.stores"
          :key="s.id"
          @click="router.push(`/spaces/${route.params.id}/stores/${s.id}`)"
          class="bg-neutral-900 border border-neutral-800 p-4 rounded-2xl flex items-center justify-between cursor-pointer hover:bg-neutral-800 transition-colors"
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
        </div>
      </div>
    </template>

    <!-- Products View -->
    <template v-else-if="viewMode === 'products'">
      <div v-if="groupedProducts.length === 0" class="flex flex-col items-center justify-center py-20 bg-neutral-900 rounded-2xl border border-neutral-800 border-dashed text-neutral-500">
        <Icon icon="mdi:package-variant" class="text-5xl mb-4 opacity-20" />
        <p class="text-sm">尚未建立任何商品紀錄</p>
      </div>

      <div v-else class="flex flex-col gap-4">
        <div
          v-for="p in groupedProducts"
          :key="p.name"
          @click="router.push({ path: `/spaces/${route.params.id}/products/view`, query: { name: p.name } })"
          class="bg-neutral-900 border border-neutral-800 p-4 rounded-2xl flex items-center justify-between cursor-pointer hover:bg-neutral-800 transition-colors"
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
        </div>
      </div>
    </template>

    <!-- FAB for adding store -->
    <button
      @click="router.push(`/spaces/${route.params.id}/stores/add`)"
      class="fixed bottom-32 right-6 w-14 h-14 bg-indigo-500 shadow-lg rounded-full flex items-center justify-center text-white z-20 cursor-pointer border-0 transition-transform active:scale-95"
    >
      <Icon icon="mdi:plus" class="text-3xl" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import SpaceHeader from '~/components/SpaceHeader.vue'

definePageMeta({
  layout: 'default'
})

const route = useRoute()
const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const viewMode = ref<'stores' | 'products'>('stores')

const formatPrice = (price: number | string) => {
  const num = typeof price === 'string' ? parseFloat(price) : price
  return num.toLocaleString('zh-TW', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

const groupedProducts = computed(() => {
  const groups: Record<string, any[]> = {}
  detailStore.products.forEach(p => {
    if (!groups[p.name]) groups[p.name] = []
    groups[p.name].push(p)
  })

  return Object.keys(groups).map(name => {
    const prods = groups[name]
    const prices = prods.map(p => parseFloat(p.price))
    const min = Math.min(...prices)
    const max = Math.max(...prices)
    const currency = prods[0].currency

    let priceRange = `${currency} ${formatPrice(min)}`
    if (min !== max) {
      priceRange += ` ~ ${formatPrice(max)}`
    }

    return {
      name,
      count: prods.length,
      priceRange
    }
  }).sort((a, b) => a.name.localeCompare(b.name))
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
    detailStore.fetchStores(),
    detailStore.fetchProducts()
  ])
})
</script>
