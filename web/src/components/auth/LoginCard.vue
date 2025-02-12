<script setup>
import { ref, reactive } from "vue";
import { useRouter } from 'vue-router';
import { useToast } from "@/composables/useToast";
import apiClient from "@/utils/api";
import { EyeIcon, EyeOffIcon } from "lucide-vue-next";
import GoogleAuthButton from "./GoogleAuthButton.vue";

const router = useRouter();
const { showToast } = useToast();

const formState = reactive({
    identifier: '',
    password: '',
    isLoading: false,
    errors: {
        identifier: '',
        password: '',
        general: ''
    }
});

const showPassword = ref(false);

const validateForm = () => {
    let isValid = true;
    formState.errors = {
        identifier: '',
        password: '',
        general: ''
    };

    if (!formState.identifier.trim()) {
        formState.errors.identifier = 'Username or email is required';
        isValid = false;
    }

    if (!formState.password) {
        formState.errors.password = 'Password is required';
        isValid = false;
    }

    if (formState.password && formState.password.length < 8) {
        formState.errors.password = 'Password must be at least 8 characters';
        isValid = false;
    }

    return isValid;
};

const togglePasswordVisibility = () => {
    showPassword.value = !showPassword.value;
};

const handleLogin = async (e) => {
    e?.preventDefault();

    if (!validateForm()) return false;

    formState.isLoading = true;
    formState.errors.general = '';

    try {
        const response = await apiClient.post('/auth/login', {
            email_or_username: formState.identifier.trim(),
            password: formState.password
        });

        console.log('Login response:', response.data); // Untuk debugging

        if (response.data.status) {
            localStorage.setItem("token", response.data.data.access_token);
            localStorage.setItem("isAdmin", response.data.data.is_admin);
            
            showToast("Login successful", "success");

            // redirect base on role
            const redirectPath = response.data.data.is_admin ? "/admin-dashboard" : "/dashboard";
            await router.push(redirectPath);
        } else {
            throw new Error(response.data.message || "Login failed. Please try again.");
        }
    } catch (error) {
        console.error('Login error:', error); // Untuk debugging

        let errorMessage = "An unexpected error occurred";
        
        if (error.response) {
            errorMessage = error.response.data.message || "Login failed. Please try again.";
            if (error.response.status === 401) {
                errorMessage = "Invalid email/username or password";
            } else if (error.response.status === 400) {
                errorMessage = error.response.data.message || "Invalid input";
            }
        } else if (error.request) {
            errorMessage = "Network error. Please check your connection";
        } else if (error.message) {
            errorMessage = error.message;
        }

        formState.errors.general = errorMessage;
        showToast(errorMessage, "error");
    } finally {
        formState.isLoading = false;
    }

    return false;
};
</script>

<template>
    <div class="bg-white rounded-lg shadow-md p-8">
        <div class="text-center">
            <h2 class="text-2xl font-bold text-gray-900">Login</h2>
            <p class="mt-2 text-gray-600">Don't have an account?
                <router-link to="/register" class="text-indigo-600 hover:text-indigo-500 font-medium">
                    Register Here
                </router-link>
            </p>
        </div>

        <!-- error message -->
        <div v-if="formState.errors.general" class="mt-4 p-3 bg-red-100 text-red-700 rounded-md">
            {{ formState.errors.general }}
        </div>

        <div class="mt-6">
            <div class="mt-6">
                <GoogleAuthButton :is-loading="formState.isLoading" />
            </div>

            <div class="relative mt-6">
                <div class="relative flex justify-center text-sm z-50">
                    <span class="px-2 bg-white text-gray-500">Or continue with</span>
                </div>

                <div class="absolute inset-0 flex items-center">
                    <div class="w-full border-t border-gray-300"></div>
                </div>
            </div>
        </div>

        <!-- Gunakan .prevent di sini -->
        <form @submit.prevent.stop="handleLogin" class="mt-6 space-y-6" novalidate>
            <div>
                <label for="identifier" class="block text-sm font-medium text-gray-700">
                    Username or Email
                </label>
                <input 
                    id="identifier"
                    v-model="formState.identifier"
                    type="text"
                    :class="[
                        'mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500',
                        formState.errors.identifier ? 'border-red-300' : 'border-gray-300'
                    ]" 
                    :disabled="formState.isLoading"
                    placeholder="Enter your valid username or email"
                />
                <p v-if="formState.errors.identifier" class="mt-1 text-sm text-red-600">
                    {{ formState.errors.identifier }}
                </p>
            </div>

            <div>
                <label for="password" class="block text-sm font-medium text-gray-700">
                    Password
                </label>
                <div class="relative mt-1">
                    <input 
                        id="password"
                        v-model="formState.password"
                        :type="showPassword ? 'text' : 'password'"
                        :class="[
                            'block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500',
                            formState.errors.password ? 'border-red-500' : 'border-gray-300'
                        ]"
                        :disabled="formState.isLoading"
                        placeholder="********"
                    />
                    <button
                        type="button"
                        @click="togglePasswordVisibility"
                        class="absolute inset-y-0 right-0 flex items-center pr-3"
                        :disabled="formState.isLoading"
                    >
                        <EyeIcon v-if="!showPassword" class="h-7 w-7 text-gray-400 pr-2" />
                        <EyeOffIcon v-else class="h-7 w-7 text-gray-400 pr-2" />
                    </button>
                </div>
                <p v-if="formState.errors.password" class="mt-1 text-sm text-red-600">
                    {{ formState.errors.password }}
                </p>
            </div>

            <button
                type="submit"
                :disabled="formState.isLoading"
                :class="[
                    'w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm font-medium text-white',
                    formState.isLoading
                        ? 'bg-gray-400 cursor-not-allowed'
                        : 'bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                ]"
            >
                <span v-if="formState.isLoading">Loading...</span>
                <span v-else>Login</span>
            </button>
        </form>
    </div>
</template>