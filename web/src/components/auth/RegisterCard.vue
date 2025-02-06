<script>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import apiClient from "@/utils/api";
import { EyeIcon, EyeOffIcon } from "lucide-vue-next";

export default {
    name: 'RegisterCard',
    components: {
        EyeIcon,
        EyeOffIcon
    },
    setup() {
        const router = useRouter();
        const form = reactive({
            name: '',
            email: '',
            username: '',
            password: '',
            confirmPassword: ''
        })

        const showPassword = ref(false)
        const showConfirmPassword = ref(false)

        const register = async () => {
            try {
                const response = await apiClient.post("/auth/register", {
                    name: form.name,
                    email: form.email,
                    username: form.username,
                    password: form.password,
                    confirm_password: form.confirmPassword
                })

                if (response.data.status) {
                    alert(response.data.message)
                    router.push("/login")
                }
            } catch (error) {
                console.log("Register error:", error)
                alert(error.response?.data?.message || "Register failed")
            }
        }

        return {
            form,
            showPassword,
            showConfirmPassword,
            register
        }
    }
}
</script>

<template>
    <div class="bg-white rounded-lg shadow-md p-8">
        <div class="text-center">
            <h2 class="text-2xl font-bold text-gray-900">Create Account</h2>
            <p class="mt-2 text-gray-600">Have an account?
                <router-link to="/login" class="text-blue-500 hover:text-blue-600 font-medium">
                    Login Here
                </router-link>
            </p>
        </div>

        <form @submit.prevent="register" class="mt-8 space-y-4">
            <div>
                <label for="name" class="block text-sm font-medium text-gray-700">Fullname</label>
                <input 
                    v-model="form.name"
                    type="text"
                    required
                    class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    placeholder="Enter your fullname"
                >
            </div>

            <div>
                <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
                <input 
                    v-model="form.username"
                    type="text"
                    required
                    class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    placeholder="Enter your username"
                />
            </div>

            <div>
                <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                <input 
                    v-model="form.email"
                    type="email"
                    required
                    class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    placeholder="member@fintrack.com"
                />
            </div>

            <div>
                <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
                <div class="relative mt-1">
                    <input 
                        v-model="form.password"
                        :type="showPassword ? 'text' : 'password'"
                        name="password"
                        required
                        class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        placeholder="********"
                    >
                    <button 
                        type="button"    
                        @click="showPassword = !showPassword"
                        class="absolute inset-y-0 right-0 pr-3 flex items-center"
                    >
                        <EyeIcon v-if="!showPassword" class="h-7 w-7 text-gray-400 pr-2" />
                        <EyeOffIcon v-else class="h-7 w-7 text-gray-400 pr-2" />
                    </button>
                </div>
            </div>

            <div>
                <label for="password" class="block text-sm font-medium text-gray-700">Confirmation Password</label>
                <div class="relative mt-1">
                    <input 
                        v-model="form.confirmPassword"
                        :type="showConfirmPassword ? 'text' : 'password'"
                        name="password"
                        required
                        class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        placeholder="********"
                    >
                    <button 
                        type="button"    
                        @click="showConfirmPassword  = !showConfirmPassword "
                        class="absolute inset-y-0 right-0 pr-3 flex items-center"
                    >
                        <EyeIcon v-if="!showConfirmPassword" class="h-7 w-7 text-gray-400 pr-2" />
                        <EyeOffIcon v-else class="h-7 w-7 text-gray-400 pr-2" />
                    </button>
                </div>
            </div>

            <button
                type="submit"
                class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-500 hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
                Register
            </button>
        </form>
    </div>
</template>