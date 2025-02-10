<script setup>
import NavbarUser from '../components/NavbarUser.vue';
</script>

<template>
    <div class="min-h-screen bg-gray-50">
        <NavbarUser />
        <main>
            <slot></slot>
        </main>
    </div>
</template>