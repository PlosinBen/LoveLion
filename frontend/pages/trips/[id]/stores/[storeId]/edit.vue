<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">編輯商店</h1>
      <div class="w-10"></div>
    </header>

    <div v-if="loading" class="flex justify-center items-center text-neutral-400 py-10">
      <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-5">
        <!-- Basic Info Card -->
        <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-5">
            <BaseInput
                v-model="form.name"
                label="商店名稱"
                placeholder="輸入商店名稱"
            />
            <BaseInput
                v-model="form.google_map_url"
                label="Google Maps URL"
                placeholder="輸入 Google Maps URL"
            />
        </div>

        <!-- Photo Management Card -->
        <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-2">
            <label class="text-sm text-neutral-400 font-medium flex items-center gap-2 mb-1">
                <Icon icon="mdi:image-multiple" /> 管理照片
            </label>
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
        
        <!-- Save Button -->
        <button @click="saveStore" class="w-full px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors border-0 cursor-pointer disabled:opacity-50" :disabled="saving">
            {{ saving ? '儲存中...' : '儲存' }}
        </button>

        <!-- Danger Zone -->
        <div class="mt-4">
            <button @click="deleteStore" class="w-full py-3 rounded-xl border border-red-500/30 text-red-500 bg-red-500/10 hover:bg-red-500/20 transition-colors flex items-center justify-center gap-2 cursor-pointer">
                <Icon icon="mdi:trash-can-outline" /> 刪除商店
            </button>
            <p class="text-xs text-center text-neutral-500 mt-2">刪除商店將會一併刪除所有商品資料</p>
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
import BaseInput from '~/components/BaseInput.vue'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const loading = ref(true)
const saving = ref(false)
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
    if (!form.value.name) {
        alert('請輸入商店名稱')
        return
    }

    saving.value = true
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
    } finally {
        saving.value = false
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
