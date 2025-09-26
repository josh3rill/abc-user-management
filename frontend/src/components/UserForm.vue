<template>
  <el-dialog
    :model-value="visible"
    :title="mode === 'create' ? 'Add New User' : 'Edit User'"
    width="500px"
    @close="handleClose"
  >
    <el-form
      ref="userForm"
      :model="formData"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="Name" prop="name">
        <el-input v-model="formData.name" placeholder="Enter full name" />
      </el-form-item>
      
      <el-form-item label="Email" prop="email">
        <el-input v-model="formData.email" placeholder="user@example.com" />
      </el-form-item>
      
      <el-form-item label="Age" prop="age">
        <el-input-number
          v-model="formData.age"
          :min="1"
          :max="120"
          placeholder="Must be 18 or older"
        />
      </el-form-item>
      
      <el-form-item label="Password" prop="password" v-if="mode === 'create'">
        <el-input
          v-model="formData.password"
          type="password"
          placeholder="Minimum 6 characters"
          show-password
        />
      </el-form-item>
      
      <el-form-item label="Role" prop="role">
        <el-select v-model="formData.role" placeholder="Select role">
          <el-option label="User" value="user" />
          <el-option label="Admin" value="admin" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="Active" prop="active">
        <el-switch v-model="formData.active" />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <el-button @click="handleClose">Cancel</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">
        {{ mode === 'create' ? 'Create' : 'Update' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
import { ref, watch } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'

export default {
  name: 'UserForm',
  props: {
    visible: Boolean,
    user: Object,
    mode: {
      type: String,
      default: 'create'
    }
  },
  emits: ['update:visible', 'saved'],
  setup(props, { emit }) {
    const store = useStore()
    const userForm = ref(null)
    const loading = ref(false)
    
    // Form data with defaults
    const formData = ref({
      name: '',
      email: '',
      age: 18,
      password: '',
      role: 'user',
      active: true
    })
    
    // Custom validation for age requirement
    const validateAge = (rule, value, callback) => {
      if (value < 18) {
        callback(new Error('User must be 18 years or older'))
      } else {
        callback()
      }
    }
    
    // Form validation rules
    const rules = {
      name: [
        { required: true, message: 'Name is required', trigger: 'blur' },
        { min: 2, max: 100, message: 'Name must be 2-100 characters', trigger: 'blur' }
      ],
      email: [
        { required: true, message: 'Email is required', trigger: 'blur' },
        { type: 'email', message: 'Invalid email format', trigger: 'blur' }
      ],
      age: [
        { required: true, message: 'Age is required', trigger: 'blur' },
        { validator: validateAge, trigger: 'blur' }
      ],
      password: [
        { required: props.mode === 'create', message: 'Password is required', trigger: 'blur' },
        { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
      ],
      role: [
        { required: true, message: 'Role is required', trigger: 'change' }
      ]
    }
    
    // Watch for user prop changes to populate form
    watch(() => props.user, (newUser) => {
      if (newUser && props.mode === 'edit') {
        formData.value = { ...newUser }
      } else {
        // Reset form for new user
        formData.value = {
          name: '',
          email: '',
          age: 18,
          password: '',
          role: 'user',
          active: true
        }
      }
    }, { immediate: true })
    
    // Handle form submission
    const handleSubmit = async () => {
      // Validate all fields first
      const valid = await userForm.value.validate().catch(() => false)
      if (!valid) return
      
      loading.value = true
      
      try {
        if (props.mode === 'create') {
          await store.dispatch('users/createUser', formData.value)
        } else {
          await store.dispatch('users/updateUser', {
            id: props.user.id,
            data: formData.value
          })
        }
        emit('saved')
        handleClose()
      } catch (error) {
        ElMessage.error(error.message || 'Operation failed')
      } finally {
        loading.value = false
      }
    }
    
    // Close dialog and reset form
    const handleClose = () => {
      userForm.value?.resetFields()
      emit('update:visible', false)
    }
    
    return {
      userForm,
      formData,
      rules,
      loading,
      handleSubmit,
      handleClose
    }
  }
}
</script>