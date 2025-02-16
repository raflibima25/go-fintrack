export const useIndexedDB = () => {
  const DB_NAME = 'chatAssistantDB'
  const DB_VERSION = 1
  const CHAT_STORE = 'chats'

  const initDB = () => {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION)

      request.onerror = () => {
        reject(request.error)
      }

      request.onsuccess = () => {
        resolve(request.result)
      }

      request.onupgradeneeded = (event) => {
        const db = event.target.result
        if (!db.objectStoreNames.contains(CHAT_STORE)) {
          const store = db.createObjectStore(CHAT_STORE, { keyPath: 'id' })
          store.createIndex('createdAt', 'createdAt', { unique: false })
        }
      }
    })
  }

  const getDB = async () => {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION)
      request.onerror = () => reject(request.error)
      request.onsuccess = () => resolve(request.result)
    })
  }

  const saveChat = async (chat) => {
    const db = await getDB()
    return new Promise((resolve, reject) => {
      try {
        // Buat objek baru yang hanya berisi properti yang diperlukan
        const chatToSave = {
          id: chat.id,
          title: chat.title,
          createdAt: chat.createdAt,
          lastMessage: chat.lastMessage,
          messages: chat.messages.map((msg) => ({
            chatId: msg.chatId,
            role: msg.role,
            content: msg.content,
            timestamp: msg.timestamp
          }))
        }

        const transaction = db.transaction(CHAT_STORE, 'readwrite')
        const store = transaction.objectStore(CHAT_STORE)
        const request = store.put(chatToSave)

        request.onerror = () => reject(request.error)
        request.onsuccess = () => resolve(request.result)
      } catch (error) {
        reject(error)
      }
    })
  }

  const getChats = async () => {
    const db = await getDB()
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(CHAT_STORE, 'readonly')
      const store = transaction.objectStore(CHAT_STORE)
      const request = store.getAll()

      request.onerror = () => reject(request.error)
      request.onsuccess = () => {
        // Sort chats by creation date, newest first
        const chats = request.result.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
        resolve(chats)
      }
    })
  }

  const getChatHistory = async (chatId) => {
    const db = await getDB()
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(CHAT_STORE, 'readonly')
      const store = transaction.objectStore(CHAT_STORE)
      const request = store.get(chatId)

      request.onerror = () => reject(request.error)
      request.onsuccess = () => {
        resolve(request.result?.messages || [])
      }
    })
  }

  const deleteChat = async (chatId) => {
    const db = await getDB()
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(CHAT_STORE, 'readwrite')
      const store = transaction.objectStore(CHAT_STORE)
      const request = store.delete(chatId)

      request.onerror = () => reject(request.error)
      request.onsuccess = () => resolve()
    })
  }

  return {
    initDB,
    saveChat,
    getChats,
    getChatHistory,
    deleteChat
  }
}
