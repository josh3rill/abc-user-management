import { createStore } from 'vuex'
import auth from './modules/auth'
import users from './modules/users'

// Create Vuex store with modules for separation of concerns
export default createStore({
  modules: {
    auth,
    users
  },
  
  // Global state if needed
  state: {
    appName: 'ABC User Management'
  },
  
  // Enable strict mode in development for better debugging
  strict: process.env.NODE_ENV !== 'production'
})