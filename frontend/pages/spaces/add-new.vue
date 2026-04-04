<template>
  <div class="add-space-page">
    <PageTitle
      title="建立新空間"
      :breadcrumbs="[{ label: '空間列表', to: '/' }]"
    />

    <div>
      <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">基本資訊</label>
          <BaseCard padding="p-6" class="flex flex-col gap-6 shadow-sm">
            <BaseInput 
              v-model="form.name" 
              label="空間名稱" 
              placeholder="例如：日本旅行、個人記帳" 
              required
              autofocus
            />

            <BaseSelect
              v-model="form.base_currency"
              label="本位幣別"
              :options="currencyOptions"
            />
          </BaseCard>
        </div>

        <p class="text-xs text-neutral-500 px-4 leading-relaxed font-medium">
          建立空間後，您可以邀請成員共同記帳、上傳收據照片，或是進行商品比價與費用分攤。
        </p>

        <BaseButton
          type="submit"
          class="w-full mt-4"
        >
          建立空間
        </BaseButton>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '~/composables/useApi'
import { useLoading } from '~/composables/useLoading'
import { useToast } from '~/composables/useToast'
import PageTitle from '~/components/PageTitle.vue'
import BaseCard from '~/components/BaseCard.vue'

definePageMeta({
  layout: 'default',
  hideGlobalNav: true
})

const router = useRouter()
const api = useApi()
const { showLoading, hideLoading } = useLoading()
const toast = useToast()

const form = ref({
  name: '',
  base_currency: 'TWD',
  type: 'personal',
  categories: ['餐飲', '交通', '購物', '娛樂', '生活', '其他'],
  currencies: ['TWD']
})

const currencyOptions = [
  { label: 'TWD - 新台幣', value: 'TWD' },
  { label: 'JPY - 日圓', value: 'JPY' },
  { label: 'USD - 美金', value: 'USD' },
  { label: 'EUR - 歐元', value: 'EUR' }
]

const handleSubmit = async () => {
  if (!form.value.name.trim()) return
  
  showLoading()
  try {
    if (!form.value.currencies.includes(form.value.base_currency)) {
        form.value.currencies.push(form.value.base_currency)
    }

    await api.post('/api/spaces', form.value)
    router.push('/')
  } catch (e: any) {
    toast.error(e.message || '建立失敗')
  } finally {
    hideLoading()
  }
}

</script>
