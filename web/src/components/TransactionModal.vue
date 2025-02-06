<script>
import { ref, computed, watch } from 'vue'
import BaseModal from '@/components/BaseModal.vue'
import CategoryFormModal from '@/components/CategoryFormModal.vue'
import { useToast } from '@/composables/useToast'
import { categoryService } from '@/services/categoryService'

export default {
  name: 'TransactionModal',

  components: {
    BaseModal,
    CategoryFormModal
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

  emits: ['close', 'submit', 'category-added'],

  setup(props, { emit }) {
    const { showToast } = useToast()
    const showCategoryFormModal = ref(false)
    const isSubmitting = ref(false)

    // Form state
    const form = ref({
      type: 'expense',
      category_id: '',
      amount: '',
      date: new Date().toISOString().split('T')[0],
      description: ''
    })

    // Computed
    const isEdit = computed(() => !!props.transaction)

    // Watchers
    watch(() => props.transaction, (newVal) => {
      if (newVal) {
        // Edit mode - Fill form with transaction data
        form.value = {
          type: newVal.type,
          category_id: newVal.category_id,
          amount: newVal.amount,
          date: newVal.date.split('T')[0],
          description: newVal.description
        }
      } else {
        // Create mode - Reset form
        resetForm()
      }
    }, { immediate: true })

    // Watch untuk categories
    watch(
      () => props.categories,
      (newCategories, oldCategories) => {
        // Inisialisasi awal
        if (!oldCategories && newCategories.length > 0 && !form.value.category_id) {
          form.value.category_id = newCategories[0].id
          return
        }

        // Handle kategori baru ditambahkan
        if (oldCategories && newCategories.length > oldCategories.length) {
          const newCategory = newCategories[newCategories.length - 1]
          if (newCategory) {
            form.value.category_id = newCategory.id
          }
        }
      },
      { immediate: true, deep: true }
    )

    // Methods
    const resetForm = () => {
      form.value = {
        type: 'expense',
        category_id: props.categories[0]?.id || '',
        amount: '',
        date: new Date().toISOString().split('T')[0],
        description: ''
      }
    }

    const closeModal = () => {
      if (!isSubmitting.value) {
        resetForm()
        emit('close')
      }
    }

    const handleSubmit = () => {
      // Validasi
      if (!form.value.category_id || !form.value.amount || !form.value.date) {
        showToast('Please fill in all required fields', 'error')
        return
      }

      isSubmitting.value = true

      try {
        // Format data
        const formData = {
          ...form.value,
          amount: Number(form.value.amount),
          date: new Date(form.value.date).toISOString().split('T')[0]
        }

        emit('submit', formData)
      } catch (error) {
        showToast('Error processing form data', 'error')
      } finally {
        isSubmitting.value = false
      }
    }

    const handleCategorySubmit = async (categoryData, { onSuccess, onError }) => {
      try {
        const response = await categoryService.createCategory(categoryData)

        if (response.data.status) {
          showToast(response.data.message || 'Category created successfully')
          
          // Update categories list
          await emit('category-added')
          
          // Set new category as selected
          if (response.data.data?.id) {
            form.value.category_id = response.data.data.id
          }
          
          onSuccess()
          showCategoryFormModal.value = false
        } else {
          onError(response.data.message)
          showToast(response.data.message, 'error')
        }
      } catch (error) {
        const errorMessage = error.response?.data?.message || 'Error creating category'
        onError(errorMessage)
        showToast(errorMessage, 'error')
      }
    }

    return {
      form,
      isEdit,
      isSubmitting,
      showCategoryFormModal,
      closeModal,
      handleSubmit,
      handleCategorySubmit
    }
  }
}
</script>

<template>
  <div>
    <!-- Main Modal -->
    <BaseModal :show="show" @close="closeModal">
      <template #title>
        {{ isEdit ? 'Edit Transaction' : 'Create New Transaction' }}
      </template>
      <template #content>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- Type Field -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <select v-model="form.type" required
                    class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500">
              <option value="income">Income</option>
              <option value="expense">Expense</option>
            </select>
          </div>

          <!-- Category Field -->
          <div>
            <div class="flex items-center justify-between">
              <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
              <div class="text-sm">
                <button 
                  type="button"
                  @click="showCategoryFormModal = true" 
                  class="font-semibold text-indigo-600 hover:text-indigo-500"
                >
                  Add category
                </button>
              </div>
            </div>
            <select v-model="form.category_id" required
                    class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500">
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>

          <!-- Amount Field -->
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

          <!-- Date Field -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
            <input
                type="date"
                v-model="form.date"
                required
                class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
            >
          </div>

          <!-- Description Field -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
                v-model="form.description"
                rows="3"
                class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter description"
            ></textarea>
          </div>

          <!-- Form Actions -->
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
                :disabled="isSubmitting"
                class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50"
            >
              {{ isEdit ? 'Update' : 'Create' }}
            </button>
          </div>
        </form>
      </template>
    </BaseModal>

    <!-- Category Modal -->
    <CategoryFormModal
      v-if="showCategoryFormModal"
      :show="showCategoryFormModal"
      :category="null"
      @close="showCategoryFormModal = false"
      @submit="handleCategorySubmit"
    />
  </div>
</template>