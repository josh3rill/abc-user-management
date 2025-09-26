import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8085/api/v1'

// Configure axios instance for user endpoints
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add authentication token to all requests
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

export default {
  // Get paginated users with search
  getUsers(params = {}) {
    return apiClient.get('/users', { params })
  },
  
  // Get single user by ID
  getUser(id) {
    return apiClient.get(`/users/${id}`)
  },
  
  // Create new user
  createUser(userData) {
    return apiClient.post('/users', userData)
  },
  
  // Update existing user
  updateUser(id, userData) {
    return apiClient.put(`/users/${id}`, userData)
  },
  
  // Delete user
  deleteUser(id) {
    return apiClient.delete(`/users/${id}`)
  }
}