# CategoryModal.vue
<template>
  <BaseModal :show="show" @close="closeModal">
    <template #title>
      {{ isEdit ? 'Edit Category' : 'Create New Category' }}
    </template>
    <template #content>
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input
              type="text"
              v-model="form.name"
              required
              class="w-full border rounded-lg px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter category name"
              :class="{ 'border-red-500': error }"
          >
          <!-- Error message -->
          <p v-if="error" class="mt-1 text-sm text-red-500">{{ error }}</p>
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
              :disabled="isSubmitting"
              class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isSubmitting ? 'Saving...' : (isEdit ? 'Update' : 'Create') }}
          </button>
        </div>
      </form>
    </template>
  </BaseModal>
</template>

<script>
import { ref, computed, watch } from 'vue'
import BaseModal from './BaseModal.vue'

export default {
  name: 'CategoryModal',

  components: {
    BaseModal
  },

  props: {
    show: {
      type: Boolean,
      required: true
    },
    category: {
      type: Object,
      default: null
    }
  },

  emits: ['close', 'submit'],

  setup(props, { emit }) {
    const form = ref({
      name: '',
      color: 'bg-blue-100',     
      icon_color: 'text-blue-500'
    })
    
    const error = ref('')
    const isSubmitting = ref(false)
    
    const isEdit = computed(() => !!props.category)

    // Watch untuk mengisi form ketika editing
    watch(() => props.category, (newVal) => {
      if (newVal) {
        form.value = {
          name: newVal.name,
          color: newVal.color || 'bg-blue-100',
          icon_color: newVal.icon_color || 'text-blue-500'
        }
      } else {
        // Reset form ketika membuat baru
        form.value = {
          name: '',
          color: 'bg-blue-100',
          icon_color: 'text-blue-500'
        }
      }
      // Reset error saat category berubah
      error.value = ''
    }, { immediate: true })

    const closeModal = () => {
      if (!isSubmitting.value) {
        error.value = ''
        form.value = {
          name: '',
          color: 'bg-blue-100',
          icon_color: 'text-blue-500'
        }
        emit('close')
      }
    }

    const handleSubmit = () => {
      // Reset error sebelum validasi
      error.value = ''

      // Validasi
      if (!form.value.name.trim()) {
        error.value = 'Category name is required'
        return
      }

      // Set submitting state
      isSubmitting.value = true

      // Emit data ke parent
      emit('submit', {
        name: form.value.name.trim(),
        color: form.value.color,
        icon_color: form.value.icon_color
      }, {
        onSuccess: () => {
          isSubmitting.value = false
          closeModal()
        },
        onError: (errorMessage) => {
          isSubmitting.value = false
          error.value = errorMessage
        }
      })
    }

    // Watch untuk show prop
    watch(() => props.show, (newVal) => {
      if (!newVal) {
        // Reset form dan error saat modal ditutup
        error.value = ''
        form.value.name = ''
      }
    })

    return {
      form,
      error,
      isEdit,
      isSubmitting,
      closeModal,
      handleSubmit
    }
  }
}
</script>