<template>
  <div class="p-6">
    <!-- Header Section -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Transactions</h1>
      <div class="flex gap-3">
        <button @click="openCreateModal"
                class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600">
          Add Transaction
        </button>
        <button @click="exportData"
                class="bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600">
          Export Excel
        </button>
      </div>
    </div>

    <!-- Filter Section -->
    <div class="bg-white p-4 rounded-lg shadow mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
          <input type="date" v-model="filter.startDate"
                 class="w-full border rounded-lg px-3 py-2" @change="fetchTransactions">
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">End Date</label>
          <input type="date" v-model="filter.endDate"
                 class="w-full border rounded-lg px-3 py-2" @change="fetchTransactions">
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select v-model="filter.categoryId"
                  class="w-full border rounded-lg px-3 py-2" @change="fetchTransactions">
            <option value="">All Categories</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select v-model="filter.type"
                  class="w-full border rounded-lg px-3 py-2" @change="fetchTransactions">
            <option value="">All Types</option>
            <option value="income">Income</option>
            <option value="expense">Expense</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
      <div class="bg-white p-4 rounded-lg shadow">
        <h3 class="text-lg font-semibold text-gray-600">Total Income</h3>
        <p class="text-2xl font-bold text-green-500">
          {{ formatCurrency(summary.totalIncome) }}
        </p>
      </div>
      <div class="bg-white p-4 rounded-lg shadow">
        <h3 class="text-lg font-semibold text-gray-600">Total Expense</h3>
        <p class="text-2xl font-bold text-red-500">
          {{ formatCurrency(summary.totalExpense) }}
        </p>
      </div>
      <div class="bg-white p-4 rounded-lg shadow">
        <h3 class="text-lg font-semibold text-gray-600">Balance</h3>
        <p class="text-2xl font-bold" :class="summary.balance >= 0 ? 'text-green-500' : 'text-red-500'">
          {{ formatCurrency(summary.balance) }}
        </p>
      </div>
    </div>

    <!-- Transactions Table -->
    <div class="bg-white rounded-lg shadow overflow-x-auto">
      <table class="min-w-full">
        <thead class="bg-gray-50">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
        </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
        <tr v-for="transaction in transactions" :key="transaction.id">
          <td class="px-6 py-4 whitespace-nowrap">{{ formatDate(transaction.date) }}</td>
          <td class="px-6 py-4 whitespace-nowrap">
              <span :class="transaction.type === 'income' ? 'text-green-500' : 'text-red-500'">
                {{ transaction.type }}
              </span>
          </td>
          <td class="px-6 py-4 whitespace-nowrap">{{ transaction.category }}</td>
          <td class="px-6 py-4">{{ transaction.description }}</td>
          <td class="px-6 py-4 whitespace-nowrap"
              :class="transaction.type === 'income' ? 'text-green-500' : 'text-red-500'">
            {{ formatCurrency(transaction.amount) }}
          </td>
          <td class="px-6 py-4 whitespace-nowrap">
            <button @click="editTransaction(transaction)"
                    class="text-blue-500 hover:text-blue-700 mr-3">
              Edit
            </button>
            <button @click="confirmDelete(transaction)"
                    class="text-red-500 hover:text-red-700">
              Delete
            </button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="mt-4 flex justify-between items-center">
      <div class="text-sm text-gray-700">
        Showing {{ ((pagination.currentPage - 1) * pagination.itemPerPage) + 1 }}
        to {{ Math.min(pagination.currentPage * pagination.itemPerPage, pagination.totalItems) }}
        of {{ pagination.totalItems }} entries
      </div>
      <div class="flex gap-2">
        <button
            @click="changePage(pagination.currentPage - 1)"
            :disabled="pagination.currentPage === 1"
            class="px-3 py-1 rounded border"
            :class="pagination.currentPage === 1 ? 'bg-gray-100 text-gray-400' : 'hover:bg-gray-100'">
          Previous
        </button>
        <button
            @click="changePage(pagination.currentPage + 1)"
            :disabled="pagination.currentPage === pagination.totalPage"
            class="px-3 py-1 rounded border"
            :class="pagination.currentPage === pagination.totalPage ? 'bg-gray-100 text-gray-400' : 'hover:bg-gray-100'">
          Next
        </button>
      </div>
    </div>

    <!-- Transaction Modal -->
    <TransactionModal
        v-if="showModal"
        :show="showModal"
        :categories="categories"
        :transaction="selectedTransaction"
        @close="closeModal"
        @submit="handleSubmit" />

    <!-- Delete Confirmation Modal -->
    <ConfirmationModal
        v-if="showDeleteModal"
        :show="showDeleteModal"
        @close="showDeleteModal = false"
        @confirm="handleDelete">
      <template #title>Delete Transaction</template>
      <template #content>
        Are you sure you want to delete this transaction? This action cannot be undone.
      </template>
    </ConfirmationModal>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { transactionService } from '@/services/transactionService'
