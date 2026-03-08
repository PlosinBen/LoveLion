<template>
  <div class="store-edit-page">
    <SpaceHeader title="編輯商店" :back-to="`/spaces/${spaceId}/stores/${storeId}`" />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <form v-else @submit.prevent="handleUpdate" class="flex flex-col gap-6">
      <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-6 flex flex-col gap-6">
        <BaseInput 
          v-model="form.name" 
          label="商店名稱" 
          placeholder="例如：全家便利商店、宜得利" 
          required
        />

        <BaseInput 
          v-model="form.google_map_url" 
          label="Google Maps 連結 (選填)" 
          placeholder="貼上地圖連結以方便導航" 
        />
        
        <BaseTextarea 
          v-model="form.description" 
          label="簡介或備註" 
          placeholder="輸入關於此商店的筆記..." 
          :rows="3"
        />

        <div class="mt-2 flex gap-3">
          <button 
            type="button"
            @click="router.push(`/spaces/${spaceId}/stores/${storeId}`)"
            class="flex-1 py-4 rounded-xl bg-neutral-800 text-white font-bold border-0 cursor-pointer"
          >
            取消
          </button>
          <button 
            type="submit"
            :disabled="updating"
            class="flex-[2] py-4 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer shadow-lg"
          >
            {{ updating ? '儲存中...' : '儲存變更' }}
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'
import BaseInput from '~/components/BaseInput.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'

definePageMeta({
  layout: 'default'
})

const route = useRoute()
const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const loading = ref(true)
const updating = ref(false)
const spaceId = route.params.id as string
const storeId = route.params.storeId as string

const form = ref({
  name: '',
  description: '',
  google_map_url: ''
})

const fetchData = async () => {
  try {
    const data = await api.get<any>(`/api/spaces/${spaceId}/stores/${storeId}`)
    form.value = {
      name: data.name,
      description: data.description || '',
      google_map_url: data.google_map_url || ''
    }
  } catch (e) {
    console.error('Failed to fetch store:', e)
    router.push(`/spaces/${spaceId}/stores`)
  } finally {
    loading.value = false
  }
}

const handleUpdate = async () => {
  updating.value = true
  try {
    await api.patch(`/api/spaces/${spaceId}/stores/${storeId}`, form.value)
    router.push(`/spaces/${spaceId}/stores/${storeId}`)
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    updating.value = false
  }
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
