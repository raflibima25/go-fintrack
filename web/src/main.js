import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/tailwind.css'
import './assets/index.css'

const app = createApp(App)

// Error handler global
app.config.errorHandler = (err, vm, info) => {
  console.error('Global error:', err, info)
}

app.use(router)
app.mount('#app')
