<template>
  <div class="add-store-page min-h-screen bg-neutral-900 text-neutral-50">
    <SpaceHeader 
      title="新增店家" 
      :show-back="true"
      class="pt-0 px-2"
    />

    <div class="p-4 pt-0">
      <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
        <div class="flex flex-col gap-2">
          <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">店家資訊</label>
          <div class="bg-neutral-900 rounded-3xl border border-neutral-800 p-5 flex flex-col gap-5">
            <BaseInput 
              v-model="form.name" 
              label="店家名稱" 
              placeholder="例如：全聯、Costco、7-11" 
              required
              autofocus
            />
            <BaseTextarea 
              v-model="form.description" 
              label="備註 (選填)" 
              placeholder="店家位置、特點或常用優惠資訊" 
              rows="3"
            />
          </div>
        </div>

        <button 
          type="submit" 
          class="w-full py-4 bg-indigo-500 text-white rounded-2xl font-bold hover:bg-indigo-600 transition-all active:scale-[0.98] border-0 cursor-pointer shadow-lg shadow-indigo-500/20 disabled:opacity-50"
          :disabled="submitting"
        >
          {{ submitting ? '儲存中...' : '儲存店家' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  hideGlobalNav: true
})
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const submitting = ref(false)
const form = ref({
  name: '',
  description: ''
})

const handleSubmit = async () => {
  if (!form.value.name) return
  submitting.value = true
  
  try {
    await api.post(`/api/spaces/${route.params.id}/stores`, form.value)
    // Return to space detail page, showing the comparison tab
    router.push({
      path: `/spaces/${route.params.id}`,
      query: { tab: 'comparison' }
    })
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
