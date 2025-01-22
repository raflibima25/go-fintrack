# ConfirmationModal.vue
<template>
  <BaseModal :show="show" @close="$emit('close')">
    <template #title>
      <slot name="title">Confirmation</slot>
    </template>
    <template #content>
      <div class="space-y-4">
        <div class="text-gray-600">
          <slot name="content">Are you sure?</slot>
        </div>

        <div class="flex justify-end gap-3">
          <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 border rounded-lg hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
              type="button"
              @click="$emit('confirm')"
              class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
          >
            Confirm
          </button>
        </div>
      </div>
    </template>
  </BaseModal>
</template>

<script>
import BaseModal from './BaseModal.vue'

export default {
  name: 'ConfirmationModal',

  components: {
    BaseModal
  },

  props: {
    show: {
      type: Boolean,
      required: true
    }
  },

  emits: ['close', 'confirm']
}
</script>