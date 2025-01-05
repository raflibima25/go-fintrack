import { createRouter, createWebHistory } from 'vue-router';
import Login from "@/components/UserLogin.vue";
import Register from "@/components/UserRegister.vue";
import DashboardUser from "@/components/DashboardUser.vue";

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
        meta: { requiresAuth: true }
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

// navigation guard
router.beforeEach((to, from, next) => {
    const token = localStorage.getItem("token");

    if (to.meta.requiresAuth && !token) {
        next('/login');
    } else if (token && (to.path === '/login' || to.path === '/register')) {
        next('/dashboard');
    } else {
        next();
    }
});

export default router;