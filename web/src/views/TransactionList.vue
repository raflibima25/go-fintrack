// views/TransactionList.vue
<script setup>
import { ref, onMounted, computed } from 'vue';
import { 
  PlusIcon, 
  FileSpreadsheetIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  SearchIcon,
  FilterIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  MoreHorizontalIcon,
  ChevronDownIcon
} from 'lucide-vue-next';
import { 
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useToast } from '@/composables/useToast';
import { transactionService } from '@/services/transactionService';
import { categoryService } from '@/services/categoryService';
import AddTransactionModal from '@/components/transaction/AddTransactionModal.vue';
import UserLayout from '@/layouts/UserLayout.vue';
import { formatDate, formatCurrency, parseDate } from '@/utils/formatters';

const { showToast } = useToast();
const transactions = ref([]);
const categories = ref([]);
const isLoading = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const selectedTransaction = ref(null);
const searchQuery = ref('');
const isFilterExpanded = ref(false);

const summary = ref({
  total_income: 0,
  total_expense: 0,
  balance: 0
});

const pagination = ref({
  current_page: 1,
  total_page: 1,
  total_items: 0,
  item_per_page: 10
});

const filter = ref({
  start_date: '',
  end_date: '',
  category_id: '',
  type: '',
  page: 1,
  limit: 10
});

const fetchTransactions = async () => {
  try {
    isLoading.value = true;
    const response = await transactionService.getTransactions(filter.value);
    if (response.data.status) {
      transactions.value = response.data.data.transactions;
      summary.value = response.data.data.summary;
      pagination.value = response.data.data.pagination;
    } else {
      showToast(response.data.message || 'Error fetching transactions', 'error');
    }
  } catch (error) {
    showToast(error.response?.data?.message || 'Error fetching transactions', 'error');
  } finally {
    isLoading.value = false;
  }
};

const fetchCategories = async () => {
  try {
    const response = await categoryService.getCategories();
    if (response.data.status) {
      categories.value = response.data.data.categories;
    }
  } catch (error) {
    showToast(error.response?.data?.message || 'Error fetching categories', 'error');
  }
};

const handleTransactionAdded = () => {
  showAddDialog.value = false;
  fetchTransactions();
};

const handleTransactionEdited = () => {
  showEditDialog.value = false;
  selectedTransaction.value = null;
  fetchTransactions();
};

const openEditDialog = (transaction) => {
  selectedTransaction.value = transaction;
  showEditDialog.value = true;
};

const confirmDelete = (transaction) => {
  selectedTransaction.value = transaction;
  showDeleteDialog.value = true;
};

const handleDelete = async () => {
  try {
    const response = await transactionService.deleteTransaction(selectedTransaction.value.id);
    if (response.data.status) {
      showToast(response.data.message || 'Transaction deleted successfully');
      await fetchTransactions();
    } else {
      showToast(response.data.message || 'Error deleting transaction', 'error');
    }
  } catch (error) {
    showToast(error.response?.data?.message || 'Error deleting transaction', 'error');
  } finally {
    showDeleteDialog.value = false;
    selectedTransaction.value = null;
  }
};

const changePage = (page) => {
  if (page >= 1 && page <= pagination.value.total_page) {
    filter.value.page = page;
    fetchTransactions();
  }
};

const exportData = async () => {
  try {
    const response = await transactionService.exportTransactions(filter.value);
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `transactions_${new Date().toISOString().split('T')[0]}.xlsx`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
    showToast('Export successful');
  } catch (error) {
    showToast(error.response?.data?.message || 'Error exporting transactions', 'error');
  }
};

const filteredTransactions = computed(() => {
  if (!searchQuery.value) return transactions.value;
  
  const query = searchQuery.value.toLowerCase();
  return transactions.value.filter(transaction => 
    transaction.description?.toLowerCase().includes(query) ||
    transaction.category?.toLowerCase().includes(query) ||
    transaction.amount.toString().includes(query)
  );
});

onMounted(() => {
  fetchTransactions();
  fetchCategories();
});
</script>

