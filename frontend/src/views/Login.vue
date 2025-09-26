<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h1>ABC User Management</h1>
          <p>Sign in to continue</p>
        </div>
      </template>
      
      <!-- Admin credentials display for easy testing -->
      <el-alert
        title="Demo Admin Credentials"
        type="info"
        :closable="false"
        show-icon
        class="demo-alert"
      >
        <div class="demo-credentials">
          <div class="credential-item" @click="copyToClipboard('admin@abc.com')">
            <span class="label">Email:</span>
            <span class="value">admin@abc.com</span>
            <el-icon class="copy-icon"><DocumentCopy /></el-icon>
          </div>
          <div class="credential-item" @click="copyToClipboard('Admin@123456')">
            <span class="label">Password:</span>
            <span class="value">Admin@123456</span>
            <el-icon class="copy-icon"><DocumentCopy /></el-icon>
          </div>
        </div>
      </el-alert>
      
      <!-- Login form with validation -->
      <el-form 
        ref="loginForm" 
        :model="credentials" 
        :rules="loginRules"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="email">
          <el-input
            v-model="credentials.email"
            placeholder="Email Address"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="credentials.password"
            type="password"
            placeholder="Password"
            prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            Sign In
          </el-button>
        </el-form-item>
        
        <!-- Quick fill button for demo -->
        <el-button
          type="success"
          size="small"
          @click="fillDemoCredentials"
          style="width: 100%"
        >
          Use Demo Admin Credentials
        </el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'

export default {
  name: 'LoginView',
  setup() {
    const router = useRouter()
    const store = useStore()
    const loginForm = ref(null)
    const loading = ref(false)
    
    // Form data model
    const credentials = ref({
      email: '',
      password: ''
    })
    
    // Validation rules for form fields
    const loginRules = {
      email: [
        { required: true, message: 'Please enter email', trigger: 'blur' },
        { type: 'email', message: 'Please enter valid email', trigger: 'blur' }
      ],
      password: [
        { required: true, message: 'Please enter password', trigger: 'blur' },
        { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
      ]
    }
    
    // Handle login submission
    const handleLogin = async () => {
      // Validate form before submission
      const valid = await loginForm.value.validate().catch(() => false)
      if (!valid) return
      
      loading.value = true
      
      try {
        // Attempt authentication through store action
        await store.dispatch('auth/login', credentials.value)
        ElMessage.success('Login successful!')
        router.push('/dashboard')
      } catch (error) {
        ElMessage.error(error.message || 'Login failed. Please check your credentials.')
      } finally {
        loading.value = false
      }
    }
    
    // Copy text to clipboard functionality
    const copyToClipboard = (text) => {
      navigator.clipboard.writeText(text)
      ElMessage.success('Copied to clipboard!')
    }
    
    // Auto-fill demo credentials for easy testing
    const fillDemoCredentials = () => {
      credentials.value = {
        email: 'admin@abc.com',
        password: 'Admin@123456'
      }
      ElMessage.info('Demo credentials filled!')
    }
    
    return {
      loginForm,
      credentials,
      loginRules,
      loading,
      handleLogin,
      copyToClipboard,
      fillDemoCredentials
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 450px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  border-radius: 12px;
}

.card-header {
  text-align: center;
}

.card-header h1 {
  color: #303133;
  margin-bottom: 8px;
  font-size: 24px;
}

.card-header p {
  color: #909399;
  font-size: 14px;
}

.demo-alert {
  margin-bottom: 24px;
}

.demo-credentials {
  margin-top: 12px;
}

.credential-item {
  display: flex;
  align-items: center;
  padding: 8px;
  margin: 4px 0;
  background: #f5f7fa;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.3s;
}

.credential-item:hover {
  background: #e9ebef;
}

.credential-item .label {
  font-weight: 600;
  margin-right: 8px;
  min-width: 70px;
}

.credential-item .value {
  flex: 1;
  font-family: 'Courier New', monospace;
}

.copy-icon {
  margin-left: 8px;
  color: #409eff;
}
</style>