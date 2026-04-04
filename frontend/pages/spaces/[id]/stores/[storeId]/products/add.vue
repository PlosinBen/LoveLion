<template>
  <div class="add-product-page">
    <PageTitle
      title="新增商品"
      :back-to="`/spaces/${spaceId}/stores/${storeId}`"
      class="px-2"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${spaceId}/ledger` }, { label: '比價', to: `/spaces/${spaceId}/stores` }, { label: storeName || '店家', to: `/spaces/${spaceId}/stores/${storeId}` }]"
    />

    <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
      <BaseCard padding="p-5" class="flex flex-col gap-5">
        <BaseInput
          v-model="form.name"
          label="商品名稱"
          placeholder="例如：pocky、午後紅茶"
          required
          autofocus
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
      </BaseCard>

      <BaseButton
        type="submit"
        variant="primary"
        class="w-full"
      >
        新增商品
      </BaseButton>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useLoading } from '~/composables/useLoading'
import { useToast } from '~/composables/useToast'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import BaseInput from '~/components/BaseInput.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'
import BaseCard from '~/components/BaseCard.vue'

definePageMeta({
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { showLoading, hideLoading } = useLoading()
const toast = useToast()
const detailStore = useSpaceDetailStore()

const spaceId = route.params.id as string
const storeId = route.params.storeId as string

const storeName = ref('')
const form = ref({
  name: '',
  price: 0,
  currency: 'TWD',
  unit: '',
  note: ''
})

const handleSubmit = async () => {
  if (!form.value.name.trim() || !form.value.price) return

  showLoading()
  try {
    await api.post(`/api/spaces/${spaceId}/stores/${storeId}/products`, form.value)
    detailStore.invalidate('stores')
    router.push(`/spaces/${spaceId}/stores/${storeId}`)
  } catch (e: any) {
    toast.error(e.message || '新增失敗')
  } finally {
    hideLoading()
  }
}

onMounted(async () => {
  detailStore.setSpaceId(spaceId)
  detailStore.fetchSpace()
  try {
    const store = await api.get<any>(`/api/spaces/${spaceId}/stores/${storeId}`)
    storeName.value = store.name || ''
  } catch {}
})
</script>
