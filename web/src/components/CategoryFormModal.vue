<script>
import { ref, computed, watch } from 'vue'
import BaseModal from './BaseModal.vue'
import { Tag } from 'lucide-vue-next'

export default {
  name: 'CategoryFormModal',

  components: {
    BaseModal,
    Tag
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
    // Predefined colors
    const backgroundColors = [
      { value: 'bg-blue-100' },
      { value: 'bg-green-100' },
      { value: 'bg-yellow-100' },
      { value: 'bg-red-100' },
      { value: 'bg-purple-100' },
      { value: 'bg-pink-100' },
      { value: 'bg-indigo-100' },
      { value: 'bg-gray-100' }
    ]

    const iconColors = [
      { value: 'text-blue-500' },
      { value: 'text-green-500' },
      { value: 'text-yellow-500' },
      { value: 'text-red-500' },
      { value: 'text-purple-500' },
      { value: 'text-pink-500' },
      { value: 'text-indigo-500' },
      { value: 'text-gray-500' }
    ]

    const form = ref({
      name: '',
      color: 'bg-blue-100',
      icon_color: 'text-blue-500'
    })
    
    const error = ref('')
    const isSubmitting = ref(false)
    const isEdit = computed(() => !!props.category)

    watch(() => props.category, (newVal) => {
      if (newVal) {
        form.value = {
          name: newVal.name,
          color: newVal.color || 'bg-blue-100',
          icon_color: newVal.icon_color || 'text-blue-500'
        }
      } else {
        form.value = {
          name: '',
          color: 'bg-blue-100',
          icon_color: 'text-blue-500'
        }
      }
      error.value = ''
    }, { immediate: true })

    const closeModal = () => {
      if (!isSubmitting.value) {
        error.value = ''
        emit('close')
      }
    }

    const handleSubmit = async () => {
      if (!form.value.name.trim()) {
        error.value = 'Category name is required'
        return
      }

      isSubmitting.value = true
      try {
        await emit('submit', {
          name: form.value.name.trim(),
          color: form.value.color,
          icon_color: form.value.icon_color
        })
        error.value = ''
      } catch (err) {
        error.value = err.message || 'Error saving category'
      } finally {
        isSubmitting.value = false
      }
    }

    return {
      form,
      error,
      isEdit,
      isSubmitting,
      backgroundColors,
      iconColors,
      closeModal,
      handleSubmit
    }
  }
}
</script>

<template>
  <BaseModal :show="show" @close="closeModal">
    <template #title>
      {{ isEdit ? 'Edit Category' : 'Create New Category' }}
    </template>
    <template #content>
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <!-- Name Field -->
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
          <p v-if="error" class="mt-1 text-sm text-red-500">{{ error }}</p>
        </div>

        <!-- Background Color -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Background Color</label>
          <div class="grid grid-cols-8 gap-3">
            <button
              v-for="color in backgroundColors"
              :key="color.value"
              type="button"
              @click="form.color = color.value"
              :class="[
                'w-8 h-8 rounded-lg border-2',
                color.value,
                form.color === color.value ? 'border-blue-500' : 'border-transparent'
              ]"
            />
          </div>
        </div>

        <!-- Icon Color -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Icon Color</label>
          <div class="grid grid-cols-8 gap-3">
            <button
              v-for="color in iconColors"
              :key="color.value"
              type="button"
              @click="form.icon_color = color.value"
              class="w-8 h-8 rounded-lg border-2 flex items-center justify-center"
              :class="[form.icon_color === color.value ? 'border-blue-500' : 'border-transparent']"
            >
              <Tag :class="[color.value]" />
            </button>
          </div>
        </div>

        <!-- Submit Buttons -->
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
            {{ isSubmitting ? 'Saving...' : (isEdit ? 'Update' : 'Create') }}
          </button>
        </div>
      </form>
    </template>
  </BaseModal>
</template>