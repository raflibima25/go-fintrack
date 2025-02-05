<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { marked } from 'marked'

const messages = ref([])
const newMessage = ref('')
const isStreaming = ref(false)
let currentStreamingMessage = ref('')
let controller = null
const messagesContainer = ref(null)

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

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
        message += data + ' '  // Add space after each data chunk
      }
    }
  }

  return message
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || isStreaming.value) return

  messages.value.push({
    role: 'user',
    content: newMessage.value,
    timestamp: new Date().toLocaleTimeString()
  })

  messages.value.push({
    role: 'assistant',
    content: '',
    timestamp: new Date().toLocaleTimeString()
  })

  await scrollToBottom()
  isStreaming.value = true
  currentStreamingMessage.value = ''

  try {
    if (controller) {
      controller.abort()
    }

    controller = new AbortController()

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

    newMessage.value = ''
  } catch (error) {
    if (error.name === 'AbortError') {
      console.log('Request aborted')
    } else {
      console.error('Error:', error)
      messages.value[messages.value.length - 1].content = 'Maaf, terjadi kesalahan dalam memproses pesan.'
    }
  } finally {
    isStreaming.value = false
    controller = null
    await scrollToBottom()
  }
}

onMounted(() => {
  if (controller) {
    controller.abort()
  }
})
</script>

<template>
  <div class="max-w-4xl mx-auto h-[600px] flex flex-col border border-slate-200 rounded-lg shadow-sm">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-slate-200 bg-white">
      <h2 class="text-lg font-semibold text-slate-800">AI Assistant</h2>
      <p class="text-sm text-slate-500">Powered by Deepseek</p>
    </div>

    <!-- Messages Container -->
    <div 
      ref="messagesContainer"
      class="flex-1 overflow-y-auto p-4 bg-slate-50 space-y-4 scroll-smooth"
    >
      <!-- Welcome Message -->
      <div v-if="messages.length === 0" class="flex justify-center items-center h-full">
        <div class="text-center text-slate-500">
          <p class="text-lg font-medium mb-2">👋 Selamat datang!</p>
          <p class="text-sm">Saya siap membantu Anda. Silakan ajukan pertanyaan.</p>
        </div>
      </div>

      <!-- Chat Messages -->
      <div v-for="(message, index) in messages"
           :key="index"
           :class="[
             'flex',
             message.role === 'user' ? 'justify-end' : 'justify-start'
           ]">
        <div :class="[
          'max-w-[80%] group',
          message.role === 'user' ? 'items-end' : 'items-start'
        ]">
          <!-- Timestamp -->
          <div :class="[
            'text-xs text-slate-400 mb-1',
            message.role === 'user' ? 'text-right' : 'text-left'
          ]">
            {{ message.timestamp }}
          </div>
          
          <!-- Message Content -->
          <div :class="[
            'p-3 rounded-lg',
            message.role === 'user' 
              ? 'bg-blue-500 text-white rounded-br-none' 
              : 'bg-white text-slate-800 rounded-bl-none shadow-sm'
          ]">
            <div v-if="message.role === 'assistant' && index === messages.length - 1 && isStreaming"
                 class="prose prose-sm max-w-none"
                 v-html="formatMessage(message.content) + '<span class=\'inline-block w-2 animate-blink\'>▋</span>'">
            </div>
            <div v-else 
                 class="prose prose-sm max-w-none"
                 v-html="formatMessage(message.content)">
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Input Container -->
    <div class="p-4 border-t border-slate-200 bg-white">
      <div class="flex gap-3">
        <input 
          v-model="newMessage"
          @keyup.enter="sendMessage"
          type="text"
          placeholder="Ketik pesan Anda di sini..."
          :disabled="isStreaming"
          class="flex-1 px-4 py-2 border border-slate-200 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-slate-50 disabled:cursor-not-allowed"
        />
        <button
          @click="sendMessage" 
          :disabled="isStreaming || !newMessage.trim()"
          class="px-6 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
        >
          <span>{{ isStreaming ? 'Mengirim...' : 'Kirim' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style>
@keyframes blink {
  from, to { opacity: 1; }
  50% { opacity: 0; }
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