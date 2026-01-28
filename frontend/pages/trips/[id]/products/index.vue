<template>
  <div class="flex flex-col gap-4">
    <header class="flex-none pb-2 bg-neutral-900/90 backdrop-blur-md sticky top-0 z-10 -mx-4 px-4 pt-2 border-b border-white/5">
      <div class="flex justify-between items-center mb-4">
        <button @click="router.push(`/trips/${route.params.id}`)" class="flex justify-center items-center w-10 h-10 rounded-full bg-neutral-800 text-white border-0 cursor-pointer hover:bg-neutral-700 transition-colors">
          <Icon icon="mdi:arrow-left" class="text-xl" />
        </button>
        <h1 class="text-xl font-bold">比價清單</h1>
        <div class="w-10"></div> <!-- Spacer -->
      </div>

      <!-- Tabs -->
      <div class="flex p-1 bg-neutral-800 rounded-xl relative">
         <div class="absolute inset-y-1 w-[calc(50%-4px)] bg-neutral-700 rounded-lg shadow-sm transition-all duration-300 left-[50%]"></div>
         <button @click="router.push(`/trips/${route.params.id}/stores`)" 
                 class="flex-1 relative z-10 py-2 text-sm font-medium transition-colors text-center border-0 bg-transparent cursor-pointer text-neutral-400 hover:text-neutral-200">
             商店列表
         </button>
         <button class="flex-1 relative z-10 py-2 text-sm font-medium transition-colors text-center border-0 bg-transparent cursor-pointer text-white">
             商品比價
         </button>
      </div>
    </header>

    <div v-if="loading" class="flex justify-center items-center text-neutral-400 py-10">
       <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <!-- Mode: Products -->
    <div v-else class="">
        <div v-if="productGroups.length === 0" class="flex flex-col items-center justify-center py-20 text-neutral-500">
            <Icon icon="mdi:shopping-search" class="text-6xl mb-4 opacity-20" />
            <p>尚無可比價的商品</p>
        </div>

        <div v-else class="flex flex-col gap-3">
            <div v-for="group in productGroups" :key="group.name" 
                 class="bg-neutral-900 rounded-xl border border-neutral-800 overflow-hidden">
                 <!-- Header -->
                 <div class="p-4 flex items-center gap-3 cursor-pointer" @click="group.expanded = !group.expanded">
                     <div class="w-12 h-12 rounded-lg bg-indigo-500/10 flex items-center justify-center text-indigo-400 font-bold">
                         <Icon icon="mdi:tag-outline" class="text-2xl" />
                     </div>
                     <div class="flex-1">
                         <h3 class="font-bold text-base">{{ group.name }}</h3>
                         <div class="text-xs text-neutral-400 mt-1 flex gap-2">
                             <span>{{ group.items.length }} 家販售</span>
                             <span class="text-indigo-400 font-medium">
                                 {{ group.currency }} 
                                 <span v-if="group.minPrice !== group.maxPrice">{{ formatPrice(group.minPrice) }} ~ {{ formatPrice(group.maxPrice) }}</span>
                                 <span v-else>{{ formatPrice(group.minPrice) }}</span>
                             </span>
                         </div>
                     </div>
                     <Icon :icon="group.expanded ? 'mdi:chevron-up' : 'mdi:chevron-down'" class="text-neutral-500 text-xl" />
                 </div>

                 <!-- Expanded List -->
                 <div v-if="group.expanded" class="border-t border-neutral-800 bg-neutral-800/20">
                     <div v-for="item in group.items" :key="item.id" 
                          class="p-3 py-2 flex items-center justify-between border-b border-neutral-800/50 last:border-0 hover:bg-white/5 cursor-pointer"
                          @click="router.push(`/trips/${route.params.id}/stores/${item.store_id}`)">
                         <div class="flex items-center gap-2">
                             <div class="text-xs text-neutral-300">{{ item.store?.name }}</div>
                             <span v-if="item.price == group.minPrice" class="bg-red-500/20 text-red-500 text-[10px] px-1.5 py-0.5 rounded font-bold">最低</span>
                         </div>
                         <div class="text-sm font-medium">{{ item.currency }} {{ formatPrice(item.price) }}</div>
                     </div>
                 </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const productGroups = ref<any[]>([])
const loading = ref(true)

const fetchData = async () => {
  try {
    const products = await api.get<any[]>(`/api/trips/${route.params.id}/products`)
    
    // Group products
    const groups: Record<string, any> = {}
    
    products.forEach(p => {
        const nameKey = p.name.trim().toLowerCase()
        if (!groups[nameKey]) {
            groups[nameKey] = {
                name: p.name,
                currency: p.currency,
                items: [],
                expanded: false,
                minPrice: Infinity,
                maxPrice: -Infinity
            }
        }
        groups[nameKey].items.push(p)
        groups[nameKey].minPrice = Math.min(groups[nameKey].minPrice, p.price)
        groups[nameKey].maxPrice = Math.max(groups[nameKey].maxPrice, p.price)
    })

    productGroups.value = Object.values(groups).map(g => {
        // Sort items by price
        g.items.sort((a: any, b: any) => parseFloat(a.price) - parseFloat(b.price))
        return g
    })

  } catch (e) {
    console.error('Failed to fetch data:', e)
  } finally {
    loading.value = false
  }
}

const formatPrice = (price: number | string) => {
  const num = typeof price === 'string' ? parseFloat(price) : price
  return num.toLocaleString('zh-TW', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchData()
})
</script>
