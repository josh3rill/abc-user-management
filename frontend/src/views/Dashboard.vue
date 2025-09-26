<template>
  <div class="dashboard-container">
    <!-- Navigation header with user info -->
    <el-header class="dashboard-header">
      <div class="header-content">
        <h1>ABC User Management System</h1>
        <div class="user-info">
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-icon><UserFilled /></el-icon>
              {{ currentUser?.name || currentUser?.email }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">Profile</el-dropdown-item>
                <el-dropdown-item command="logout" divided>Logout</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-header>
    
    <!-- Main content area with user management -->
    <el-main class="dashboard-main">
      <Users />
    </el-main>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import Users from './Users.vue'

export default {
  name: 'DashboardView',
  components: {
    Users
  },
  setup() {
    const router = useRouter()
    const store = useStore()
    
    // Get current user from store
    const currentUser = computed(() => store.state.auth.user)
    
    // Handle dropdown menu actions
    const handleCommand = (command) => {
      switch (command) {
        case 'logout':
          store.dispatch('auth/logout')
          ElMessage.success('Logged out successfully')
          router.push('/login')
          break
        case 'profile':
          ElMessage.info('Profile page coming soon!')
          break
      }
    }
    
    return {
      currentUser,
      handleCommand
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.dashboard-header {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  padding: 0 24px;
  height: 60px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-content h1 {
  font-size: 20px;
  color: #303133;
  margin: 0;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
}

.user-dropdown:hover {
  color: #409eff;
}

.dashboard-main {
  padding: 24px;
}
</style>