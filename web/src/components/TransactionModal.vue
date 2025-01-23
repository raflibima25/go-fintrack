<template>
  <BaseModal :show="show" @close="closeModal">
    <template #title>
      {{ isEdit ? 'Edit Transaction' : 'Create New Transaction' }}
    </template>
    <template #content>
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select v-model="form.type" required
                  class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500">
            <option value="income">Income</option>
            <option value="expense">Expense</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select v-model="form.category_id" required
                  class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500">
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Amount</label>
          <input
              type="number"
              v-model="form.amount"
              required
              min="0"
              step="0.01"
              class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter amount"
          >
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
          <input
              type="date"
              v-model="form.date"
              required
              class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
          >
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea
              v-model="form.description"
              rows="3"
              class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter description"
          ></textarea>
        </div>

        <div class="flex justify-end gap-3 pt-4">
          <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 border rounded-lg hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
              type="submit"
              class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
          >
            {{ isEdit ? 'Update' : 'Create' }}
          </button>
        </div>
      </form>
    </template>
  </BaseModal>
</template>

<script>
import { ref, computed, watch } from 'vue'
import BaseModal from '@/components/BaseModal.vue'

export default {
  name: 'TransactionModal',

  components: {
    BaseModal
  },

  props: {
    show: {
      type: Boolean,
      required: true
    },
    categories: {
      type: Array,
      required: true
    },
    transaction: {
      type: Object,
      default: null
    }
  },

  emits: ['close', 'submit'],

  setup(props, { emit }) {
    const form = ref({
      type: 'expense',
      category_id: '',
      amount: '',
      date: new Date().toISOString().split('T')[0],
      description: ''
    })

    const isEdit = computed(() => !!props.transaction)

    // Watch untuk mengisi form ketika editing
    watch(() => props.transaction, (newVal) => {
      if (newVal) {
        form.value = {
          type: newVal.type,
          category_id: newVal.category_id,
          amount: newVal.amount,
          date: newVal.date.split('T')[0], // Mengambil hanya tanggal dari ISO string
          description: newVal.description
        }
      } else {
        // Reset form ketika membuat baru
        form.value = {
          type: 'expense',
          category_id: props.categories[0]?.id || '',
          amount: '',
          date: new Date().toISOString().split('T')[0],
          description: ''
        }
      }
    }, { immediate: true })

    const closeModal = () => {
      emit('close')
    }

    const handleSubmit = () => {
      // Validasi basic
      if (!form.value.category_id || !form.value.amount || !form.value.date) {
        return
      }

      // Convert amount ke number dan format tanggal
      const formData = {
        ...form.value,
        amount: Number(form.value.amount),
        date: new Date(form.value.date).toISOString().split('T')[0] // Format tanggal sesuai API
      }

      emit('submit', formData)
    }

    return {
      form,
      isEdit,
      closeModal,
      handleSubmit
    }
  }
}
</script>