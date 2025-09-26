<template>
  <div class="user-list-container">
    <!-- Search bar section -->
    <div class="search-section">
      <el-input
        v-model="searchQuery"
        placeholder="Search users by name or email..."
        clearable
        size="large"
        @input="debouncedSearch"
        @clear="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>
    
    <!-- Users table with loading state -->
    <el-table
      :data="users"
      v-loading="loading"
      style="width: 100%"
      @sort-change="handleSortChange"
      empty-text="No users found"
    >
      <!-- ID column -->
      <el-table-column
        prop="id"
        label="ID"
        width="80"
        sortable="custom"
      />
      
      <!-- Name column with search highlight -->
      <el-table-column
        prop="name"
        label="Name"
        sortable="custom"
      >
        <template #default="scope">
          <span v-html="highlightText(scope.row.name)"></span>
        </template>
      </el-table-column>
      
      <!-- Email column with search highlight -->
      <el-table-column
        prop="email"
        label="Email"
        sortable="custom"
      >
        <template #default="scope">
          <span v-html="highlightText(scope.row.email)"></span>
        </template>
      </el-table-column>
      
      <!-- Age column -->
      <el-table-column
        prop="age"
        label="Age"
        width="100"
        sortable="custom"
      />
      
      <!-- Role column with colored tags -->
      <el-table-column
        prop="role"
        label="Role"
        width="120"
      >
        <template #default="scope">
          <el-tag
            :type="getRoleTagType(scope.row.role)"
            effect="dark"
          >
            {{ scope.row.role }}
          </el-tag>
        </template>
      </el-table-column>
      
      <!-- Status column -->
      <el-table-column
        prop="active"
        label="Status"
        width="100"
      >
        <template #default="scope">
          <el-tag
            :type="scope.row.active ? 'success' : 'danger'"
            effect="light"
          >
            {{ scope.row.active ? 'Active' : 'Inactive' }}
          </el-tag>
        </template>
      </el-table-column>
      
      <!-- Created date column -->
      <el-table-column
        prop="created_at"
        label="Created"
        width="180"
        sortable="custom"
      >
        <template #default="scope">
          {{ formatDate(scope.row.created_at) }}
        </template>
      </el-table-column>
      
      <!-- Actions column -->
      <el-table-column
        label="Actions"
        width="200"
        fixed="right"
        align="center"
      >
        <template #default="scope">
          <el-button-group>
            <el-button
              type="primary"
              size="small"
              icon="Edit"
              @click="handleEdit(scope.row)"
            >
              Edit
            </el-button>
            <el-button
              type="danger"
              size="small"
              icon="Delete"
              @click="handleDelete(scope.row)"
            >
              Delete
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    
    <!-- Pagination component -->
    <Pagination
      v-model="pagination"
      :total="totalUsers"
      @change="fetchUsers"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import Pagination from './Pagination.vue'
import { debounce } from 'lodash'

export default {
  name: 'UserListComponent',
  components: {
    Pagination
  },
  emits: ['edit', 'delete'],
  setup(props, { emit }) {
    // State management
    const users = ref([])
    const loading = ref(false)
    const searchQuery = ref('')
    const totalUsers = ref(0)
    const pagination = ref({
      page: 1,
      limit: 10
    })
    const sortConfig = ref({
      prop: 'created_at',
      order: 'descending'
    })
    
    // Fetch users from API
    const fetchUsers = async () => {
      loading.value = true
      
      try {
        // Build query parameters
        const params = {
          ...pagination.value,
          search: searchQuery.value,
          sort: sortConfig.value.prop,
          order: sortConfig.value.order === 'ascending' ? 'asc' : 'desc'
        }
        
        // Make API call
        const response = await api.getUsers(params)
        users.value = response.data
        totalUsers.value = response.total
      } catch (error) {
        ElMessage.error('Failed to load users')
        console.error('Error fetching users:', error)
      } finally {
        loading.value = false
      }
    }
    
    // Debounced search to reduce API calls
    const debouncedSearch = debounce(() => {
      pagination.value.page = 1 // Reset to first page on search
      fetchUsers()
    }, 500)
    
    // Handle immediate search on clear
    const handleSearch = () => {
      if (!searchQuery.value) {
        fetchUsers()
      }
    }
    
    // Handle table sorting
    const handleSortChange = ({ prop, order }) => {
      sortConfig.value = { prop, order }
      fetchUsers()
    }
    
    // Highlight search text in results
    const highlightText = (text) => {
      if (!searchQuery.value || !text) return text
      
      const regex = new RegExp(`(${searchQuery.value})`, 'gi')
      return text.replace(regex, '<mark>$1</mark>')
    }
    
    // Get tag type based on role
    const getRoleTagType = (role) => {
      const types = {
        admin: 'danger',
        manager: 'warning',
        user: 'primary'
      }
      return types[role] || 'info'
    }
    
    // Format date for display
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }
    
    // Handle edit button click
    const handleEdit = (user) => {
      emit('edit', user)
    }
    
    // Handle delete button click
    const handleDelete = (user) => {
      emit('delete', user)
    }
    
    // Load initial data
    onMounted(() => {
      fetchUsers()
    })
    
    return {
      users,
      loading,
      searchQuery,
      pagination,
      totalUsers,
      debouncedSearch,
      handleSearch,
      handleSortChange,
      highlightText,
      getRoleTagType,
      formatDate,
      handleEdit,
      handleDelete,
      fetchUsers
    }
  }
}
</script>

<style scoped>
.user-list-container {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.search-section {
  margin-bottom: 20px;
}

/* Highlight search results */
:deep(mark) {
  background-color: #ffd04b;
  padding: 0 2px;
  border-radius: 2px;
}

/* Custom table styles */
:deep(.el-table) {
  font-size: 14px;
}

:deep(.el-table th) {
  background-color: #f5f7fa;
  font-weight: 600;
}
</style>