<template>
  <UserLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-6">
      <!-- Header Section with Glassmorphism -->
      <div class="bg-white/30 border border-gray-200/50 rounded-2xl p-6 shadow-sm">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
          <div>
            <h1 class="text-2xl font-bold text-black bg-clip-text">
              Transactions
            </h1>
            <p class="text-gray-500 mt-1">Manage your financial activities</p>
          </div>
          
          <div class="flex flex-wrap items-center gap-3">
            <div class="relative">
              <SearchIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <Input 
                v-model="searchQuery"
                type="text"
                placeholder="Search transactions..."
                class="pl-10 w-64 outline-none focus-visible:ring-0 focus-visible:ring-indigo-500 focus-visible:border-indigo-500"
              />
            </div>
            
            <Dialog :open="showAddDialog" @update:open="(value) => showAddDialog = value">
              <DialogTrigger asChild>
                <Button class="bg-gradient-to-r bg-indigo-600 text-white hover:bg-indigo-500 shadow-lg shadow-indigo-500/20">
                  <PlusIcon class="w-4 h-4 mr-2" />
                  New Transaction
                </Button>
              </DialogTrigger>
              <DialogContent class="sm:max-w-[425px]">
                <DialogHeader>
                  <DialogTitle>Add New Transaction</DialogTitle>
                </DialogHeader>
                <AddTransactionModal @transaction-added="handleTransactionAdded" />
              </DialogContent>
            </Dialog>

            <Button 
              @click="exportData"
              class="bg-white text-gray-700 border-2 border-green-500 hover:bg-green-50 shadow-sm"
            >
              <FileSpreadsheetIcon class="w-4 h-4 mr-2" />
              Export
            </Button>
          </div>
        </div>

        <!-- Summary Cards with Modern Design -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6">
          <div class="bg-white rounded-xl p-6 border border-gray-500 shadow-sm">
            <div class="flex items-center justify-between">
              <h3 class="text-gray-600 font-medium">Balance</h3>
              <div class="p-2 bg-gray-100 rounded-lg">
                <MoreHorizontalIcon class="w-4 h-4 text-gray-500" />
              </div>
            </div>
            <p class="text-2xl font-bold mt-2" :class="summary.balance >= 0 ? 'text-green-500' : 'text-red-500'">
              {{ formatCurrency(summary.balance) }}
            </p>
          </div>
          
          <div class="bg-white rounded-xl p-6 border border-green-500">
            <div class="flex items-center justify-between">
              <h3 class="text-gray-600 font-medium">Income</h3>
              <div class="p-2 bg-green-100 rounded-lg">
                <ArrowUpIcon class="w-4 h-4 text-green-500" />
              </div>
            </div>
            <p class="text-2xl font-bold mt-2 text-green-500">
              {{ formatCurrency(summary.total_income) }}
            </p>
          </div>
          
          <div class="bg-white rounded-xl p-6 border border-red-500 shadow-sm">
            <div class="flex items-center justify-between">
              <h3 class="text-gray-600 font-medium">Expenses</h3>
              <div class="p-2 bg-red-100 rounded-lg">
                <ArrowDownIcon class="w-4 h-4 text-red-500" />
              </div>
            </div>
            <p class="text-2xl font-bold mt-2 text-red-500">
              {{ formatCurrency(summary.total_expense) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Filter Section with Collapsible Design -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200">
        <div 
          class="p-4 flex justify-between items-center cursor-pointer"
          @click="isFilterExpanded = !isFilterExpanded"
        >
          <div class="flex items-center gap-2">
            <FilterIcon class="w-4 h-4 text-gray-500" />
            <span class="font-medium">Filters</span>
          </div>
          <ChevronDownIcon
            :class="`w-4 h-4 text-gray-500 transform transition-transform duration-200 ${isFilterExpanded ? 'rotate-180' : ''}`"
          />
        </div>
        
        <div 
          v-show="isFilterExpanded"
          class="p-4 border-t border-gray-100 grid grid-cols-1 md:grid-cols-4 gap-4"
        >
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Start Date</label>
            <Input 
              type="date" 
              v-model="filter.start_date"
              @change="fetchTransactions"
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">End Date</label>
            <Input 
              type="date" 
              v-model="filter.end_date"
              @change="fetchTransactions"
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Category</label>
            <select 
              v-model="filter.category_id"
              class="w-full rounded-lg border border-gray-200 p-2 focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-shadow"
              @change="fetchTransactions"
            >
              <option value="">All Categories</option>
              <option 
                v-for="category in categories" 
                :key="category.id" 
                :value="category.id"
              >
                {{ category.name }}
              </option>
            </select>
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Type</label>
            <select 
              v-model="filter.type"
              class="w-full rounded-lg border border-gray-200 p-2 focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-shadow"
              @change="fetchTransactions"
            >
              <option value="">All Types</option>
              <option value="income">Income</option>
              <option value="expense">Expense</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Transactions Table with Modern Styling -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50/50">
              <tr>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <template v-if="!isLoading && filteredTransactions.length">
                <tr 
                  v-for="transaction in filteredTransactions" 
                  :key="transaction.id"
                  class="hover:bg-gray-50/50 transition-colors"
                >
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    {{ formatDate(parseDate(transaction.date)) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span 
                      class="px-3 py-1 rounded-full text-xs font-medium"
                      :class="transaction.type === 'income' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
                    >
                      {{ transaction.type }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    {{ transaction.category }}
                  </td>
                  <td class="px-6 py-4 text-sm">
                    {{ transaction.description || '-' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium"
                      :class="transaction.type === 'income' ? 'text-green-600' : 'text-red-600'">
                    {{ formatCurrency(transaction.amount) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm">
                          <MoreHorizontalIcon class="w-4 h-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem @click="openEditDialog(transaction)">
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem 
                          @click="confirmDelete(transaction)"
                          class="text-red-600"
                        >
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
                </tr>
              </template>
              <tr v-else-if="isLoading">
                <td colspan="6">
                  <div class="flex items-center justify-center py-8">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
                  </div>
                </td>
              </tr>
              <tr v-else>
                <td colspan="6">
                  <div class="flex flex-col items-center justify-center py-8 text-gray-500">
                    <div class="rounded-full bg-gray-100 p-3 mb-2">
                      <SearchIcon class="w-6 h-6" />
                    </div>
                    <p>No transactions found</p>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Modern Pagination -->
        <div class="px-6 py-4 flex items-center justify-between border-t border-gray-100">
          <p class="text-sm text-gray-700">
            Showing 
            <span class="font-medium">{{ ((pagination.current_page - 1) * pagination.item_per_page) + 1 }}</span>
            to
            <span class="font-medium">{{ Math.min(pagination.current_page * pagination.item_per_page, pagination.total_items) }}</span>
            of
            <span class="font-medium">{{ pagination.total_items }}</span>
            results
          </p>
          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              @click="changePage(pagination.current_page - 1)"
              :disabled="pagination.current_page === 1"
              class="flex items-center gap-1"
            >
              <ChevronLeftIcon class="w-4 h-4" />
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              @click="changePage(pagination.current_page + 1)"
              :disabled="pagination.current_page === pagination.total_page"
              class="flex items-center gap-1"
            >
              Next
              <ChevronRightIcon class="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>

      <!-- Edit Transaction Dialog -->
      <Dialog 
        :open="showEditDialog"
        @update:open="(value) => showEditDialog = value"
      >
        <DialogContent class="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Edit Transaction</DialogTitle>
          </DialogHeader>
          <AddTransactionModal 
            v-if="selectedTransaction"
            :transaction="selectedTransaction"
            @transaction-added="handleTransactionEdited"
          />
        </DialogContent>
      </Dialog>

      <!-- Delete Confirmation Dialog -->
      <AlertDialog :open="showDeleteDialog" @update:open="(value) => showDeleteDialog = value">
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Transaction</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this transaction? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              class="bg-red-600 text-white hover:bg-red-500"
              @click="handleDelete"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  </UserLayout>
</template>