import { categoryService } from '@/services/categoryService'
import TransactionModal from '../components/TransactionModal.vue'
import ConfirmationModal from '@/components/ConfirmationModal.vue'
import { formatCurrency, formatDate } from '@/utils/formatters'
import { useToast } from '@/composables/useToast'

export default {
  name: 'TransactionList',
  components: {
    TransactionModal,
    ConfirmationModal
  },

  setup() {
    const { showToast } = useToast()
    const transactions = ref([])
    const categories = ref([])
    const showModal = ref(false)
    const showDeleteModal = ref(false)
    const selectedTransaction = ref(null)
    const summary = ref({
      totalIncome: 0,
      totalExpense: 0,
      balance: 0
    })
    const pagination = ref({
      currentPage: 1,
      totalPage: 1,
      totalItems: 0,
      itemPerPage: 10
    })
    const filter = ref({
      startDate: '',
      endDate: '',
      categoryId: '',
      type: '',
      page: 1,
      limit: 10
    })

    const fetchTransactions = async () => {
      try {
        const response = await transactionService.getTransactions(filter.value)
        transactions.value = response.data.data.transactions
        summary.value = response.data.data.summary
        pagination.value = response.data.data.pagination
      } catch (error) {
        showToast('Error fetching transactions', 'error')
      }
    }

    const fetchCategories = async () => {
      try {
        const response = await categoryService.getCategories()
        categories.value = response.data.data
      } catch (error) {
        showToast('Error fetching categories', 'error')
      }
    }

    const openCreateModal = () => {
      selectedTransaction.value = null
      showModal.value = true
    }

    const editTransaction = (transaction) => {
      selectedTransaction.value = transaction
      showModal.value = true
    }

    const closeModal = () => {
      showModal.value = false
      selectedTransaction.value = null
    }

    const handleSubmit = async (formData) => {
      try {
        if (selectedTransaction.value) {
          await transactionService.updateTransaction(selectedTransaction.value.id, formData)
          showToast('Transaction updated successfully')
        } else {
          await transactionService.createTransaction(formData)
          showToast('Transaction created successfully')
        }
        await fetchTransactions()
        closeModal()
      } catch (error) {
        showToast(error.response?.data?.message || 'Error processing transaction', 'error')
      }
    }

    const confirmDelete = (transaction) => {
      selectedTransaction.value = transaction
      showDeleteModal.value = true
    }

    const handleDelete = async () => {
      try {
        await transactionService.deleteTransaction(selectedTransaction.value.id)
        showToast('Transaction deleted successfully')
        await fetchTransactions()
        showDeleteModal.value = false
      } catch (error) {
        showToast('Error deleting transaction', 'error')
      }
    }

    const changePage = (page) => {
      if (page >= 1 && page <= pagination.value.totalPage) {
        filter.value.page = page
        fetchTransactions()
      }
    }

    const exportData = async () => {
      try {
        const response = await transactionService.exportTransactions(filter.value)
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `transactions_${new Date().toISOString().split('T')[0]}.xlsx`)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        showToast('Export successful')
      } catch (error) {
        showToast('Error exporting transactions', 'error')
      }
    }

    onMounted(() => {
      fetchTransactions()
      fetchCategories()
    })

    return {
      transactions,
      categories,
      showModal,
      showDeleteModal,
      selectedTransaction,
      summary,
      pagination,
      filter,
      formatCurrency,
      formatDate,
      openCreateModal,
      editTransaction,
      closeModal,
      handleSubmit,
      confirmDelete,
      handleDelete,
      changePage,
      exportData,
      fetchTransactions
    }
  }
}
</script>