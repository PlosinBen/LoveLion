<template>
  <div class="edit-product-page">
    <PageTitle
      title="編輯商品"
      :back-to="`/spaces/${spaceId}/stores/${storeId}`"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${spaceId}/ledger` }, { label: '比價', to: `/spaces/${spaceId}/stores` }, { label: storeName || '店家', to: `/spaces/${spaceId}/stores/${storeId}` }]"
    />
    />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <form v-else @submit.prevent="handleSubmit" class="flex flex-col gap-6">
      <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-5 flex flex-col gap-5">
        <BaseInput
          v-model="form.name"
          label="商品名稱"
          placeholder="例如：pocky、午後紅茶"
          required
        />

        <div class="grid grid-cols-2 gap-4">
          <BaseInput
            v-model.number="form.price"
            label="價格"
            type="number"
            step="0.01"
            placeholder="0"
            required
          />
          <BaseInput
            v-model="form.currency"
            label="幣別"
            placeholder="TWD"
          />
        </div>

        <BaseInput
          v-model="form.unit"
          label="單位 (選填)"
          placeholder="例如：包、瓶、100g"
        />

        <BaseTextarea
          v-model="form.note"
          label="備註 (選填)"
          placeholder="口味、規格等補充資訊"
          :rows="2"
        />
      </div>

      <div class="flex flex-col gap-3">
        <BaseButton
          type="submit"
          variant="primary"
          fullWidth
          :loading="submitting"
        >
          儲存變更
        </BaseButton>

        <BaseButton
          type="button"
          @click="handleDelete"
          variant="danger"
          fullWidth
        >
          刪除此商品
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '~/components/BaseButton.vue'
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import BaseInput from '~/components/BaseInput.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'

definePageMeta({
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const spaceId = route.params.id as string
const storeId = route.params.storeId as string
const productId = route.params.productId as string

const storeName = ref('')
const loading = ref(true)
const submitting = ref(false)
const form = ref({
  name: '',
  price: 0,
  currency: 'TWD',
  unit: '',
  note: ''
})

const fetchProduct = async () => {
  try {
    const data = await api.get<any>(`/api/spaces/${spaceId}/stores/${storeId}/products/${productId}`)
    form.value = {
      name: data.name,
      price: parseFloat(data.price),
      currency: data.currency || 'TWD',
      unit: data.unit || '',
      note: data.note || ''
    }
  } catch (e) {
    console.error('Failed to fetch product:', e)
    router.push(`/spaces/${spaceId}/stores/${storeId}`)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!form.value.name.trim() || !form.value.price) return

  submitting.value = true
  try {
    await api.put(`/api/spaces/${spaceId}/stores/${storeId}/products/${productId}`, form.value)
    detailStore.invalidate('stores')
    router.push(`/spaces/${spaceId}/stores/${storeId}`)
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除此商品嗎？')) return

  try {
    await api.del(`/api/spaces/${spaceId}/stores/${storeId}/products/${productId}`)
    detailStore.invalidate('stores')
    router.push(`/spaces/${spaceId}/stores/${storeId}`)
  } catch (e: any) {
    alert(e.message || '刪除失敗')
  }
}

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  detailStore.setSpaceId(spaceId)
  detailStore.fetchSpace()
  try {
    const store = await api.get<any>(`/api/spaces/${spaceId}/stores/${storeId}`)
    storeName.value = store.name || ''
  } catch {}
  fetchProduct()
})
</script>
