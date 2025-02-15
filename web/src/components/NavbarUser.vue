<script setup>
import { computed, ref } from 'vue'
import { useAuth } from '../composables/useAuth'
import { 
  Menu as MenuIcon, 
  X, 
  LayoutDashboard, 
  Receipt, 
  Tag, 
  BotMessageSquare 
} from 'lucide-vue-next'
import logoFintrack from '@/assets/logo-fintrack-new.webp'

const { userName, logout } = useAuth()
const showMobileMenu = ref(false)
const showUserMenu = ref(false)

// Tidak perlu menyimpan dalam variabel terpisah karena akan langsung digunakan
const userInitials = computed(() => {
  const name = userName.value || 'User'
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
})

const displayName = computed(() => {
  return userName.value || 'User'
})

const handleLogout = async () => {
  await logout()
}
</script>

<template>
  <nav class="bg-white border-b">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <!-- Left section: Logo & Main Navigation -->
        <div class="flex">
          <!-- Logo -->
          <div class="flex-shrink-0 flex items-center">
            <router-link to="/" class="text-2xl font-bold italic text-blue-600 font-montserrat">
              <img class="h-14 w-auto" :src=logoFintrack alt="Fintrack Logo">
            </router-link>
          </div>

          <!-- Main Navigation -->
          <div class="hidden sm:ml-8 sm:flex sm:space-x-4">
            <router-link 
              to="/dashboard" 
              class="inline-flex items-center px-1 pt-1 text-sm font-medium"
              :class="[$route.path === '/dashboard' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700']"
            >
              <LayoutDashboard class="w-4 h-4 mr-2" />
              Dashboard
            </router-link>

            <router-link 
              to="/transactions" 
              class="inline-flex items-center px-1 pt-1 text-sm font-medium"
              :class="[$route.path === '/transactions' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700']"
            >
              <Receipt class="w-4 h-4 mr-2" />
              Transactions
            </router-link>

            <router-link 
              to="/categories" 
              class="inline-flex items-center px-1 pt-1 text-sm font-medium"
              :class="[$route.path === '/categories' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700']"
            >
              <Tag class="w-4 h-4 mr-2" />
              Categories
            </router-link>

            <router-link 
              to="/chat-assistant" 
              class="inline-flex items-center px-1 pt-1 text-sm font-medium"
              :class="[$route.path === '/chat-assistant' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700']"
            >
              <BotMessageSquare class="w-4 h-4 mr-2" />
              AI Assistant
            </router-link>
          </div>
        </div>

        <!-- Right section: User menu -->
        <div class="flex items-center">
          <!-- User dropdown -->
          <div class="ml-3 relative">
            <div>
              <button
                @click="showUserMenu = !showUserMenu"
                class="flex text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <span class="inline-flex items-center justify-center h-8 w-8 rounded-full bg-blue-100">
                  <span class="text-sm font-medium leading-none text-blue-700">
                    {{ userInitials }}
                  </span>
                </span>
              </button>
            </div>

            <!-- Dropdown menu -->
            <div
              v-if="showUserMenu"
              class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 divide-y divide-gray-100 focus:outline-none"
            >
              <div class="py-1">
                <div class="px-4 py-2 text-sm text-gray-700">
                  {{ displayName }}
                </div>
              </div>
              <div class="py-1">
                <button
                  @click="handleLogout"
                  class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100"
                >
                  Logout
                </button>
              </div>
            </div>
          </div>

          <!-- Mobile menu button -->
          <div class="flex items-center sm:hidden">
            <button
              @click="showMobileMenu = !showMobileMenu"
              class="inline-flex items-center justify-center p-2 rounded-md text-blue-500 hover:text-blue-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
            >
              <MenuIcon v-if="!showMobileMenu" class="block h-6 w-6" />
              <X v-else class="block h-6 w-6" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile menu -->
    <div v-if="showMobileMenu" class="sm:hidden">
      <div class="pt-2 pb-3 space-y-1">
        <router-link
          to="/dashboard"
          class="flex items-center px-3 py-2 text-base font-medium"
          :class="[$route.path === '/dashboard' ? 'text-blue-600 bg-blue-50' : 'text-gray-600 hover:bg-gray-50']"
        >
          <LayoutDashboard class="w-5 h-5 mr-3" />
          Dashboard
        </router-link>

        <router-link
          to="/transactions"
          class="flex items-center px-3 py-2 text-base font-medium"
          :class="[$route.path === '/transactions' ? 'text-blue-600 bg-blue-50' : 'text-gray-600 hover:bg-gray-50']"
        >
          <Receipt class="w-5 h-5 mr-3" />
          Transactions
        </router-link>

        <router-link
          to="/categories"
          class="flex items-center px-3 py-2 text-base font-medium"
          :class="[$route.path === '/categories' ? 'text-blue-600 bg-blue-50' : 'text-gray-600 hover:bg-gray-50']"
        >
          <Tag class="w-5 h-5 mr-3" />
          Categories
        </router-link>

        <router-link
          to="/chat-assistant"
          class="flex items-center px-3 py-2 text-base font-medium"
          :class="[$route.path === '/chat-assistant' ? 'text-blue-600 bg-blue-50' : 'text-gray-600 hover:bg-gray-50']"
        >
          <BotMessageSquare class="w-5 h-5 mr-3" />
          AI Assistant
        </router-link>
      </div>
    </div>
  </nav>
</template>