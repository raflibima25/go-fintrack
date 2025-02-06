<script>
import { ref } from "vue";
import { useRouter } from 'vue-router';
import apiClient from "@/utils/api";
import { EyeIcon, EyeOffIcon } from "lucide-vue-next";

export default {
    name: 'LoginCard',
    components: {
        EyeIcon,
        EyeOffIcon
    },
    setup() {
        const router = useRouter();
        const identifier = ref('');
        const password = ref('');
        const showPassword = ref(false);

        const login = async () => {
            try {
                const response = await apiClient.post("/auth/login", {
                    email_or_username: identifier.value, // email_or_username dan password merupakan field json backend
                    password: password.value
                });

                if (response.data.status) {
                    localStorage.setItem("token", response.data.data.access_token);
                    localStorage.setItem("isAdmin", response.data.data.is_admin);

                    // redirect berdasarkan role
                    if (response.data.data.is_admin) {
                        router.push("/admin-dashboard");
                    } else {
                        router.push("/dashboard");
                    }
                } 
            } catch (error) {
                console.error("Login error:", error)
                alert(error.response?.data?.message || "Login failed")
            }
        }

        return {
            identifier,
            password,
            showPassword,
            login
        }
    }
};
</script>

<template>
    <div class="bg-white rounded-lg shadow-md p-8">
        <div class="text-center">
            <h2 class="text-2xl font-bold text-gray-900">Login</h2>
            <p class="mt-2 text-gray-600">Don't have an account?
                <router-link to="/register" class="text-blue-500 hover:text-blue-600 font-medium">
                    Register Here
                </router-link>
            </p>
        </div>

        <form @submit.prevent="login" class="mt-8 space-y-6">
            <div>
                <label for="username_or_email" class="block text-sm font-medium text-gray-700">
                    Username or Email
                </label>
                <input 
                    id="username_or_email"
                    v-model="identifier"
                    type="text"
                    required
                    class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500" 
                    placeholder="Enter your valid username or email"
                />
            </div>

            <div>
                <label for="password" class="block text-sm font-medium text-gray-700">
                    Password
                </label>
                <div class="relative mt-1">
                    <input 
                        id="password"
                        v-model="password"
                        :type="showPassword ? 'text' : 'password'"
                        class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                        placeholder="********"
                    />
                    <button
                        type="button"
                        @click="showPassword = !showPassword"
                        class="absolute inset-y-0 right-0 flex items-center pr-3"
                    >
                        <EyeIcon v-if="!showPassword" class="h-7 w-7 text-gray-400 pr-2" />
                        <EyeOffIcon v-else class="h-7 w-7 text-gray-400 pr-2" />
                    </button>
                </div>
            </div>

            <button
                type="submit"
                class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm font-medium text-white bg-blue-500 hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-pink-500"
            >
                Login
            </button>
        </form>
    </div>
</template>