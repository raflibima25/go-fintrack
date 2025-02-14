import { createRouter, createWebHistory } from 'vue-router'

// Gunakan lazy loading untuk semua komponen
const Login = () => import('@/views/LoginAuth.vue')
const Register = () => import('@/views/RegisterAuth.vue')
const DashboardUser = () => import('@/views/DashboardUser.vue')
const DashboardAdmin = () => import('@/views/DashboardAdmin.vue')
const PageNotFound = () => import('@/views/PageNotFound.vue')
const TransactionList = () => import('@/views/TransactionList.vue')
const CategoryList = () => import('@/views/CategoryList.vue')
const ChatAssistant = () => import('@/views/ChatAssistant.vue')
const LandingPage = () => import('@/views/LandingPage.vue')

const routes = [
  // public routes
  {
    path: '/',
    name: 'LandingPage',
    component: LandingPage,
    meta: { requiresAuth: false }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: PageNotFound,
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'LoginAuth',
    component: Login,
    meta: {
      requiresAuth: false,
      title: 'Login'
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: {
      requiresAuth: false,
      title: 'Register'
    }
  },
  {
    path: '/auth/google/callback',
    name: 'GoogleCallback',
    component: () => import('@/views/GoogleCallback.vue'),
    meta: {
      requiresAuth: false
    }
  },
  // user routes
  {
    path: '/dashboard',
    name: 'DashboardUser',
    component: DashboardUser,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Dashboard',
      layout: 'UserLayout'
    }
  },
  {
    path: '/transactions',
    name: 'TransactionsUser',
    component: TransactionList,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Transactions',
      layout: 'UserLayout'
    }
  },
  {
    path: '/categories',
    name: 'CategoriesUser',
    component: CategoryList,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Categories',
      layout: 'UserLayout'
    }
  },
  {
    path: '/chat-assistant',
    name: 'ChatAssistant',
    component: ChatAssistant,
    meta: {
      requiresAuth: true,
      role: 'user',
      title: 'Chat Assistant',
      layout: 'UserLayout'
    }
  },
  // admin routes
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
  document.title = `${to.meta.title || 'Financial Tracker'}`

  // Jika route memerlukan auth dan user tidak terautentikasi
  if (to.meta.requiresAuth && !isAuthenticated()) {
    next({
      name: 'LoginAuth',
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
