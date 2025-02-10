<template>
  <div id="app" class="font-poppins ">
    <router-view></router-view>
    <ToastContainer />
  </div>
</template>

<script setup>
import ToastContainer from "@/components/ToastContainer.vue";
</script>