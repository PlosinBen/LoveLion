<template>
  <div class="flex flex-col h-[calc(100vh-64px)] bg-black text-white">
    <!-- Header -->
    <header class="flex-none p-4 flex justify-between items-center bg-neutral-900 border-b border-neutral-800">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-full bg-neutral-800 text-white border-0 cursor-pointer hover:bg-neutral-700 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-xl" />
      </button>
      <h1 class="text-xl font-bold">編輯商店</h1>
      <button @click="saveStore" class="flex justify-center items-center px-4 h-10 rounded-lg bg-indigo-600 text-white border-0 cursor-pointer hover:bg-indigo-500 transition-colors font-medium">
        儲存
      </button>
    </header>

    <div v-if="loading" class="flex-1 flex justify-center items-center text-neutral-400">
      <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex-1 overflow-y-auto p-4 space-y-6">
        <!-- Basic Info -->
        <div class="space-y-4">
            <div class="space-y-2">
                <label class="text-sm text-neutral-400 font-medium">商店名稱</label>
                <input v-model="form.name" type="text" placeholder="輸入商店名稱" class="w-full bg-neutral-900 border border-neutral-800 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-indigo-500 transition-colors" />
            </div>
             <div class="space-y-2">
                <label class="text-sm text-neutral-400 font-medium">Google Maps URL</label>
                <input v-model="form.google_map_url" type="text" placeholder="輸入 Google Maps URL" class="w-full bg-neutral-900 border border-neutral-800 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-indigo-500 transition-colors" />
            </div>
        </div>

        <div class="border-t border-neutral-800 my-4"></div>

        <!-- Photo Management -->
         <div class="space-y-4">
            <label class="text-sm text-neutral-400 font-medium flex items-center gap-2">
                <Icon icon="mdi:image-multiple" /> 管理照片
            </label>
            <div class="bg-neutral-900 border border-neutral-800 rounded-xl p-4">
                <ImageManager 
                    ref="imageManagerRef"
                    :entity-id="route.params.storeId as string" 
                    entity-type="store" 
                    :max-count="5" 
                    :allow-reorder="true"
                    :instant-upload="false"
                    :instant-delete="false"
                />
            </div>
        </div>
        
        <div class="border-t border-neutral-800 my-4"></div>

        <!-- Danger Zone -->
         <div class="space-y-4">
             <button @click="deleteStore" class="w-full py-3 rounded-xl border border-red-500/30 text-red-500 bg-red-500/10 hover:bg-red-500/20 transition-colors flex items-center justify-center gap-2">
                 <Icon icon="mdi:trash-can-outline" /> 刪除商店
             </button>
         </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import ImageManager from '~/components/ImageManager.vue'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const loading = ref(true)
const form = ref({
    name: '',
    google_map_url: ''
})
const imageManagerRef = ref<InstanceType<typeof ImageManager> | null>(null)

const fetchStore = async () => {
    try {
        const store = await api.get<any>(`/api/trips/${route.params.id}/stores/${route.params.storeId}`)
        form.value.name = store.name
        form.value.google_map_url = store.google_map_url || ''
    } catch (e) {
        console.error('Failed to fetch store', e)
    } finally {
        loading.value = false
    }
}

const saveStore = async () => {
    try {
        // 1. Update info
        await api.put(`/api/trips/${route.params.id}/stores/${route.params.storeId}`, {
            name: form.value.name,
            google_map_url: form.value.google_map_url
        })

        // 2. Commit images
        if (imageManagerRef.value) {
            await imageManagerRef.value.commit()
        }

        router.back()
    } catch (e: any) {
        alert(e.message || '儲存失敗')
    }
}

const deleteStore = async () => {
  if (!confirm('確定要刪除此商店？所有的商品也會被刪除。')) return
  try {
    await api.del(`/api/trips/${route.params.id}/stores/${route.params.storeId}`)
    router.push(`/trips/${route.params.id}/stores`)
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
  fetchStore()
})
</script>
