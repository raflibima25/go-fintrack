<script>
import BaseAuth from '../components/BaseAuth.vue';
import RegisterCard from '../components/auth/RegisterCard.vue';

export default {
  name: 'RegisterAuth',
  components: {
    BaseAuth,
    RegisterCard
  }
}
</script>

<template>
  <BaseAuth illustration-src="/vector/auth.svg">
    <RegisterCard />
  </BaseAuth>
</template>