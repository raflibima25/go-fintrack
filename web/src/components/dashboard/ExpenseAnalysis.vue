<script setup>
import { ref, onMounted } from 'vue'
import { Line, Pie, Bar } from 'vue-chartjs'
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    ArcElement,
    BarElement,
    plugins,
    scales
} from 'chart.js'
import { dashboardService } from '../../services/dashboardService';
import { useToast } from '@/composables/useToast';
import { formatCurrency } from '@/utils/formatters'

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    ArcElement,
    BarElement
);

const incomeVsExpenseData = ref(null)
const categoryDistributionData = ref(null)
const topExpensesData = ref(null)

const lineChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            position: 'top',
        },
        title: {
            display: true,
            text: 'Income vs Expense (Last 6 Months)',
        }
    },
    scales: {
        y: {
            beginAtZero: true,
            ticks: {
                callback: (value) => `${formatCurrency(value)}`
            }
        }
    }
};

const pieChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            position: 'right',
        },
        title: {
            display: true,
            text: 'Expense Distribution by Category',
        }
    }
}

const barChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            display: false,
        },
        title: {
            display: true,
            text: 'Top Expenses Categories',
        }
    },
    scales: {
        y: {
            beginAtZero: true,
            ticks: {
                callback: (value) => `${formatCurrency(value)}`
            }
        }
    }
};

// fetch data from API
const fetchChartData = async () => {
    try {
        const response = await dashboardService.getExpenseAnalysis()
        console.log(response.data)
        const data = await response.data

        if (data.status) {
            incomeVsExpenseData.value = {
                labels: data.data.income_vs_expense.labels,
                datasets: data.data.income_vs_expense.datasets.map(dataset => ({
                ...dataset,
                borderColor: dataset.border_color,
                backgroundColor: dataset.background_color
                }))
            };

            categoryDistributionData.value = {
                labels: data.data.category_distribution.labels,
                datasets: data.data.category_distribution.datasets.map(dataset => ({
                data: dataset.data,
                backgroundColor: dataset.background_color
                }))
            };

            topExpensesData.value = {
                labels: data.data.top_expenses.labels,
                datasets: data.data.top_expenses.datasets.map(dataset => ({
                    data: dataset.data,
                    backgroundColor: dataset.background_color
                }))
            };
        }
    } catch (error) {
        console.error("Error fetching chart data: ", error)
        useToast("Error fetching chart data", "error")
    }
};

onMounted(() => {
    fetchChartData()
})
</script>

<template>
    <div class="bg-white rounded-xl shadow-sm p-6 mb-6">
        <h2 class="text-xl font-semibold mb-6">Expense Analysis</h2>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- income vs expense Line Chart -->
            <div class="col-span-2 h-72">
                <Line 
                    v-if="incomeVsExpenseData"
                    :data="incomeVsExpenseData"
                    :options="lineChartOptions"
                />
            </div>

            <!-- category distribution Pie chart -->
            <div class="h-72">
                <Pie 
                    v-if="categoryDistributionData"
                    :data="categoryDistributionData"
                    :options="pieChartOptions"
                />
            </div>

            <!-- top expense bar chart -->
            <div class="h-72">
                <Bar 
                    v-if="topExpensesData"
                    :data="topExpensesData"
                    :options="barChartOptions"
                />
            </div>

        </div>
    </div>
</template>