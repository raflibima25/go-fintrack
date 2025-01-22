import apiClient from "@/utils/api";

export const transactionService = {
    async getTransactions(params = {}) {
        return apiClient.get('/transaction', { params })
    },

    async createTransaction(data) {
        return apiClient.post('/transaction', data)
    },

    async updateTransaction(id, data) {
        return apiClient.put(`/transaction/${id}`, data)
    },

    async deleteTransaction(id) {
        return apiClient.delete(`/transaction/${id}`)
    },

    async exportTransactions(params = {}) {
        return apiClient.get('/transaction/export', {
            params,
            responseType: 'blob'
        })
    }
}