import { createRouter, createWebHistory } from 'vue-router'

// Gunakan lazy loading untuk semua komponen
const Login = () => import('@/views/UserLogin.vue')
const Register = () => import('@/views/UserRegister.vue')
const DashboardUser = () => import('@/views/DashboardUser.vue')
const DashboardAdmin = () => import('@/views/DashboardAdmin.vue')
const PageNotFound = () => import('@/views/PageNotFound.vue')
const TransactionList = () => import('@/views/TransactionList.vue')

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: PageNotFound,
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'UserLogin',
    component: Login,
    meta: {
      requiresAuth: false,
      title: 'Login'
    }
  },
  {
    path: '/register',
    name: 'UserRegister',
    component: Register,
    meta: {
      requiresAuth: false,
      title: 'Register'
    }
  },
  {
    path: '/dashboard',
    name: 'DashboardUser',
    component: DashboardUser,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Dashboard'
    }
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: TransactionList,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Transactions'
    }
  },
  {
    path: '/admin-dashboard',
    name: 'DashboardAdmin',
    component: DashboardAdmin,
    meta: {
      requiresAuth: true,
      role: 'admin',
      title: 'Admin Dashboard'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Fungsi helper untuk validasi
const isAuthenticated = () => !!localStorage.getItem('token')
const isAdmin = () => localStorage.getItem('isAdmin') === 'true'
const getUserRole = () => (isAdmin() ? 'admin' : 'user')

// Navigation guard yang dioptimalkan
router.beforeEach(async (to, from, next) => {
  // Update title
  document.title = `${to.meta.title || 'Manajemen Keuangan'}`

  // Jika route memerlukan auth dan user tidak terautentikasi
  if (to.meta.requiresAuth && !isAuthenticated()) {
    next({
      name: 'UserLogin',
      query: { redirect: to.fullPath }
    })
    return
  }

  // Jika user sudah login
  if (isAuthenticated()) {
    const userRole = getUserRole()

    // Mencegah akses ke login/register
    if (to.path === '/login' || to.path === '/register') {
      next(isAdmin() ? '/admin-dashboard' : '/dashboard')
      return
    }

    // Validasi role-based access
    if (to.meta.role && to.meta.role !== userRole) {
      next(userRole === 'admin' ? '/admin-dashboard' : '/dashboard')
      return
    }
  }

  next()
})

// Error handling untuk route
router.onError((error) => {
  console.error('Router error:', error)
  router.push({ name: 'NotFound' })
})

export default router
