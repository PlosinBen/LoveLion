<template>
  <OverlayPage>
    <PageTitle
      title="編輯出入金"
      :show-back="true"
      :breadcrumbs="[{ label: '出入金', to: '/investments/transactions' }]"
    />
    <MemberTransactionForm :transaction-id="route.params.id as string" @saved="router.back()" @deleted="router.back()" />
  </OverlayPage>
</template>

<script setup lang="ts">
import PageTitle from '~/components/PageTitle.vue'
import OverlayPage from '~/components/OverlayPage.vue'
import MemberTransactionForm from '~/components/investment/MemberTransactionForm.vue'

const route = useRoute()
const router = useRouter()
</script>
