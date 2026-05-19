<template>
  <OverlayPage>
    <PageTitle
      title="新增出入金"
      :show-back="true"
      :breadcrumbs="[{ label: '出入金', to: '/investments/transactions' }]"
    />
    <MemberTransactionForm @saved="router.back()" />
  </OverlayPage>
</template>

<script setup lang="ts">
import PageTitle from '~/components/PageTitle.vue'
import OverlayPage from '~/components/OverlayPage.vue'
import MemberTransactionForm from '~/components/investment/MemberTransactionForm.vue'

const router = useRouter()
</script>
