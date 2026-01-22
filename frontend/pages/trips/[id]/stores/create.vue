<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">新增商店</h1>
      <div class="w-10"></div>
    </header>

    <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
      <BaseInput v-model="name" placeholder="商店名稱" input-class="mb-4" @keyup.enter="handleSubmit" :auto-focus="true" />
      
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

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const name = ref('')
const submitting = ref(false)

const handleSubmit = async () => {
  if (!name.value.trim()) return
  
  submitting.value = true
  try {
    await api.post(`/api/trips/${route.params.id}/stores`, { name: name.value.trim() })
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
