import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from './useToast'
import apiClient from '../utils/api'

export function useGoogleAuth() {
  const router = useRouter()
  const { showToast } = useToast()
  const loading = ref(false)

  const handleCallback = async (code, state) => {
    try {
      loading.value = true
      console.log('Processing callback with:', { code, state }) // debug

      const savedState = localStorage.getItem('googleAuthState')
      if (state !== savedState) {
        console.log('State mismatch:', { saved: savedState, received: state }) // debug
        throw new Error('Invalid state parameter')
      }

      const response = await apiClient.get('/auth/google/callback', {
        params: { code, state }
      })

      console.log('Callback response:', response.data) // debug

      if (response.data.status) {
        localStorage.setItem('token', response.data.data.access_token)
        localStorage.setItem('isAdmin', response.data.data.is_admin || false)

        showToast('Successfully logged in with Google', 'success')

        const redirectPath = response.data.data.is_admin ? '/admin/dashboard' : '/dashboard'
        await router.push(redirectPath)
      } else {
        throw new Error(response.data.message || 'Authentication failed')
      }
    } catch (error) {
      console.error('Google callback error:', error)
      showToast(
        error.response?.data?.message || 'Failed to complete Google authentication',
        'error'
      )
      router.push('/login')
    } finally {
      loading.value = false
      localStorage.removeItem('googleAuthState')
    }
  }

  return {
    loading,
    handleCallback
  }
}
