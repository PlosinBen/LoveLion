<template>
  <div class="add-store-page">
    <PageTitle
      title="新增店家"
      :back-to="`/spaces/${route.params.id}/stores`"
      class="px-2"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}` }, { label: '比價', to: `/spaces/${route.params.id}/stores` }]"
    />

    <div>
      <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">店家資訊</label>
          <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-5 flex flex-col gap-5">
            <BaseInput 
              v-model="form.name" 
              label="店家名稱" 
              placeholder="例如：唐吉訶德、大國藥妝" 
              required
              autofocus
            />
          </div>
        </div>

        <p class="text-xs text-neutral-500 px-4 leading-relaxed">
          建立店家後，您可以開始記錄該店家的商品價格，方便在不同店家間進行比價。
        </p>

        <button 
          type="submit" 
          :disabled="submitting"
          class="w-full py-4 bg-indigo-500 text-white rounded-xl font-bold hover:bg-indigo-600 transition-all active:scale-95 border-0 cursor-pointer shadow-lg disabled:opacity-50"
        >
          {{ submitting ? '建立中...' : '建立店家' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'

definePageMeta({
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const submitting = ref(false)
const form = ref({
  name: ''
})

const handleSubmit = async () => {
  if (!form.value.name.trim()) return
  
  submitting.value = true
  try {
    await api.post(`/api/spaces/${route.params.id}/stores`, form.value)
    detailStore.invalidate('stores')
    router.push(`/spaces/${route.params.id}/stores`)
  } catch (e: any) {
    alert(e.message || '建立失敗')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()
})
</script>
