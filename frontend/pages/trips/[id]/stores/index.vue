<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.push(`/trips/${route.params.id}`)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">比價商店</h1>
      <button @click="router.push(`/trips/${route.params.id}/stores/create`)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-indigo-500 text-white border-0 cursor-pointer hover:bg-indigo-600 transition-colors">
        <Icon icon="mdi:plus" class="text-2xl" />
      </button>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <div v-else-if="stores.length === 0" class="text-center py-16 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
      <Icon icon="mdi:store-outline" class="text-6xl text-neutral-500 mb-4" />
      <h2 class="text-lg font-semibold mb-2">還沒有商店</h2>
      <p class="text-neutral-400 mb-5">新增商店來開始比較價格</p>
      <button @click="router.push(`/trips/${route.params.id}/stores/create`)" class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm border-0 cursor-pointer">
        新增商店
      </button>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div
        v-for="store in stores"
        :key="store.id"
        class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 cursor-pointer transition-all duration-200 hover:-translate-y-0.5 hover:border-indigo-500"
        @click="router.push(`/trips/${route.params.id}/stores/${store.id}`)"
      >
        <div class="flex items-center gap-3">
          <div class="flex items-center justify-center w-12 h-12 rounded-xl bg-indigo-500/20">
            <Icon icon="mdi:store" class="text-2xl text-indigo-500" />
          </div>
          <div class="flex-1">
            <h3 class="text-base font-semibold mb-1">{{ store.name }}</h3>
            <p class="text-xs text-neutral-400">{{ store.products?.length || 0 }} 個商品</p>
          </div>
          <button @click.stop="deleteStore(store.id)" class="flex justify-center items-center w-8 h-8 rounded-lg bg-red-500/20 text-red-500 border-0 cursor-pointer hover:bg-red-500/30 transition-colors">
            <Icon icon="mdi:delete" />
          </button>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const stores = ref<any[]>([])
const loading = ref(true)

const fetchStores = async () => {
  try {
    stores.value = await api.get<any[]>(`/api/trips/${route.params.id}/stores`)
  } catch (e) {
    console.error('Failed to fetch stores:', e)
  } finally {
    loading.value = false
  }
}



const deleteStore = async (storeId: string) => {
  if (!confirm('確定要刪除此商店？')) return
  try {
    await api.del(`/api/trips/${route.params.id}/stores/${storeId}`)
    fetchStores()
  } catch (e: any) {
    alert(e.message || '刪除失敗')
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchStores()
})
</script>
