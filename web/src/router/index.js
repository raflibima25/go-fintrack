import { createRouter, createWebHistory } from 'vue-router';
import Login from "@/components/UserLogin.vue";
import Register from "@/components/UserRegister.vue";
import DashboardUser from "@/components/DashboardUser.vue";
import DashboardAdmin from "@/components/DashboardAdmin.vue";

const routes = [
    {
        path: "/",
        redirect: "/login"
    },
    {
        path: "/login",
        name: "UserLogin",
        component: Login,
        meta: { requiresAuth: false }
    },
    {
        path: "/register",
        name: "UserRegister",
        component: Register,
        meta: { requiresAuth: false }
    },
    {
        path: "/dashboard",
        name: "DashboardUser",
        component: DashboardUser,
        meta: { requiresAuth: true, role: 'user' }
    },
    {
        path: "/admin-dashboard",
        name: "DashboardAdmin",
        component: DashboardAdmin,
        meta: { requiresAuth: true, role: 'admin' }
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

// navigation guard
router.beforeEach((to, from, next) => {
    const token = localStorage.getItem("token");
    const isAdmin = localStorage.getItem("isAdmin") === 'true';
    console.log("Navigation guard - isAdmin:", isAdmin) // debug

    if (!token && to.meta.requiresAuth) {
        next('/login');
        return;
    }

    if (token) {
        // jika mencoba akses login/register saat sudah login
        if (to.path === '/login' || to.path === '/register') {
            next(isAdmin ? '/admin-dashboard' : '/dashboard');
            return;
        }

        // jika user mencoba akses halaman admin
        if (to.meta.role === 'admin' && !isAdmin) {
            next('/dashboard');
            return;
        }

        // jika admin mencoba akses halaman user
        if (to.meta.role === 'user' && isAdmin) {
            next('/admin-dashboard');
            return;
        }
    }

    next()
});

export default router;