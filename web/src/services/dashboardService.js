import apiClient from '../utils/api'

export const dashboardService = {
  async getFinancialOverview() {
    return apiClient.get('/dashboard/overview')
  },

  async getExpenseAnalysis() {
    return apiClient.get('/dashboard/charts')
  }
}
