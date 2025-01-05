<template>
  <div class="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-sm">
      <img class="mx-auto h-10 w-auto" src="https://tailwindui.com/plus/img/logos/mark.svg?color=indigo&shade=600" alt="Your Company" />
      <h2 class="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">Register new account</h2>
    </div>

    <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
      <form class="space-y-6" @submit.prevent="register">
        <div>
          <label for="name" class="block text-sm/6 font-medium text-gray-900">Name</label>
          <div class="mt-2">
            <input
                v-model="name"
                id="name"
                name="name"
                type="text"
                required
                class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            />
          </div>
        </div>

        <div>
          <label for="email" class="block text-sm/6 font-medium text-gray-900">Email address</label>
          <div class="mt-2">
            <input
                v-model="email"
                id="email"
                name="email"
                type="email"
                autocomplete="email"
                required
                class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            />
          </div>
        </div>

        <div>
          <label for="username" class="block text-sm/6 font-medium text-gray-900">Username</label>
          <div class="mt-2">
            <input
                v-model="username"
                id="username"
                name="username"
                type="text"
                required
                class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            />
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm/6 font-medium text-gray-900">Password</label>
          <div class="mt-2">
            <input
                v-model="password"
                id="password"
                name="password"
                type="password"
                required
                class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            />
          </div>
        </div>

        <div>
          <label for="confirm-password" class="block text-sm/6 font-medium text-gray-900">Confirm Password</label>
          <div class="mt-2">
            <input
                v-model="confirmPassword"
                id="confirm-password"
                name="confirm-password"
                type="password"
                required
                class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            />
          </div>
        </div>

        <div>
          <button
              type="submit"
              class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Register
          </button>
        </div>
      </form>

      <p class="mt-10 text-center text-sm/6 text-gray-500">
        Sudah punya akun?
        {{ ' ' }}
        <a href="/login" class="font-semibold text-indigo-600 hover:text-indigo-500">Login</a>
      </p>
    </div>
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
        if (response.data.status) {
          this.$router.push("/login");
        }
      } catch (error) {
        alert(error.response.data.message || "Registration failed");
      }
    },
  },
};
</script>