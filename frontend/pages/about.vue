<template>
  <div class="about-page">
    <PageTitle
      title="關於 LoveLion"
      :breadcrumbs="[{ label: '設定', to: '/settings' }]"
    />

    <div class="flex flex-col gap-8 pb-20">
      <!-- Product Intro -->
      <section class="flex flex-col gap-4">
        <BaseCard padding="p-6" class="shadow-sm">
          <div class="prose prose-invert prose-sm max-w-none">
            <p>
              <strong>LoveLion</strong> 最初是因為和朋友出國時的分帳問題而誕生的。
            </p>
            <p>
              多人旅行中，代付、外幣找零、統一刷卡的匯率與手續費，每一筆都讓結算變得更複雜。以前我都是先用記帳 app 記錄，回國後再手動整理到 Google Sheet，逐筆填入匯率和手續費，最後再一個一個算出每個人該付多少。這個過程既分散又容易出錯。
            </p>
            <p>
              另一個困擾是比價。旅途中看到想買的東西，會沿路比較不同店家的價格，但每次都是拍照記錄，等到要確認的時候又得翻相簿翻很久，很難快速對照。
            </p>
            <p>
              所以我做了 LoveLion — 把記帳、比價、匯率換算、分帳和結算收在同一個地方。旅途中就能即時記錄，每個人隨時都能看到消費明細和自己還需要付多少錢，比價紀錄也不用再靠照片大海撈針。
            </p>
            <p>
              後來發現，日常生活中的記帳需求其實也差不多，就把它一起收進來了。
            </p>
            <p>
              這是一個為了解決自己問題而做的小工具，希望也能幫到有同樣困擾的人。如果你有任何功能建議或想法，歡迎到
              <a href="https://github.com/PlosinBen/LoveLion" target="_blank" rel="noopener noreferrer" class="text-indigo-400 hover:text-indigo-300 transition-colors">GitHub</a>
              上跟我聊聊。
            </p>
          </div>
        </BaseCard>
      </section>

      <!-- App Info -->
      <section class="flex flex-col items-center gap-2 mt-4">
        <img src="/dark-logo.png" alt="LoveLion" class="w-32 h-auto object-contain opacity-50">
        <span class="text-xs text-neutral-600 font-bold">{{ appVersion }}</span>
        <span class="text-xs text-neutral-700 font-bold uppercase tracking-widest">Built by Ben &copy; 2026</span>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import PageTitle from '~/components/PageTitle.vue'
import BaseCard from '~/components/BaseCard.vue'

definePageMeta({
  layout: 'default'
})

const config = useRuntimeConfig()
const appVersion = computed(() => {
  const version = config.public.appVersion as string
  return version ? `v${version}` : 'dev'
})
</script>

<style scoped>
.about-page {
  @apply max-w-lg mx-auto;
}
</style>
