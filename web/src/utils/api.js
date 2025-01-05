import axios from "axios";

const apiClient = axios.create({
    // baseURL: (process.env.VUE_APP_API_BASE_URL || "http://localhost:8080") + "/api",
    baseURL: '/api',
    headers: {
        "Content-Type": "application/json",
    },
});

apiClient.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export default apiClient;