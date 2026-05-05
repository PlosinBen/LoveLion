<template>
  <OverlayPage>
    <PageTitle
      title="新增交易"
      :show-back="true"
      :breadcrumbs="[{ label: '交易', to: '/investments/trades' }]"
    />
    <StockTradeForm @saved="router.back()" />
  </OverlayPage>
</template>

<script setup lang="ts">
import PageTitle from '~/components/PageTitle.vue'
import OverlayPage from '~/components/OverlayPage.vue'
import StockTradeForm from '~/components/investment/StockTradeForm.vue'

const router = useRouter()
</script>
