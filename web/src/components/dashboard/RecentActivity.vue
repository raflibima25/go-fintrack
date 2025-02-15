// components/dashboard/RecentActivity.vue
<script setup>
import { ref, onMounted } from 'vue';
import { PlusIcon } from 'lucide-vue-next';
import { 
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { transactionService } from '@/services/transactionService';
import { Button } from '@/components/ui/button';
import AddTransactionModal from '../transaction/AddTransactionModal.vue';
import { formatDate, formatCurrency, parseDate } from '@/utils/formatters';

const transactions = ref([]);
const isLoading = ref(false);
const showAddDialog = ref(false);

const fetchRecentTransactions = async () => {
  try {
    isLoading.value = true;
    const response = await transactionService.getTransactions({
      page: 1,
      limit: 5,
    });
    
    if (response.data.status) {
      transactions.value = response.data.data.transactions;
    }
  } catch (error) {
    console.error('Error fetching recent transactions:', error);
  } finally {
    isLoading.value = false;
  }
};

const handleTransactionAdded = () => {
  showAddDialog.value = false;
  fetchRecentTransactions();
};

onMounted(() => {
  fetchRecentTransactions();
});
</script>

<template>
  <div class="bg-white rounded-lg shadow p-4">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-xl font-semibold">Recent Activity</h2>
      
      <Dialog 
        :open="showAddDialog"
         @update:open="(value) => showAddDialog = value"
        >
        <DialogTrigger asChild>
          <Button class="flex items-center gap-2 bg-indigo-600 text-white hover:bg-indigo-500">
            <PlusIcon class="w-4 h-4" />
            Add Transaction
          </Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Add New Transaction</DialogTitle>
          </DialogHeader>
          <!-- Component Add Transaction -->
          <AddTransactionModal @transaction-added="handleTransactionAdded" />
        </DialogContent>
      </Dialog>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b">
            <th class="px-3 py-1 text-left font-medium">Date</th>
            <th class="px-3 py-1 text-left font-medium">Category</th>
            <th class="px-3 py-1 text-left font-medium hidden sm:table-cell">Description</th>
            <th class="px-3 py-1 text-right font-medium">Amount</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="!isLoading && transactions.length">
            <tr v-for="transaction in transactions" 
                :key="transaction.id"
                class="border-b hover:bg-gray-50">
              <td class="px-3 py-2">{{ formatDate(parseDate(transaction.date)) }}</td>
              <td class="px-3 py-2">{{ transaction.category }}</td>
              <td class="px-3 py-2 hidden sm:table-cell">{{ transaction.description || '-' }}</td>
              <td :class="[
                'px-3 py-2 text-right font-medium',
                transaction.type === 'income' ? 'text-green-600' : 'text-red-600'
              ]">
                {{ formatCurrency(transaction.amount) }}
              </td>
            </tr>
          </template>
          <tr v-else-if="isLoading">
            <td colspan="4" class="px-4 py-8 text-center text-gray-500">
              Loading...
            </td>
          </tr>
          <tr v-else>
            <td colspan="4" class="px-4 py-8 text-center text-gray-500">
              No recent transactions
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>