<script setup>
import { ref, onMounted, nextTick, watch } from 'vue'
import { marked } from 'marked'
import { PlusCircle, Settings, Search, Trash2, Send, Copy, Check, Edit2, X, User, Bot, Menu } from 'lucide-vue-next'
import { useIndexedDB } from '../composables/useIndexedDB'
import { useClipboard } from '@vueuse/core'
import { useToast } from '@/composables/useToast'

const {
  initDB,
  saveChat,
  getChats,
  getChatHistory,
  deleteChat
} = useIndexedDB()

const messages = ref([])
const chats = ref([])
const newMessage = ref('')
const isStreaming = ref(false)
const currentStreamingMessage = ref('')
const selectedChatId = ref(null)
const searchTerm = ref('')
let controller = null
const messagesContainer = ref(null)
const { showToast } = useToast()
const { copy, copied } = useClipboard()
const editingChatId = ref(null)
const editedTitle = ref('')
const isSidebarOpen = ref(false);

// Initialize IndexedDB and load chats
onMounted(async () => {
  await initDB()
  await loadChats()
  if (chats.value.length === 0) {
    await createNewChat()
  } else {
    selectedChatId.value = chats.value[0].id
    await loadChatHistory(selectedChatId.value)
  }
})

const loadChats = async () => {
  chats.value = await getChats()
}

