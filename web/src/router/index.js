import { createRouter, createWebHistory } from 'vue-router';
import Login from "@/views/UserLogin.vue";
import Register from "@/views/UserRegister.vue";
import DashboardUser from "@/views/DashboardUser.vue";
import DashboardAdmin from "@/views/DashboardAdmin.vue";
import PageNotFound from "@/views/PageNotFound.vue";

const routes = [
    {
        path: "/",
        redirect: "/login"
    },
    {
        path: "/:pathMatch(.*)*",
        name: "NotFound",
        component: PageNotFound
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