<script setup>
import { onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useGoogleAuth } from '../composables/useGoogleAuth';

const route = useRoute();
const router = useRouter();
const { loading, handleCallback } = useGoogleAuth(); 

onMounted(async () => {
    const { code, state } = route.query;
    console.log("Callback mounted with:", { code,state }) // debug

    if (!code || !state) {
        console.log("Missing parameters");
        router.push("/login");
        return;
    }

    try {
        await handleCallback(code, state);
    } catch (error) {
        console.log("Callback error:", error)
        router.push("/login");
    }
})
</script>

<template>
    <div class="min-h-screen flex items-center justify-center">
        <div v-if="loading" class="text-center">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
            <p class="mt-4 text-gray-600">Completing authentication...</p>
        </div>
    </div>
</template>