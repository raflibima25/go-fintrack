<script setup>
import { useToast } from '@/composables/useToast'
import {
  CheckCircle2,
  XCircle,
  AlertCircle,
  Info
} from 'lucide-vue-next'

const { toasts } = useToast()

const getIcon = (type) => {
  const icons = {
    success: CheckCircle2,
    error: XCircle,
    warning: AlertCircle,
    info: Info
  }
  return icons[type]
}
</script>

<template>
  <div class="fixed bottom-4 right-4 z-50 space-y-2">
    <TransitionGroup name="toast">
      <div
          v-for="toast in toasts"
          :key="toast.id"
          class="flex items-center gap-3 px-7 py-6 rounded-lg shadow-lg max-w-sm"
          :class="{
          'bg-green-500 text-white': toast.type === 'success',
          'bg-red-500 text-white': toast.type === 'error',
          'bg-yellow-500 text-white': toast.type === 'warning',
          'bg-blue-500 text-white': toast.type === 'info'
        }"
      >
        <component 
          :is="getIcon(toast.type)"
          class="w-5 h-5 flex-shrink-0"
        />
        <span class="text-sm">{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

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