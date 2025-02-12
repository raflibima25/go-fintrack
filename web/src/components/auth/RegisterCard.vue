<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from '@/composables/useToast';
import apiClient from "@/utils/api";
import { EyeIcon, EyeOffIcon } from "lucide-vue-next";
import GoogleAuthButton from '@/components/auth/GoogleAuthButton.vue';

const router = useRouter();
const { showToast } = useToast();

const formState = reactive({
    name: '',
    email: '',
    username: '',
    password: '',
    confirmPassword: '',
    isLoading: false,
    errors: {
        name: '',
        email: '',
        username: '',
        password: '',
        confirmPassword: '',
        general: ''
    }
});

const showPassword = ref(false);
const showConfirmPassword = ref(false);

const validateForm = () => {
    let isValid = true;
    formState.errors = {
        name: '',
        email: '',
        username: '',
        password: '',
        confirmPassword: '',
        general: ''
    };

    if (!formState.name.trim()) {
        formState.errors.name = 'Name is required';
        isValid = false;
    }

    if (!formState.email.trim()) {
        formState.errors.email = 'Email is required';
        isValid = false;
    }

    if (!formState.username.trim()) {
        formState.errors.username = 'Username is required';
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

    if (formState.password !== formState.confirmPassword) {
        formState.errors.confirmPassword = 'Passwords do not match';
        isValid = false;
    }

    return isValid;
};

const register = async () => {
    if (!validateForm()) return;

    formState.isLoading = true;
    formState.errors.general = '';

    try {
        const response = await apiClient.post("/auth/register", {
            name: formState.name,
            email: formState.email,
            username: formState.username,
            password: formState.password,
            confirm_password: formState.confirmPassword
        });

        if (response.data.status) {
            showToast(response.data.message, "success");
            router.push("/login");
        } else {
            throw new Error(response.data.message || "Registration failed");
        }
    } catch (error) {
        console.error("Register error:", error);
        formState.errors.general = error.response?.data?.message || "Registration failed";
        showToast(formState.errors.general, "error");
    } finally {
        formState.isLoading = false;
    }
};
</script>

<template>
    <div class="bg-white rounded-lg shadow-md p-8">
        <div class="text-center">
            <h2 class="text-2xl font-bold text-gray-900">Create Account</h2>
            <p class="mt-2 text-gray-600">Have an account?
                <router-link to="/login" class="text-indigo-600 hover:text-indigo-500 font-medium">
                    Login Here
                </router-link>
            </p>
        </div>

        <div class="mt-6">
            <GoogleAuthButton :is-loading="formState.isLoading" text="Sign up with Google" />

            <div class="relative mt-6">
                <div class="relative flex justify-center text-sm z-50">
                    <span class="px-2 bg-white text-gray-500">Or continue with</span>
                </div>

                <div class="absolute inset-0 flex items-center">
                    <div class="w-full border-t border-gray-300"></div>
                </div>
            </div>
        </div>

        <!-- Error message -->
        <div v-if="formState.errors.general" class="mt-4 p-3 bg-red-100 text-red-700 rounded-md">
            {{ formState.errors.general }}
        </div>

        <form @submit.prevent="register" class="mt-6 space-y-4" novalidate>
            <div>
                <label for="name" class="block text-sm font-medium text-gray-700">Fullname</label>
                <input 
                    id="name"
                    v-model="formState.name"
                    type="text"
                    :disabled="formState.isLoading"
                    :class="[
                        'mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm',
                        formState.errors.name ? 'border-red-300' : 'border-gray-300'
                    ]"
                    placeholder="Enter your fullname"
                >
                <p v-if="formState.errors.name" class="mt-1 text-sm text-red-600">
                    {{ formState.errors.name }}
                </p>
            </div>

            <div>
                <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
                <input 
                    id="username"
                    v-model="formState.username"
                    type="text"
                    :disabled="formState.isLoading"
                    :class="[
                        'mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm',
                        formState.errors.username ? 'border-red-300' : 'border-gray-300'
                    ]"
                    placeholder="Enter your username"
                />
                <p v-if="formState.errors.username" class="mt-1 text-sm text-red-600">
                    {{ formState.errors.username }}
                </p>
            </div>

            <div>
                <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                <input 
                    id="email"
                    v-model="formState.email"
                    type="email"
                    :disabled="formState.isLoading"
                    :class="[
                        'mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm',
                        formState.errors.email ? 'border-red-300' : 'border-gray-300'
                    ]"
                    placeholder="member@fintrack.com"
                />
                <p v-if="formState.errors.email" class="mt-1 text-sm text-red-600">
                    {{ formState.errors.email }}
                </p>
            </div>

            <div>
                <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
                <div class="relative mt-1">
                    <input 
                        id="password"
                        v-model="formState.password"
                        :type="showPassword ? 'text' : 'password'"
                        :disabled="formState.isLoading"
                        :class="[
                            'block w-full px-3 py-2 border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500',
                            formState.errors.password ? 'border-red-300' : 'border-gray-300'
                        ]"
                        placeholder="********"
                    >
                    <button 
                        type="button"    
                        @click="showPassword = !showPassword"
                        :disabled="formState.isLoading"
                        class="absolute inset-y-0 right-0 pr-3 flex items-center"
                    >
                        <EyeIcon v-if="!showPassword" class="h-7 w-7 text-gray-400 pr-2" />
                        <EyeOffIcon v-else class="h-7 w-7 text-gray-400 pr-2" />
                    </button>
                </div>
                <p v-if="formState.errors.password" class="mt-1 text-sm text-red-600">
                    {{ formState.errors.password }}
                </p>
            </div>

            <div>
                <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirmation Password</label>
                <div class="relative mt-1">
                    <input 
                        id="confirmPassword"
                        v-model="formState.confirmPassword"
                        :type="showConfirmPassword ? 'text' : 'password'"
                        :disabled="formState.isLoading"
                        :class="[
                            'block w-full px-3 py-2 border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500',
                            formState.errors.confirmPassword ? 'border-red-300' : 'border-gray-300'
                        ]"
                        placeholder="********"
                    >
                    <button 
                        type="button"    
                        @click="showConfirmPassword = !showConfirmPassword"
                        :disabled="formState.isLoading"
                        class="absolute inset-y-0 right-0 pr-3 flex items-center"
                    >
                        <EyeIcon v-if="!showConfirmPassword" class="h-7 w-7 text-gray-400 pr-2" />
                        <EyeOffIcon v-else class="h-7 w-7 text-gray-400 pr-2" />
                    </button>
                </div>
                <p v-if="formState.errors.confirmPassword" class="mt-1 text-sm text-red-600">
                    {{ formState.errors.confirmPassword }}
                </p>
            </div>

            <button
                type="submit"
                :disabled="formState.isLoading"
                class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
                <span v-if="formState.isLoading">Registering...</span>
                <span v-else>Register</span>
            </button>
        </form>
    </div>
</template>