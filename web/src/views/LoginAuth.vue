<script>
import BaseAuth from '../components/BaseAuth.vue';
import LoginCard from '../components/auth/LoginCard.vue';

export default {
  name: 'LoginAuth',
  components: {
    BaseAuth,
    LoginCard
  }
}
</script>

<template>
  <BaseAuth illustration-src="/vector/auth.svg">
    <LoginCard />  
  </BaseAuth>
</template>