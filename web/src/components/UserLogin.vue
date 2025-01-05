<template>
  <div class="login-container">
    <h1>Login</h1>
    <form @submit.prevent="login">
      <input v-model="identifier" placeholder="Email or Username" required />
      <input v-model="password" type="password" placeholder="Password" required />
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script>
import apiClient from "@/utils/api";

export default {
  data() {
    return {
      identifier: "",
      password: "",
    };
  },
  methods: {
    async login() {
      try {
        const response = await apiClient.post("/user/login", {
          email_or_username: this.identifier, // email_or_username dan password merupakan field json backend
          password: this.password,
        });
        if (response.data.status) {
          localStorage.setItem("token", response.data.data.access_token);
          this.$router.push("/dashboard");
        } else {
          alert(response.data.message || "Login gagal");
        }
      } catch (error) {
        console.error("Login error:", error.response?.data);
        alert(error.response?.data?.message || "Login gagal ");
      }
    },
  },
};
</script>