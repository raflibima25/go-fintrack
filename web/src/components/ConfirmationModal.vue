<script setup>
import { AlertTriangle } from 'lucide-vue-next'
import BaseModal from './BaseModal.vue'

defineProps({
  show: {
    type: Boolean,
    required: true
  }
})

defineEmits(['close', 'confirm'])
</script>

<template>
  <BaseModal :show="show" @close="$emit('close')">
    <template #title>
      <div class="flex items-center">
        <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-red-100">
          <AlertTriangle class="size-6 text-red-600" aria-hidden="true" />
        </div>
        <h3 class="ml-4 text-base font-semibold text-gray-900">
          <slot name="title">Confirmation</slot>
        </h3>
      </div>
    </template>

    <template #content>
      <div class="space-y-4">
        <div class="ml-14 text-sm text-gray-500">
          <slot name="content">Are you sure?</slot>
        </div>

        <div class="mt-5 flex justify-end gap-3 sm:mt-4">
          <button
            type="button"
            class="inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 shadow-xs ring-gray-300 ring-inset hover:bg-gray-50 sm:w-auto"
            @click="$emit('close')"
          >
            Cancel
          </button>
          <button
            type="button"
            class="inline-flex w-full justify-center rounded-md bg-red-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-red-500 sm:w-auto"
            @click="$emit('confirm')"
          >
            Confirm
          </button>
        </div>
      </div>
    </template>
  </BaseModal>
</template>