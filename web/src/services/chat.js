import apiClient from '../utils/api'

export const chatService = {
  async createChatStream(message) {
    const response = await apiClient.post(
      '/chat/stream',
      { message },
      {
        responseType: 'text',
        headers: {
          Accept: 'text/event-stream',
          'Cache-Control': 'no-cache',
          Connection: 'keep-alive'
        }
      }
    )
    return response
  }
}
