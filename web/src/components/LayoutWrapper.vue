<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import UserLayout from '../layouts/UserLayout.vue';

const route = useRoute();

const defaultLayout = { 
    template: '<div><slot></slot></div>',
}

const layout = computed(() => {
  const layoutName = route.meta.layout;
  
  switch(layoutName) {
    case 'UserLayout':
      return UserLayout;
    default:
      return defaultLayout;
  }
});
</script>

<template>
    <component :is="layout">
      <slot></slot>
    </component>
</template>