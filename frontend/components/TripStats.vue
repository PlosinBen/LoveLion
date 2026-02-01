<template>
    <div class="flex flex-col gap-6 px-4">
        
        <!-- Total Spending Card -->
        <div class="bg-neutral-800 rounded-2xl p-6 flex flex-col items-center justify-center relative overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/10 to-purple-500/10 z-0"></div>
            <div class="relative z-10 text-center">
                <span class="text-neutral-400 text-sm mb-1 block">總支出 ({{ baseCurrency }})</span>
                <div class="text-4xl font-bold text-white tracking-tight">
                    <span class="text-2xl text-neutral-500 mr-1">$</span>
                    {{ formatNumber(totalSpent) }}
                </div>
            </div>
        </div>

        <!-- Settlement (New) -->
        <div v-if="settlements.length > 0">
             <h2 class="text-lg font-bold mb-4 flex items-center gap-2">
                <Icon icon="mdi:handshake-outline" class="text-indigo-400" />
                應收應付 (Settlement)
            </h2>
            <div class="bg-neutral-900 rounded-xl border border-neutral-800 divide-y divide-neutral-800">
                <div v-for="member in settlements" :key="member.name" class="p-4 flex items-center justify-between">
                     
                     <!-- Left: Name -->
                     <div class="flex items-center gap-3">
                         <div class="w-8 h-8 rounded-full bg-neutral-800 flex items-center justify-center text-sm font-bold text-neutral-400 border border-neutral-700">
                             {{ member.name.charAt(0) }}
                         </div>
                         <span class="font-bold text-base">{{ member.name }}</span>
                     </div>

                     <!-- Right: Amounts -->
                     <div class="flex flex-col items-end gap-0.5">
                        
                        <!-- Base Currency -->
                        <div class="font-mono font-bold text-base" :class="member.baseAmount >= 0 ? 'text-green-400' : 'text-red-400'">
                            {{ baseCurrency }} {{ formatNumber(member.baseAmount) }}
                        </div>

                        <!-- Foreign Currencies (Unconverted) -->
                        <div v-for="(amount, cur) in member.foreign" :key="cur" class="text-xs text-neutral-400 font-mono">
                            {{ cur }} {{ formatNumber(amount) }}
                        </div>

                     </div>
                </div>
            </div>
        </div>

        <!-- Category Breakdown -->
        <div v-if="categories.length > 0">
            <h2 class="text-lg font-bold mb-4 flex items-center gap-2">
                <Icon icon="mdi:shape-outline" class="text-indigo-400" />
                分類支出
            </h2>
            <div class="flex flex-col gap-3">
                <div v-for="cat in categories" :key="cat.name" class="bg-neutral-900 rounded-xl p-3 border border-neutral-800">
                     <div class="flex justify-between items-center mb-2">
                         <div class="flex items-center gap-2">
                             <div class="w-8 h-8 rounded-full bg-neutral-800 flex items-center justify-center text-lg">
                                 <Icon :icon="getCategoryIcon(cat.name)" class="text-neutral-400" />
                             </div>
                             <span class="font-medium text-sm">{{ cat.name }}</span>
                         </div>
                         <div class="text-right">
                             <div class="font-bold text-sm">{{ formatNumber(cat.amount) }}</div>
                             <div class="text-xs text-neutral-500">{{ cat.percentage }}%</div>
                         </div>
                     </div>
                     <!-- Progress Bar -->
                     <div class="h-1.5 bg-neutral-800 rounded-full overflow-hidden">
                         <div class="h-full bg-indigo-500 rounded-full transition-all duration-1000" :style="{ width: cat.percentage + '%' }"></div>
                     </div>
                </div>
            </div>
        </div>

        <!-- Top Spenders -->
        <div v-if="spenders.length > 0">
             <h2 class="text-lg font-bold mb-4 flex items-center gap-2">
                <Icon icon="mdi:account-cash-outline" class="text-indigo-400" />
                花費排行
            </h2>
            <div class="bg-neutral-900 rounded-xl border border-neutral-800 divide-y divide-neutral-800">
                <div v-for="(spender, index) in spenders" :key="spender.name" class="p-4 flex items-center justify-between">
                     <div class="flex items-center gap-3">
                         <div class="w-6 h-6 rounded-full bg-neutral-800 flex items-center justify-center text-xs font-bold text-neutral-500">
                             {{ index + 1 }}
                         </div>
                         <span class="font-medium">{{ spender.name }}</span>
                     </div>
                     <span class="font-bold text-indigo-400">{{ formatNumber(spender.amount) }}</span>
                </div>
            </div>
        </div>

    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
    transactions: any[]
    members: any[]
    baseCurrency: string
}>()

const totalSpent = ref(0)
const categories = ref<any[]>([])
const spenders = ref<any[]>([])
const settlements = ref<any[]>([])

