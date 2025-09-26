import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// Initialize Vue application with all plugins
const app = createApp(App)

// Register Element Plus UI library
app.use(ElementPlus)

// Register all icons globally for convenience
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// Setup routing and state management
app.use(router)
app.use(store)

// Mount application to DOM
app.mount('#app')