import apiClient from "@/utils/api";

export const categoryService = {
    async getCategories() {
        return apiClient.get('/category')
    },

    async getCategoryById(id) {
        return apiClient.get(`/category/${id}`)
    },

    async createCategory(data) {
        return apiClient.post(`/category`, data)
    },

    async updateCategory(id, data) {
        return apiClient.put(`/category/${id}`, data)
    },

    async deleteCategory(id) {
        return apiClient.delete(`/category/${id}`)
    }
}