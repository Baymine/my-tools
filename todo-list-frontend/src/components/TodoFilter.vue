<template>
  <div class="todo-filter">
    <button 
      v-for="option in filterOptions" 
      :key="option.value"
      @click="setFilter(option.value)"
      :class="{ active: filter === option.value }"
      class="filter-btn"
    >
      {{ option.label }}
    </button>
  </div>
</template>

<script>
export default {
  name: 'TodoFilter',
  data() {
    return {
      filter: 'all',
      filterOptions: [
        { value: 'all', label: '全部' },
        { value: 'active', label: '未完成' },
        { value: 'completed', label: '已完成' }
      ]
    }
  },
  methods: {
    setFilter(value) {
      this.filter = value
      this.$emit('filter-change', value)
    }
  }
}
</script>

<style scoped>
.todo-filter {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.filter-btn {
  padding: 8px 16px;
  margin: 0 5px;
  font-size: 14px;
  background-color: #f0f0f0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.filter-btn:hover {
  background-color: #e0e0e0;
}

.filter-btn.active {
  background-color: #4CAF50;
  color: white;
}
</style>
