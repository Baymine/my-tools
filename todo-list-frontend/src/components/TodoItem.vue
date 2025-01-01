<template>
  <li class="todo-item" :class="{ 'completed': todo.completed, 'high-priority': todo.priority === 'high', 'medium-priority': todo.priority === 'medium', 'low-priority': todo.priority === 'low' }">
    <div class="todo-content">
      <input type="checkbox" :checked="todo.completed" @change="toggleComplete" class="todo-checkbox">
      <span class="todo-title">{{ todo.title }}</span>
    </div>
    <div class="todo-actions">
      <select v-model="todo.priority" @change="updatePriority" class="priority-select">
        <option value="low">低</option>
        <option value="medium">中</option>
        <option value="high">高</option>
      </select>
      <button @click="$emit('remove')" class="remove-btn">删除</button>
    </div>
  </li>
</template>

<script>
export default {
  name: 'TodoItem',
  props: ['todo'],
  methods: {
    toggleComplete() {
      this.$emit('toggle-complete')
    },
    updatePriority() {
      this.$emit('update-priority', this.todo.priority)
    }
  }
}
</script>

<style scoped>
.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  margin-bottom: 10px;
  background-color: #f8f8f8;
  border-radius: 5px;
  transition: all 0.3s ease;
}

.todo-content {
  display: flex;
  align-items: center;
}

.todo-checkbox {
  margin-right: 10px;
}

.todo-title {
  font-size: 16px;
}

.todo-actions {
  display: flex;
  align-items: center;
}

.priority-select {
  margin-right: 10px;
  padding: 5px;
  border-radius: 3px;
  border: 1px solid #ddd;
}

.remove-btn {
  padding: 5px 10px;
  background-color: #ff4d4d;
  color: white;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.remove-btn:hover {
  background-color: #ff3333;
}

.completed {
  opacity: 0.6;
}

.completed .todo-title {
  text-decoration: line-through;
}

.high-priority {
  border-left: 5px solid #ff4d4d;
}

.medium-priority {
  border-left: 5px solid #ffa500;
}

.low-priority {
  border-left: 5px solid #4caf50;
}
</style>