const formatNumber = (num: number) => {
    return num.toLocaleString(undefined, { maximumFractionDigits: 1 })
}

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    'Food': 'mdi:food',
    'Transport': 'mdi:train-car',
    'Shopping': 'mdi:shopping',
    'Entertainment': 'mdi:movie',
    'Accommodation': 'mdi:bed',
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
    '住宿': 'mdi:bed',
    '生活': 'mdi:home',
    '其他': 'mdi:dots-horizontal',
  }
  return icons[category] || 'mdi:shape'
}

const calculateStats = () => {
    if (!props.transactions) return

    let total = 0
    const catMap: Record<string, number> = {}
    const spenderMap: Record<string, number> = {}
    
    // Settlement Data Structures
    const balances: Record<string, { 
        name: string, 
        baseAmount: number, 
        foreign: Record<string, number> 
    }> = {}

    // Init Members
    props.members?.forEach(m => {
        balances[m.id || m.name] = {
            name: m.name,
            baseAmount: 0,
            foreign: {}
        }
        spenderMap[m.name] = 0
    })

    props.transactions.forEach(txn => {
        const baseCurrency = props.baseCurrency
        // 1. Total & Categories (Using Converted Base Amount)
        let baseAmountForStats = 0
        const isSameCurrency = txn.currency === baseCurrency
        
        if (isSameCurrency) {
            baseAmountForStats = parseFloat(txn.total_amount)
        } else {
             // If Billing Amount exists, use it.
             const billing = parseFloat(txn.billing_amount)
             if (billing > 0) {
                 baseAmountForStats = billing
             } else {
                 baseAmountForStats = 0 
             }
        }
        
        total += baseAmountForStats

        const cat = txn.category || '未分類'
        catMap[cat] = (catMap[cat] || 0) + baseAmountForStats

        // 2. Settlement (Payables/Receivables)
        const txnTotal = parseFloat(txn.total_amount)
        if (txnTotal === 0) return 

        let rate = 0
        if (isSameCurrency) {
            rate = 1
        } else if (baseAmountForStats > 0) {
            rate = baseAmountForStats / txnTotal
        }
        
        const isUnconverted = (rate === 0 && !isSameCurrency)

        const processFlow = (memberId: string | undefined, memberName: string, amount: number, isPayer: boolean) => {
             let balanceEntry = null
             if (memberId && balances[memberId]) {
                 balanceEntry = balances[memberId]
             } else {
                 // Find by name
                 const key = Object.keys(balances).find(k => balances[k].name === memberName)
                 if (key) balanceEntry = balances[key]
             }
             
             if (!balanceEntry) return 

             const sign = isPayer ? 1 : -1
             
             if (!isPayer) {
                 if (!isUnconverted) {
                     spenderMap[memberName] = (spenderMap[memberName] || 0) + (amount * rate)
                 }
             }

             // Update Settlement
             if (isUnconverted) {
                 const cur = txn.currency
                 const val = amount * sign
                 balanceEntry.foreign[cur] = (balanceEntry.foreign[cur] || 0) + val
             } else {
                 const val = amount * rate * sign
                 balanceEntry.baseAmount += val
             }
        }

        if (txn.splits && txn.splits.length > 0) {
            txn.splits.forEach((split: any) => {
                processFlow(split.member_id, split.name, parseFloat(split.amount), split.is_payer)
            })
        } else {
            // No splits? Assume Payer paid, Payer consumed
            processFlow(txn.payer_id, txn.payer, txnTotal, true)  // Paid
            processFlow(txn.payer_id, txn.payer, txnTotal, false) // Consumed
        }
    })

    totalSpent.value = total
    
    // Categories
    categories.value = Object.entries(catMap)
        .map(([name, amount]) => ({
            name,
            amount,
            percentage: total > 0 ? Math.round((amount / total) * 100) : 0
        }))
        .sort((a, b) => b.amount - a.amount)

    // Spenders
    spenders.value = Object.entries(spenderMap)
        .map(([name, amount]) => ({ name, amount }))
        .sort((a, b) => b.amount - a.amount)

    // Settlements - cleanup zero values
    settlements.value = Object.values(balances)
        .map(b => {
             // Filter zero foreign
             const foreignClean: Record<string, number> = {}
             for (const [c, amt] of Object.entries(b.foreign)) {
                 if (Math.abs(amt) > 0.1) foreignClean[c] = amt
             }
             return {
                 name: b.name,
                 baseAmount: b.baseAmount,
                 foreign: foreignClean
             }
        })
        .filter(b => Math.abs(b.baseAmount) > 1 || Object.keys(b.foreign).length > 0)
        .sort((a, b) => b.baseAmount - a.baseAmount) 
}

watch(() => props.transactions, () => {
    calculateStats()
}, { immediate: true, deep: true })
</script>
