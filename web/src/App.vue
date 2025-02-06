<template>
  <div id="app" class="font-poppins ">
    <NavbarUser v-if="shouldShowNavbar"/>
    <router-view></router-view>
    <ToastContainer />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router';
import { useAuth } from './composables/useAuth';
import ToastContainer from "@/components/ToastContainer.vue";
import NavbarUser from "./components/NavbarUser.vue";

const route = useRoute();
const { isAuthenticated, isAdmin } = useAuth();

const shouldShowNavbar = computed(() => {
  // Daftar route yang tidak perlu navbar
  const noNavbarRoutes = ['/login', '/register'];
  
  return isAuthenticated.value && 
         !isAdmin.value && 
         !noNavbarRoutes.includes(route.path);
});
</script>