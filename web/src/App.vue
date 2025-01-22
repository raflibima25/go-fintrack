<template>
  <div id="app">
    <router-view></router-view>
    <ToastContainer />
  </div>
</template>

<script setup>
import ToastContainer from "@/components/ToastContainer.vue";
</script>