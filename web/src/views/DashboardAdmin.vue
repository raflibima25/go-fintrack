<template>
  <div class="min-h-screen bg-gray-100 p-8">
    <div class="max-w-4xl mx-auto bg-white rounded-lg shadow-md p-6">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold text-gray-800">Welcome to the Dashboard Admin</h1>
        <button
            @click="logout"
            class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
        >
          Logout
        </button>
      </div>

      <!-- Tambahkan konten dashboard di sini -->
      <div class="mt-4">
        <p class="text-gray-600">Selamat datang di dashboard aplikasi manajemen keuangan.</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DashboardAdmin',

  created() {
    this.checkAuth();
  },

  methods: {
    checkAuth() {
      const token = localStorage.getItem('token');
      const isAdmin = localStorage.getItem('isAdmin') === 'true';

      if (!token) {
        this.$router.push('/login');
        return;
      }

      if (!isAdmin) {
        this.$router.push('/dashboard');
        return;
      }
    },

    async logout() {
      try {
        localStorage.removeItem("token");
        localStorage.removeItem("isAdmin")
        this.$router.push("/login");
      } catch (error) {
        console.log('Logout error:', error)
      }
    },
  },
};
</script>

<style scoped>
.dashboard-user {
  padding: 20px;
}

h1 {
  margin-bottom: 20px;
}

button {
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}
</style>