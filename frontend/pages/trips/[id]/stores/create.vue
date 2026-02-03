<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">新增商店</h1>
      <div class="w-10"></div>
    </header>

    <div class="flex flex-col gap-5">
      <!-- Info Card -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-5">
        <BaseInput 
            v-model="name" 
            label="商店名稱"
            placeholder="商店名稱" 
            @keyup.enter="handleSubmit" 
            :auto-focus="true" 
        />
        <BaseInput 
            v-model="googleMapUrl" 
            label="Google Maps URL"
            placeholder="Google Maps URL (選填)" 
        />
      </div>
      
      <!-- Photo Card -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-2">
          <label class="text-sm text-neutral-400 font-medium flex items-center gap-2 mb-1">
              <Icon icon="mdi:image-multiple" /> 商店照片
          </label>
          <ImageManager 
            ref="imageManagerRef"
            entity-id="" 
            entity-type="store" 
            :max-count="5" 
            :instant-upload="false"
          />
      </div>

      <button 
        @click="handleSubmit" 
        class="w-full px-4 py-3 rounded-xl bg-indigo-500 text-white border-0 cursor-pointer hover:bg-indigo-600 transition-colors font-bold text-base disabled:opacity-50 disabled:cursor-not-allowed" 
        :disabled="!name.trim() || submitting"
      >
        {{ submitting ? '新增中...' : '新增商店' }}
      </button>
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

const name = ref('')
const googleMapUrl = ref('')
const submitting = ref(false)
const imageManagerRef = ref<InstanceType<typeof ImageManager> | null>(null)

const handleSubmit = async () => {
  if (!name.value.trim()) return
  
  submitting.value = true
  try {
    const store = await api.post<any>(`/api/trips/${route.params.id}/stores`, { 
        name: name.value.trim(),
        google_map_url: googleMapUrl.value.trim() 
    })
    
    // Upload images if any
    if (imageManagerRef.value) {
        await imageManagerRef.value.commit(store.id)
    }

    router.push(`/trips/${route.params.id}/stores`)
  } catch (e: any) {
    alert(e.message || '新增失敗')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
  }
})
</script>
