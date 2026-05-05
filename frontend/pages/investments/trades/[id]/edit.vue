<template>
  <OverlayPage>
    <PageTitle
      title="編輯交易"
      :show-back="true"
      :breadcrumbs="[{ label: '交易', to: '/investments/trades' }]"
    />
    <StockTradeForm :trade-id="route.params.id as string" @saved="router.back()" @deleted="router.back()" />
  </OverlayPage>
</template>

<script setup lang="ts">
import PageTitle from '~/components/PageTitle.vue'
import OverlayPage from '~/components/OverlayPage.vue'
import StockTradeForm from '~/components/investment/StockTradeForm.vue'

const route = useRoute()
const router = useRouter()
</script>
