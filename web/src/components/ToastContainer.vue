<template>
  <div class="fixed bottom-4 right-4 z-50 space-y-2">
    <TransitionGroup name="toast">
      <div
          v-for="toast in toasts"
          :key="toast.id"
          class="px-4 py-2 rounded-lg shadow-lg max-w-sm"
          :class="{
          'bg-green-500 text-white': toast.type === 'success',
          'bg-red-500 text-white': toast.type === 'error',
          'bg-yellow-500 text-white': toast.type === 'warning',
          'bg-blue-500 text-white': toast.type === 'info'
        }"
      >
        {{ toast.message }}
      </div>
    </TransitionGroup>
  </div>
</template>

<script>
import { useToast } from '@/composables/useToast'

export default {
  name: 'ToastContainer',

  setup() {
    const { toasts } = useToast()
    return { toasts }
  }
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(30px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100px);
}
</style>