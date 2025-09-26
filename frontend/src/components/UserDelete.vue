<template>
  <el-dialog
    :model-value="visible"
    title="Confirm Delete"
    width="400px"
    @close="handleClose"
  >
    <div class="delete-content">
      <el-icon class="warning-icon" :size="48" color="#F56C6C">
        <WarningFilled />
      </el-icon>
      <p class="delete-message">
        Are you sure you want to delete user 
        <strong>{{ user?.name || user?.email }}</strong>?
      </p>
      <p class="delete-warning">
        This action cannot be undone.
      </p>
    </div>
    
    <template #footer>
      <el-button @click="handleClose">Cancel</el-button>
      <el-button type="danger" @click="handleDelete" :loading="loading">
        Delete User
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
import { ref } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'

export default {
  name: 'UserDelete',
  props: {
    visible: Boolean,
    user: Object
  },
  emits: ['update:visible', 'deleted'],
  setup(props, { emit }) {
    const store = useStore()
    const loading = ref(false)
    
    // Process user deletion
    const handleDelete = async () => {
      if (!props.user) return
      
      loading.value = true
      
      try {
        await store.dispatch('users/deleteUser', props.user.id)
        emit('deleted')
        handleClose()
      } catch (error) {
        ElMessage.error(error.message || 'Failed to delete user')
      } finally {
        loading.value = false
      }
    }
    
    // Close dialog
    const handleClose = () => {
      emit('update:visible', false)
    }
    
    return {
      loading,
      handleDelete,
      handleClose
    }
  }
}
</script>

<style scoped>
.delete-content {
  text-align: center;
  padding: 20px 0;
}

.warning-icon {
  margin-bottom: 20px;
}

.delete-message {
  font-size: 16px;
  margin-bottom: 12px;
  color: #303133;
}

.delete-warning {
  color: #909399;
  font-size: 14px;
}
</style>