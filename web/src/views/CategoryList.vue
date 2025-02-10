<script setup>
import { ref, onMounted, computed } from 'vue'
import { useToast } from '@/composables/useToast'
import { categoryService } from '@/services/categoryService'
import { Tag, Edit, Trash2 } from 'lucide-vue-next'
import CategoryFormModal from '@/components/CategoryFormModal.vue'
import ConfirmationModal from '@/components/ConfirmationModal.vue'
import UserLayout from '../layouts/UserLayout.vue'

const { showToast } = useToast()
const categories = ref([])
const showModal = ref(false)
const showDeleteModal = ref(false)
const selectedCategory = ref(null)

// Computed
const mostUsedCategory = computed(() => {
  return [...categories.value].sort((a, b) => 
    (b.usage_count || 0) - (a.usage_count || 0)
  )[0]
})

// Methods
const fetchCategories = async () => {
  try {
    const response = await categoryService.getCategories()
    if (response.data.status) {
      categories.value = response.data.data.categories
    } else {
      showToast(response.data.message || 'Error fetching categories', 'error')
    }
  } catch (error) {
    showToast(error.response?.data?.message || 'Error fetching categories', 'error')
  }
}

const openCreateModal = () => {
  selectedCategory.value = null
  showModal.value = true
}

const editCategory = (category) => {
  selectedCategory.value = category
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedCategory.value = null
}

const handleSubmit = async (formData) => {
  try {
    let response
    if (selectedCategory.value) {
      response = await categoryService.updateCategory(
        selectedCategory.value.id, 
        {
          name: formData.name,
          color: formData.color,
          icon_color: formData.icon_color
        }
      )
    } else {
      response = await categoryService.createCategory({
        name: formData.name,
        color: formData.color,
        icon_color: formData.icon_color
      })
    }

    if (response.data.status) {
      showToast(response.data.message || 'Category saved successfully')
      await fetchCategories()
      closeModal()
    } else {
      showToast(response.data.message || 'Error saving category', 'error')
    }
  } catch (error) {
    showToast(
      error.response?.data?.message || 'Error saving category', 
      'error'
    )
  }
}

const confirmDelete = (category) => {
  selectedCategory.value = category
  showDeleteModal.value = true
}

const handleDelete = async () => {
  try {
    const response = await categoryService.deleteCategory(
      selectedCategory.value.id
    )
    if (response.data.status) {
      showToast(response.data.message || 'Category deleted successfully')
      await fetchCategories()
      showDeleteModal.value = false
    } else {
      showToast(response.data.message || 'Error deleting category', 'error')
    }
  } catch (error) {
    showToast(
      error.response?.data?.message || 'Error deleting category', 
      'error'
    )
  }
}

const formatPercentage = (value) => {
  return Number(value || 0).toFixed(1)
}

onMounted(() => {
  fetchCategories()
})
</script>

<template>
  <user-layout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="py-6">
        <!-- Header Section -->
        <div class="flex justify-between items-center mb-6">
          <h1 class="text-2xl font-bold">Categories</h1>
          <button 
            @click="openCreateModal"
            class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600"
          >
            Add Category
          </button>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          <div class="bg-white p-4 rounded-lg shadow">
            <h3 class="text-lg font-semibold text-gray-600">Total Categories</h3>
            <p class="text-2xl font-bold text-blue-500">{{ categories.length }}</p>
          </div>
          <div class="bg-white p-4 rounded-lg shadow">
            <h3 class="text-lg font-semibold text-gray-600">Most Used</h3>
            <p class="text-2xl font-bold text-blue-500">
              {{ mostUsedCategory?.name || 'N/A' }}
            </p>
            <p class="text-sm text-gray-500">
              {{ mostUsedCategory?.usage_count || 0 }} transactions
            </p>
          </div>
          <div class="bg-white p-4 rounded-lg shadow">
            <h3 class="text-lg font-semibold text-gray-600">Last Added</h3>
            <p class="text-2xl font-bold text-blue-500">
              {{ categories[categories.length - 1]?.name || 'N/A' }}
            </p>
          </div>
        </div>

        <!-- Categories Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <div v-for="category in categories" :key="category.id" 
              class="bg-white rounded-lg shadow p-4 flex flex-col">
            <div class="flex items-center justify-between mb-4">
              <!-- Category Icon/Color -->
              <div class="flex items-center">
                <div :class="['w-10 h-10 rounded-full flex items-center justify-center', category.color || 'bg-blue-100']">
                  <Tag :class="['w-5 h-5', category.icon_color || 'text-blue-500']" />
                </div>
                <span class="ml-3 font-medium">{{ category.name }}</span>
              </div>
              
              <!-- Actions -->
              <div class="flex items-center space-x-2">
                <button @click="editCategory(category)" 
                        class="text-blue-500 hover:text-blue-700">
                  <Edit class="w-4 h-4" />
                </button>
                <button @click="confirmDelete(category)"
                        class="text-red-500 hover:text-red-700">
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </div>

            <!-- Usage Stats -->
            <div class="mt-auto">
              <div class="text-sm text-gray-500">Usage</div>
              <div class="flex items-center justify-between">
                <div class="text-sm font-medium">
                  {{ category.usage_count || 0 }} transactions
                </div>
                <div class="text-sm text-gray-500">
                  {{ formatPercentage(category.usage_percentage) }}%
                </div>
              </div>
              <!-- Progress bar -->
              <div class="w-full h-2 bg-gray-200 rounded-full mt-1">
                <div class="h-2 bg-blue-500 rounded-full"
                    :style="{ width: (category.usage_percentage || 0) + '%' }" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Category Modal -->
      <CategoryFormModal
        v-if="showModal"
        :show="showModal"
        :category="selectedCategory"
        @close="closeModal"
        @submit="handleSubmit"
      />

      <!-- Delete Confirmation Modal -->
      <ConfirmationModal
        v-if="showDeleteModal"
        :show="showDeleteModal"
        @close="showDeleteModal = false"
        @confirm="handleDelete"
      >
        <template #title>Delete Category</template>
        <template #content>
          Are you sure you want to delete this category? This will not delete any transactions, but they will no longer be associated with this category.
        </template>
      </ConfirmationModal>
    </div>
  </user-layout>
</template>