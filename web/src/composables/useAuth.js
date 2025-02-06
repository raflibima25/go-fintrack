import { computed } from 'vue'
import { useRouter } from 'vue-router'

export function useAuth() {
  const router = useRouter()

  const isAuthenticated = computed(() => {
    return !!localStorage.getItem('token') !== null
  })

  const isAdmin = computed(() => {
    return localStorage.getItem('isAdmin') === 'true'
  })

  const userRole = computed(() => {
    return isAdmin.value ? 'admin' : 'user'
  })

  const userName = computed(() => {
    return localStorage.getItem('userName') || 'User'
  })

  const checkAuth = (requiredRole = null) => {
    // cek auth
    if (!isAuthenticated.value) {
      router.push({
        name: 'LoginAuth',
        query: { redirect: router.currentRoute.value.fullPath }
      })
      return false
    }

    // cek role
    if (requiredRole && requiredRole !== userRole.value) {
      router.push(isAdmin.value ? '/admin-dashboard' : '/dashboard')
      return false
    }

    return true
  }

  const logout = () => {
    try {
      localStorage.removeItem('token')
      localStorage.removeItem('isAdmin')
      localStorage.removeItem('userName')
      router.push({ name: 'LoginAuth' })
    } catch (error) {
      console.error('Error logout:', error)
    }
  }

  // helper untuk dapat initial user
  const getUserInitials = () => {
    return userName.value
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
  }

  return {
    isAuthenticated,
    isAdmin,
    userRole,
    userName,
    checkAuth,
    logout,
    getUserInitials
  }
}
