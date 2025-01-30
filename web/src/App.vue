<template>
  <div id="app">
    <MainNavbar />
    <router-view></router-view>
    <ToastContainer />
  </div>
</template>

<script setup>
import ToastContainer from "@/components/ToastContainer.vue";
import MainNavbar from "./components/MainNavbar.vue";
</script>