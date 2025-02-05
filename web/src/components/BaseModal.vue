<script>
export default {
  name: 'BaseModal',

  props: {
    show: {
      type: Boolean,
      required: true
    }
  },

  emits: ['close']
}
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
           @click="$emit('close')"></div>

      <!-- Modal -->
      <div class="flex min-h-full items-center justify-center p-4">
        <div class="relative w-full max-w-lg bg-white rounded-lg shadow-xl">
          <!-- Header -->
          <div class="border-b px-6 py-4">
            <h3 class="text-lg font-medium text-gray-900">
              <slot name="title"></slot>
            </h3>
          </div>

          <!-- Content -->
          <div class="px-6 py-4">
            <slot name="content"></slot>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>