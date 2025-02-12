<script setup>
import { ref, onMounted, computed } from 'vue'
import { dashboardService } from '../../services/dashboardService';
import { formatCurrency } from '../../utils/formatters';
import { Wallet, TrendingUp, TrendingDown, PiggyBank } from 'lucide-vue-next';

const overview = ref({})
const loading = ref(true)

const fetchOverview = async () => {
    try {
        loading.value = true
        const { data } = await dashboardService.getFinancialOverview()
        overview.value = data.data
    } catch (error) {
        console.error(error)
    } finally {
        loading.value = false
    }
}

const overviewCards = computed(() => [
    {
        title: 'Total Balance',
        amount: overview.value.current_balance,
        trend: 5.2,
        bgColor: 'bg-blue-100',
        iconBgColor: 'bg-blue-100',
        iconColor: 'text-blue-600',
        trendColor: 'text-green-600',
        icon: Wallet
    },
    {
        title: 'Income of the Month',
        amount: overview.value.monthly_income,
        trend: 3.1,
        bgColor: 'bg-green-100',
        iconBgColor: 'bg-green-100',
        iconColor: 'text-green-600',
        trendColor: 'text-green-600',
        icon: TrendingUp
    },
    {
        title: 'Expenses of the Month',
        amount: overview.value.monthly_expense,
        trend: -2.4,
        bgColor: 'bg-red-100',
        iconBgColor: 'bg-red-100',
        iconColor: 'text-red-600',
        trendColor: 'text-red-600',
        icon: TrendingDown
    },
    {
        title: 'Total Savings',
        amount: overview.value.total_savings,
        trend: 8.7,
        bgColor: 'bg-purple-50',
        iconBgColor: 'bg-purple-100',
        iconColor: 'text-purple-600',
        trendColor: 'text-green-600',
        icon: PiggyBank
    },
])

onMounted(() => {
    fetchOverview()
})
</script>

<template>
    <div class="bg-white rounded-xl shadow-sm p-6 mb-6">
        <h2 class="text-xl font-semibold mb-6">Financial Overview</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div v-for="(card, index) in overviewCards"
                :key="index"
                class="rounded-lg p-6 transition-all duration-200 hover:shadow-md"
                :class="card.bgColor">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-gray-700 font-medium">{{ card.title }}</h3>
                    <div class="p-2 rounded-full" :class="card.iconBgColor">
                        <component :is="card.icon"
                            :class="card.iconColor"
                            class="w-6 h-6"
                        />
                    </div>
                </div>
                <p class="text-2xl font-bold text-gray-900 mb-2">
                    {{ formatCurrency(card.amount) }}
                </p>
                <div class="flex items-center gap-2">
                    <span :class="card.trendColor" class="flex items-center text-sm font-medium">
                        <component :is="card.trend > 0 ? TrendingUp : TrendingDown"
                            :class="card.trend > 0 ? 'text-green-500' : 'text-red-500'"
                            class="w-4 h-4 mr-1" />
                        {{ Math.abs(card.trend) }}%
                    </span>
                </div>
                <p class="text-gray-500 text-sm ml-2">vs last month</p>
            </div>
        </div>
    </div>
</template>