import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import pinia from './stores'
import { useAuthStore } from '@/features/auth/stores/authStore'

const app = createApp(App)

app.use(pinia)
app.use(router)

// Initialize auth store to restore user from localStorage
const authStore = useAuthStore()
authStore.initializeAuth()

app.mount('#app')
