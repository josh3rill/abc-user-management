import api from '../../api/auth'
import router from '../../router'

export default {
  namespaced: true,
  
  state: {
    token: localStorage.getItem('token') || null,
    user: null,
    isAuthenticated: false
  },
  
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
      state.isAuthenticated = !!token
      
      // Persist token to localStorage
      if (token) {
        localStorage.setItem('token', token)
      } else {
        localStorage.removeItem('token')
      }
    },
    
    SET_USER(state, user) {
      state.user = user
    }
  },
  
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await api.login(credentials)
        const { token, user } = response.data
        
        commit('SET_TOKEN', token)
        commit('SET_USER', user)
        
        // Set default authorization header
        api.setAuthHeader(token)
        
        return response
      } catch (error) {
        commit('SET_TOKEN', null)
        commit('SET_USER', null)
        throw error
      }
    },
    
    logout({ commit }) {
      commit('SET_TOKEN', null)
      commit('SET_USER', null)
      api.setAuthHeader(null)
      router.push('/login')
    },
    
    async validateToken({ commit, state }) {
      if (!state.token) return false
      
      try {
        // Validate token with backend
        api.setAuthHeader(state.token)
        // Token is valid, keep session
        return true
      } catch {
        // Token invalid, clear session
        commit('SET_TOKEN', null)
        commit('SET_USER', null)
        return false
      }
    }
  },
  
  getters: {
    isAuthenticated: state => state.isAuthenticated,
    currentUser: state => state.user,
    token: state => state.token
  }
}