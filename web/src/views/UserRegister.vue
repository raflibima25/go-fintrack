<template>
  <div class="register-container">
    <h1>Register</h1>
    <form @submit.prevent="register">
      <input v-model="name" placeholder="Name" required />
      <input v-model="email" placeholder="Email" type="email" required />
      <input v-model="username" placeholder="Username" required />
      <input v-model="password" placeholder="Password" type="password" required />
      <input v-model="confirmPassword" placeholder="Confirm Password" type="password" required />
      <button type="submit">Register</button>
    </form>
  </div>
</template>

<script>
import apiClient from "@/utils/api";

export default {
  data() {
    return {
      name: "",
      email: "",
      username: "",
      password: "",
      confirmPassword: "",
    };
  },
  methods: {
    async register() {
      try {
        const response = await apiClient.post("/user/register", {
          name: this.name,
          email: this.email,
          username: this.username,
          password: this.password,
          confirm_password: this.confirmPassword,
        });
        alert(response.data.message);
      } catch (error) {
        alert(error.response.data.message || "Registration failed");
      }
    },
  },
};
</script>