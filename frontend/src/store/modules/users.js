import api from '../../api/users'

export default {
  namespaced: true,
  
  state: {
    users: [],
    currentUser: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0
    }
  },
  
  mutations: {
    SET_USERS(state, users) {
      state.users = users
    },
    
    SET_PAGINATION(state, pagination) {
      state.pagination = { ...state.pagination, ...pagination }
    },
    
    SET_CURRENT_USER(state, user) {
      state.currentUser = user
    }
  },
  
  actions: {
    async fetchUsers({ commit }, params = {}) {
      try {
        const response = await api.getUsers(params)
        commit('SET_USERS', response.data.data)
        commit('SET_PAGINATION', response.data.pagination)
        return response.data
      } catch (error) {
        console.error('Failed to fetch users:', error)
        throw error
      }
    },
    
    async createUser({ dispatch }, userData) {
      try {
        await api.createUser(userData)
        // Refresh user list after creation
        await dispatch('fetchUsers')
      } catch (error) {
        console.error('Failed to create user:', error)
        throw error
      }
    },
    
    async updateUser({ dispatch }, { id, data }) {
      try {
        await api.updateUser(id, data)
        // Refresh list to show updated data
        await dispatch('fetchUsers')
      } catch (error) {
        console.error('Failed to update user:', error)
        throw error
      }
    },
    
    async deleteUser({ dispatch }, id) {
      try {
        await api.deleteUser(id)
        // Refresh list after deletion
        await dispatch('fetchUsers')
      } catch (error) {
        console.error('Failed to delete user:', error)
        throw error
      }
    }
  }
}