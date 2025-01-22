import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import routes from './router'
import './assets/tailwind.css'

const router = createRouter({
  history: createWebHistory(),
  routes
})

app.config.errorHandler = (err, vm, info) => {
  console.error('Global error:', err, info)
}

const app = createApp(App)
app.use(router)
app.mount('#app')
