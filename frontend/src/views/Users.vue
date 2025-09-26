<template>
  <div class="users-container">
    <!-- Page header with actions -->
    <div class="page-header">
      <h2>User Management</h2>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        Add New User
      </el-button>
    </div>
    
    <!-- Search bar for filtering users -->
    <el-card class="search-card">
      <el-input
        v-model="searchQuery"
        placeholder="Search by name or email..."
        size="large"
        clearable
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </el-card>
    
    <!-- User list table with actions -->
    <el-card>
      <el-table
        :data="users"
        v-loading="loading"
        style="width: 100%"
        empty-text="No users found"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="Name" />
        <el-table-column prop="email" label="Email" />
        <el-table-column prop="age" label="Age" width="80" />
        <el-table-column prop="role" label="Role" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.role === 'admin' ? 'danger' : 'primary'">
              {{ scope.row.role }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Status" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.active ? 'success' : 'info'">
              {{ scope.row.active ? 'Active' : 'Inactive' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="180" fixed="right">
          <template #default="scope">
            <el-button
              size="small"
              type="primary"
              @click="handleEdit(scope.row)"
            >
              Edit
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(scope.row)"
            >
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- Pagination controls -->
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="totalUsers"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="pagination"
      />
    </el-card>
    
    <!-- User form dialog for create/edit -->
    <UserForm
      v-model:visible="formVisible"
      :user="selectedUser"
      :mode="formMode"
      @saved="handleUserSaved"
    />
    
    <!-- Delete confirmation dialog -->
    <UserDelete
      v-model:visible="deleteVisible"
      :user="selectedUser"
      @deleted="handleUserDeleted"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import UserForm from '../components/UserForm.vue'
import UserDelete from '../components/UserDelete.vue'

export default {
  name: 'UsersView',
  components: {
    UserForm,
    UserDelete
  },
  setup() {
    const store = useStore()
    
    // Data state
    const users = ref([])
    const loading = ref(false)
    const searchQuery = ref('')
    const currentPage = ref(1)
    const pageSize = ref(10)
    const totalUsers = ref(0)
    
    // Dialog state
    const formVisible = ref(false)
    const deleteVisible = ref(false)
    const selectedUser = ref(null)
    const formMode = ref('create')
    
    // Load users from API
    const fetchUsers = async () => {
      loading.value = true
      try {
        const response = await store.dispatch('users/fetchUsers', {
          page: currentPage.value,
          limit: pageSize.value,
          search: searchQuery.value
        })
        users.value = response.data
        totalUsers.value = response.pagination.total
      } catch (error) {
        ElMessage.error('Failed to load users')
      } finally {
        loading.value = false
      }
    }
    
    // Search handling with debounce
    let searchTimeout = null
    const handleSearch = () => {
      clearTimeout(searchTimeout)
      searchTimeout = setTimeout(() => {
        currentPage.value = 1
        fetchUsers()
      }, 500)
    }
    
    // Pagination handlers
    const handlePageChange = () => fetchUsers()
    const handleSizeChange = () => {
      currentPage.value = 1
      fetchUsers()
    }
    
    // CRUD operations
    const showCreateDialog = () => {
      selectedUser.value = null
      formMode.value = 'create'
      formVisible.value = true
    }
    
    const handleEdit = (user) => {
      selectedUser.value = { ...user }
      formMode.value = 'edit'
      formVisible.value = true
    }
    
    const handleDelete = (user) => {
      selectedUser.value = user
      deleteVisible.value = true
    }
    
    const handleUserSaved = () => {
      formVisible.value = false
      fetchUsers()
      ElMessage.success(`User ${formMode.value === 'create' ? 'created' : 'updated'} successfully`)
    }
    
    const handleUserDeleted = () => {
      deleteVisible.value = false
      fetchUsers()
      ElMessage.success('User deleted successfully')
    }
    
    // Load initial data
    onMounted(() => {
      fetchUsers()
    })
    
    return {
      users,
      loading,
      searchQuery,
      currentPage,
      pageSize,
      totalUsers,
      formVisible,
      deleteVisible,
      selectedUser,
      formMode,
      handleSearch,
      handlePageChange,
      handleSizeChange,
      showCreateDialog,
      handleEdit,
      handleDelete,
      handleUserSaved,
      handleUserDeleted
    }
  }
}
</script>

<style scoped>
.users-container {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  color: #303133;
}

.search-card {
  margin-bottom: 24px;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}
</style>