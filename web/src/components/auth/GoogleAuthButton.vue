<script setup>
import { ref } from 'vue';
import { useToast } from '@/composables/useToast';
import { useRouter } from 'vue-router';
import apiClient from '../../utils/api';

const props = defineProps({
    text: {
        type: String,
        default: 'Continue with Google'
    },
    isLoading: {
        type: Boolean,
        default: false
    }
});

const router = useRouter();
const { showToast } = useToast();
const loading = ref(false);

const handleGoogleLogin = async () => {
    try {
        loading.value = true;
        // get login url from backend
        const response = await apiClient.get('/auth/google/login');

        if (response.data.status) {
            localStorage.setItem('googleAuthState', response.data.data.state);
            window.location.href = response.data.data.redirect_url;
        } else {
            throw new Error('Failed to initiate Google login');
        }
    } catch (error) {
        console.log('Google login error:', error);
        showToast(
            error.response?.data?.message || 'Failed to connect with Google',
            'error'
        );
    } finally {
        loading.value = false;
    }
};
</script>

<template>
    <button
        @click="handleGoogleLogin"
        :disabled="loading || isLoading"
        class="w-full flex items-center justify-center gap-3 px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
    >
        <img src="../../assets/google-logo.svg" alt="Google Logo" class="h-5 w-5" />
        <span class="text-gray-700 font-medium">
            {{ loading ? 'Loading...' : text }}
        </span>
    </button>
</template>