const loadChatHistory = async (chatId) => {
  messages.value = await getChatHistory(chatId)
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// Format message with Markdown
marked.setOptions({
  breaks: true,
  gfm: true
})

const formatMessage = (text) => {
  let formatted = text

  // Pisahkan string menjadi paragraf
  const paragraphs = formatted.split(/(?<=\.)\s+(?=[A-Z])/)
  
  // Format setiap paragraf
  formatted = paragraphs.map(paragraph => {
    // Cek apakah paragraf mengandung list
    if (/^\d+\.\s/.test(paragraph)) {
      // Jika ini list item, tambahkan line break sebelumnya
      return '<br>' + paragraph.trim()
    }
    return paragraph.trim()
  }).join('<br><br>')

  // Format Markdown
  formatted = formatted
    // Bold
    .replace(/\*\*(.*?)\*\*/g, (_, p1) => `<strong>${p1}</strong>`)
    // Italic
    .replace(/\*(.*?)\*/g, (_, p1) => `<em>${p1}</em>`)
    // Underscore
    .replace(/_(.*?)_/g, (_, p1) => `<u>${p1}</u>`)

  // Handle numbered lists
  formatted = formatted.replace(/(\d+\.\s.*?)(?=(?:\d+\.|$))/g, (match) => {
    return `<div class="list-item">${match}</div>`
  })

  // Handle bullet points
  formatted = formatted.replace(/^[•-]\s(.*?)$/gm, (_, p1) => {
    return `<div class="list-item">• ${p1}</div>`
  })

  return formatted
}

const parseStreamChunk = (chunk) => {
  const lines = chunk.split('\n')
  let message = ''

  for (const line of lines) {
    if (line.startsWith('data:')) {
      // Remove 'data:' prefix and trim whitespace
      let data = line.slice(5).trim()
      if (data) {
        message += data + ' '
      }
    }
  }

  return message
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || isStreaming.value) return

  // Buat objek pesan user
  const messageObj = {
    chatId: selectedChatId.value,
    role: 'user',
    content: newMessage.value,
    timestamp: new Date().toISOString()
  }

  try {
    // Push pesan user
    messages.value.push(messageObj)

    // Update chat title jika ini pesan pertama
    const currentChat = chats.value.find(chat => chat.id === selectedChatId.value)
    if (currentChat && currentChat.messages?.length === 0) {
      const newTitle = generateChatTitle(newMessage.value)
      await saveChat({
        id: selectedChatId.value,
        title: newTitle,
        createdAt: currentChat.createdAt,
        lastMessage: generateMessageSummary(newMessage.value),
        messages: [...messages.value]
      })
    } else {
      await saveChat({
        id: selectedChatId.value,
        title: currentChat.title,
        createdAt: currentChat.createdAt,
        lastMessage: generateMessageSummary(newMessage.value),
        messages: [...messages.value]
      })
    }

    // Tambahkan placeholder pesan assistant
    const assistantMsg = {
      chatId: selectedChatId.value,
      role: 'assistant',
      content: '',
      timestamp: new Date().toISOString()
    }
    messages.value.push(assistantMsg)

    await scrollToBottom()
    isStreaming.value = true
    currentStreamingMessage.value = ''

    // Abort controller sebelumnya jika ada
    if (controller) {
      controller.abort()
    }
    controller = new AbortController()

    // Kirim request ke API
    const response = await fetch(`${import.meta.env.VITE_APP_API_BASE_URL}/api/chat/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ message: newMessage.value }),
      signal: controller.signal
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    // Baca response stream
    const reader = response.body.getReader()
    const decoder = new TextDecoder()

    while (true) {
      const { value, done } = await reader.read()
      if (done) break

      const chunk = decoder.decode(value)
      const parsedMessage = parseStreamChunk(chunk)
      
      if (parsedMessage) {
        currentStreamingMessage.value += parsedMessage
        messages.value[messages.value.length - 1].content = currentStreamingMessage.value.trim()
        await scrollToBottom()
      }
    }

    // Update chat dengan pesan terakhir dari AI
    await saveChat({
      id: selectedChatId.value,
      title: currentChat.title,
      createdAt: currentChat.createdAt,
      lastMessage: generateMessageSummary(currentStreamingMessage.value),
      messages: messages.value
    })

    newMessage.value = ''
    await loadChats()

  } catch (error) {
    if (error.name === 'AbortError') {
      console.log('Request aborted')
    } else {
      console.error('Error:', error)
      // Update pesan error ke UI
      if (messages.value.length > 0) {
        messages.value[messages.value.length - 1].content = 'Sorry, there was an error processing your message. Please try again.'
      }
      // Tampilkan toast error
      showToast('Failed to send message. Please try again.', 'error')
    }
  } finally {
    isStreaming.value = false
    controller = null
    await scrollToBottom()
  }
}

const deleteSelectedChat = async () => {
  if (selectedChatId.value) {
    await deleteChat(selectedChatId.value)
    await loadChats()
    if (chats.value.length > 0) {
      selectedChatId.value = chats.value[0].id
      await loadChatHistory(selectedChatId.value)
    } else {
      await createNewChat()
    }
  }
}

watch(searchTerm, async (newTerm) => {
  if (newTerm) {
    chats.value = (await getChats()).filter(chat => 
      chat.title.toLowerCase().includes(newTerm.toLowerCase()) ||
      chat.lastMessage.toLowerCase().includes(newTerm.toLowerCase())
    )
  } else {
    await loadChats()
  }
})

const generateChatTitle = (message) => {
  // Ambil 30 karakter pertama atau sampai tanda titik pertama
  const title = message.split('.')[0].substring(0, 30)
  return title.length === 30 ? title + '...' : title
}

// Tambah fungsi untuk generate ringkasan pesan
const generateMessageSummary = (message) => {
  const summary = message.split('\n')[0].substring(0, 35)
  return summary.length === 50 ? summary + '...' : summary
}

const createNewChat = async () => {
  const newChat = {
    id: Date.now(),
    title: 'New Chat',
    createdAt: new Date().toISOString(),
    lastMessage: '',
    messages: []
  }
  await saveChat(newChat)
  selectedChatId.value = newChat.id
  messages.value = []
  await loadChats()
}

const copyMessage = async (content) => {
  await copy(content)
  showToast('Message copied to clipboard!', "success")
}

// Edit chat title
const startEditingTitle = (chat) => {
  editingChatId.value = chat.id
  editedTitle.value = chat.title
}

const saveEditedTitle = async (chat) => {
  if (editedTitle.value.trim()) {
    await saveChat({
      ...chat,
      title: editedTitle.value.trim()
    })
    await loadChats()
  }
  editingChatId.value = null
}

const cancelEditingTitle = () => {
  editingChatId.value = null
  editedTitle.value = ''
}

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

onMounted(() => {
  if (controller) {
    controller.abort()
  }

  const handleResize = () => {
    if (window.innerWidth >= 1024) {
      isSidebarOpen.value = false;
    }
  };

  window.addEventListener('resize', handleResize);
  return () => window.removeEventListener('resize', handleResize);
})
</script>

<template>
  <div class="h-screen flex">
    <!-- Backdrop mobile -->
    <div 
      v-if="isSidebarOpen" 
      class="fixed inset-0 bg-black/20 backdrop-blur-sm z-20 lg:hidden"
      @click="toggleSidebar"
    ></div>
    <!-- Sidebar -->
    <div 
      :class="[
          'w-64 bg-slate-800 text-white flex flex-col z-30',
          'fixed inset-y-0 left-0 transform transition-transform duration-300 lg:relative lg:transform-none',
          isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        ]"
    >
      <!-- New Chat Button -->
      <button 
        @click="createNewChat"
        class="m-4 p-3 flex items-center font-semibold gap-2 bg-indigo-700 rounded-lg hover:bg-indigo-600 transition-colors"
      >
        <PlusCircle class="w-5 h-5" />
        <span>New Chat</span>
      </button>

      <!-- Search -->
      <div class="px-4 mb-4">
        <div class="relative">
          <Search class="w-4 h-4 absolute left-3 top-3 text-slate-400" />
          <input
            v-model="searchTerm"
            type="text"
            placeholder="Search chats..."
            class="w-full pl-10 pr-4 py-2 bg-slate-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
      </div>

      <!-- Chat List -->
      <div class="flex-1 overflow-y-auto px-2">
        <div
          v-for="chat in chats"
          :key="chat.id"
          @click="selectedChatId = chat.id; loadChatHistory(chat.id)"
          :class="[
            'p-3 rounded-lg cursor-pointer mb-2 transition-colors group',
            selectedChatId === chat.id 
              ? 'bg-slate-600'
              : 'hover:bg-slate-700'
          ]"
        >
          <div class="flex items-center justify-between">
            <div class="flex-1 mr-2">
              <div v-if="editingChatId === chat.id" class="flex items-center gap-2">
                <input
                  v-model="editedTitle"
                  type="text"
                  class="w-full px-2 py-1 text-sm bg-slate-800 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  @keyup.enter="saveEditedTitle(chat)"
                  @keyup.esc="cancelEditingTitle"
                />
                <button
                  @click.stop="saveEditedTitle(chat)"
                  class="p-1 text-green-400 hover:text-green-300"
                >
                  <Check class="w-4 h-4" />
                </button>
                <button
                  @click.stop="cancelEditingTitle"
                  class="p-1 text-red-400 hover:text-red-300"
                >
                  <X class="w-4 h-4" />
                </button>
              </div>
              <div v-else>
                <div class="text-sm font-medium truncate">{{ chat.title }}</div>
                <div class="text-xs text-slate-400 truncate">
                  {{ chat.lastMessage ? generateMessageSummary(chat.lastMessage) : 'No messages yet' }}
                </div>
              </div>
            </div>
            <button
              v-if="!editingChatId && selectedChatId === chat.id"
              @click.stop="startEditingTitle(chat)"
              class="p-1 text-slate-400 hover:text-slate-300 opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <Edit2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- Settings -->
      <div class="p-4 border-t border-slate-700">
        <button class="flex items-center gap-2 text-slate-300 hover:text-white transition-colors">
          <Settings class="w-5 h-5" />
          <span>Settings</span>
        </button>
      </div>
    </div>

    <!-- Main Chat Area -->
    <div class="flex-1 flex flex-col bg-slate-50">
      <!-- Chat Header with Toggle Button -->
      <div class="px-6 py-4 bg-white border-b border-slate-200 flex justify-between items-center">
        <div class="flex items-center gap-3">
          <!-- Toggle Sidebar Button (Mobile Only) -->
          <button 
            @click="toggleSidebar"
            class="p-2 text-slate-600 hover:bg-slate-100 rounded-lg lg:hidden"
          >
            <Menu class="w-5 h-5" />
          </button>
          <h2 class="text-lg font-semibold text-slate-800">AI Assistant</h2>
        </div>
        <button 
          @click="deleteSelectedChat"
          class="p-2 text-slate-400 hover:text-red-500 transition-colors"
          title="Delete chat"
        >
          <Trash2 class="w-5 h-5" />
        </button>
      </div>

      <!-- Messages -->
      <div 
        ref="messagesContainer"
        class="flex-1 overflow-y-auto p-6 space-y-6"
      >
        <div v-if="messages.length === 0" class="flex justify-center items-center h-full">
          <div class="text-center text-slate-500">
            <p class="text-lg font-medium mb-2">👋 Welcome!</p>
            <p class="text-sm">How can I help you today?</p>
          </div>
        </div>

        <div
          v-for="(message, index) in messages"
          :key="index"
          :class="[
            'flex gap-4',
            message.role === 'user' ? 'justify-end' : 'justify-start'
          ]"
        >
          <!-- AI Avatar -->
          <div v-if="message.role !== 'user'" class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center">
            <Bot class="w-5 h-5 text-indigo-600" />
          </div>
          
          <div 
            :class="[
              'group relative max-w-2xl',
              message.role === 'user' ? 'items-end' : 'items-start'
            ]"
          >
            <div class="text-xs text-slate-400 mb-1">
              {{ new Date(message.timestamp).toLocaleTimeString() }}
            </div>
            
            <div
              :class="[
                'p-4 rounded-2xl',
                message.role === 'user'
                  ? 'bg-indigo-500 text-white'
                  : 'bg-white shadow-sm'
              ]"
            >
              <div
                v-if="message.role === 'assistant' && index === messages.length - 1 && isStreaming"
                class="prose prose-sm max-w-none"
              >
                <div v-html="formatMessage(message.content)" />
                <div class="flex items-center gap-2 mt-2">
                  <div class="w-2 h-2 bg-indigo-500 rounded-full animate-bounce" style="animation-delay: 0ms" />
                  <div class="w-2 h-2 bg-indigo-500 rounded-full animate-bounce" style="animation-delay: 150ms" />
                  <div class="w-2 h-2 bg-indigo-500 rounded-full animate-bounce" style="animation-delay: 300ms" />
                </div>
              </div>
              <div
                v-else
                class="prose prose-sm max-w-none"
                v-html="formatMessage(message.content)"
              />
            </div>

            <!-- Copy button -->
            <button
              v-if="!isStreaming"
              @click="copyMessage(message.content)"
              class="absolute -right-12 top-8 p-2 text-slate-400 hover:text-slate-600 opacity-0 group-hover:opacity-100 transition-opacity"
              :title="copied ? 'Copied!' : 'Copy message'"
            >
              <Copy v-if="!copied" class="w-4 h-4" />
              <Check v-else class="w-4 h-4 text-green-500" />
            </button>
          </div>

          <!-- User Avatar -->
          <div v-if="message.role === 'user'" class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center">
            <User class="w-5 h-5 text-slate-600" />
          </div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="p-6 bg-white border-t border-slate-200">
        <div class="max-w-4xl mx-auto flex gap-4">
          <input
            v-model="newMessage"
            @keyup.enter="sendMessage"
            type="text"
            placeholder="Type your message..."
            :disabled="isStreaming"
            class="flex-1 px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:bg-slate-50"
          />
          <button
            @click="sendMessage"
            :disabled="isStreaming || !newMessage.trim()"
            class="px-6 py-3 bg-indigo-500 text-white rounded-xl hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
          >
            <Send class="w-5 h-5" />
            <span>{{ isStreaming ? 'Sending...' : 'Send' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.overflow-hidden {
  overflow: hidden;
}

@media (max-width: 1024px) {
  .sidebar-scroll {
    height: 100vh;
    overflow-y: auto;
  }
}

@keyframes blink {
  from, to { opacity: 1; }
  50% { opacity: 0; }
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

.animate-bounce {
  animation: bounce 0.6s infinite;
}

.animate-blink {
  animation: blink 1s step-end infinite;
}

.scroll-smooth {
  scroll-behavior: smooth;
}

/* Tailwind Typography Styles */
.prose {
  @apply text-base leading-7;
}

.prose strong {
  @apply font-semibold text-current;
}

.prose em {
  @apply italic text-current;
}

.prose u {
  @apply underline text-current;
}

.prose br {
  @apply block content-[''] my-2;
}

.prose ul {
  @apply list-disc pl-5 my-2;
}

.prose ol {
  @apply list-decimal pl-5 my-2;
}

.prose p {
  @apply mb-4;
}

/* Tambahan untuk list */
.list-item {
  @apply my-2 pl-4;
  display: block;
}

/* Spacing khusus untuk list berurutan */
div.list-item + div.list-item {
  @apply mt-2;
}

/* Pastikan list item terakhir memiliki margin bottom */
div.list-item:last-child {
  @apply mb-2;
}
</style>