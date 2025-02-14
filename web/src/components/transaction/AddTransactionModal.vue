<script setup>
import { ref, onMounted } from 'vue';
import { transactionService } from '@/services/transactionService';
import { categoryService } from '@/services/categoryService';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Calendar } from '@/components/ui/calendar';
import { CalendarIcon } from 'lucide-vue-next';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { formatDate } from '@/utils/formatters';
import { useToast } from '@/composables/useToast';
import { Money3Component } from 'v-money3'

const emit = defineEmits(['transaction-added']);

const form = ref({
  type: 'expense',
  category_id: '',
  amount: '',
  description: '',
  date: new Date()
});

const categories = ref([]);
const isLoading = ref(false);
const error = ref('');
const isDatePickerOpen = ref(false);
const { showToast } = useToast();

// config for v-money3
const moneyConfig = {
  decimal: ',',
  thousands: '.',
  prefix: 'Rp ',
  suffix: '',
  precision: 0,
  masked: false,
  disableNegative: true,
  disabled: false,
  min: 0,
  max: null,
  allowBlank: false,
  minimumNumberOfCharacters: 0,
  modelModifiers: { lazy: true },
  debounce: 0,
};

const fetchCategories = async () => {
  try {
    const response = await categoryService.getCategories();
    if (response.data.status) {
      categories.value = response.data.data.categories;
    }
  } catch (err) {
    console.error('Error fetching categories:', err);
  }
};

const handleSubmit = async () => {
  try {
    const rawAmount = form.value.amount?.toString().replace('Rp ', '').replace(/\./g, '');
    const amount = Number(rawAmount);
    
    if (!amount || amount <= 0) {
      showToast('Please enter a valid amount', 'error');
      return;
    }

    isLoading.value = true;
    error.value = '';

    const formattedDate = form.value.date instanceof Date 
      ? form.value.date.toISOString().split('T')[0]
      : new Date(form.value.date).toISOString().split('T')[0];

    const formData = {
        ...form.value,
        date: formattedDate,
        amount
    }

    const response = await transactionService.createTransaction(formData);
    
    if (response.data.status) {
        // reset form
        form.value = {
            type: 'expense',
            category_id: '',
            amount: '',
            description: '',
            date: new Date()
        };
        emit('transaction-added');
        showToast('Transaction added successfully', 'success');
    }
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to add transaction';
    console.error('Error adding transaction:', err);
    showToast(error.value, 'error');
  } finally {
    isLoading.value = false;
  }
};

const getFormattedDate = (date) => {
  try {
    if (!date) return formatDate(new Date());
    return formatDate(date instanceof Date ? date : new Date(date));
  } catch (err) {
    console.error('Error formatting date:', err);
    return formatDate(new Date());
  }
};

const handleDateSelect = (newDate) => {
    if (newDate) {
        form.value.date = new Date(newDate);
        isDatePickerOpen.value = false;
    }
};

onMounted(() => {
    form.value.date = new Date();
    fetchCategories();
});
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <div class="grid grid-cols-2 gap-4">
      <Button
        type="button"
        :class="[
            'w-full',
            form.type === 'income' 
            ? 'bg-green-600 text-white hover:bg-green-500 border-green-600' 
            : 'bg-white text-green-600 hover:bg-gray-50 border-green-600'
        ]"
        @click="form.type = 'income'"
      >
        Income
      </Button>
      <Button
        type="button"
        :class="[
            'w-full',
            form.type === 'expense' 
            ? 'bg-red-600 text-white hover:bg-red-500 border-red-600' 
            : 'bg-white text-red-600 hover:bg-gray-50 border-red-600'
        ]"
        @click="form.type = 'expense'"
      >
        Expense
      </Button>
    </div>

    <div class="space-y-2">
      <label class="text-sm font-medium">Category</label>
      <select
        v-model="form.category_id"
        class="w-full rounded-md border border-gray-300 p-2 outline-none focus:ring-1 focus:border-indigo-600 focus:ring-indigo-600"
        required
      >
        <option value="" disabled>Select category</option>
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
      <label class="text-sm font-medium">Amount</label>
      <Money3Component 
        v-model="form.amount"
        class="w-full rounded-md border border-gray-300 p-2 outline-none focus:ring-1 focus:border-indigo-600 focus:ring-indigo-600"
        v-bind="moneyConfig"
        required
      />
    </div>

    <div class="space-y-2">
      <label class="text-sm font-medium">Date</label>
      <Popover :open="isDatePickerOpen" @update:open="isDatePickerOpen = $event">
          <PopoverTrigger asChild>
              <div class="relative cursor-pointer">
                  <Input
                      :value="getFormattedDate(form.date)"
                      class="ring-offset-background focus-visible:ring-1 focus-visible:ring-indigo-600 focus-visible:border-indigo-600"
                      readonly
                  />
                  <CalendarIcon class="absolute right-2 top-2.5 h-4 w-4 text-gray-500" />
              </div>
          </PopoverTrigger>
          <PopoverContent class="w-auto p-0 bg-white">
              <Calendar
                  :selected="form.date"
                  :defaultDate="form.date"
                  :value="form.date"
                  @update:modelValue="handleDateSelect"
                  class="border-0"
                  mode="single"
                  :initialFocus="false"
                  :fromDate="new Date(2023, 0)"
                  :toDate="new Date(2025, 11)"
              />
          </PopoverContent>
      </Popover>
    </div>

    <div class="space-y-2">
      <label class="text-sm font-medium">Description (Optional)</label>
      <Input
        v-model="form.description"
        type="text"
        placeholder="Enter description"
        class="ring-offset-background focus-visible:ring-1 focus-visible:border-indigo-600 focus-visible:ring-indigo-600"
      />
    </div>

    <div v-if="error" class="text-sm text-red-600">
      {{ error }}
    </div>

    <div class="flex justify-end gap-4">
      <Button 
            type="submit" 
            :disabled="isLoading"
            class="bg-indigo-600 text-white hover:bg-indigo-500 border border-indigo-600"
        >
        {{ isLoading ? 'Adding...' : 'Add Transaction' }}
      </Button>
    </div>
  </form>
</template>