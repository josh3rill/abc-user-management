<template>
  <div class="pagination-container">
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="pageSizes"
      :total="total"
      :layout="layout"
      :background="background"
      :disabled="disabled"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'PaginationComponent',
  props: {
    modelValue: {
      type: Object,
      default: () => ({
        page: 1,
        limit: 10
      })
    },
    total: {
      type: Number,
      required: true
    },
    pageSizes: {
      type: Array,
      default: () => [10, 20, 50, 100]
    },
    layout: {
      type: String,
      default: 'total, sizes, prev, pager, next, jumper'
    },
    background: {
      type: Boolean,
      default: true
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue', 'change'],
  setup(props, { emit }) {
    // Local state synced with v-model
    const currentPage = ref(props.modelValue.page)
    const pageSize = ref(props.modelValue.limit)
    
    // Watch for external changes to pagination
    watch(() => props.modelValue, (newVal) => {
      currentPage.value = newVal.page
      pageSize.value = newVal.limit
    }, { deep: true })
    
    // Handle page size change
    const handleSizeChange = (val) => {
      pageSize.value = val
      currentPage.value = 1 // Reset to first page
      emitUpdate()
    }
    
    // Handle page number change
    const handleCurrentChange = (val) => {
      currentPage.value = val
      emitUpdate()
    }
    
    // Emit updated pagination values
    const emitUpdate = () => {
      const pagination = {
        page: currentPage.value,
        limit: pageSize.value
      }
      
      emit('update:modelValue', pagination)
      emit('change', pagination)
    }
    
    return {
      currentPage,
      pageSize,
      handleSizeChange,
      handleCurrentChange
    }
  }
}
</script>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 10px 0;
}
</style>