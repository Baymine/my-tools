<template>
  <form @submit.prevent="addTodo" class="todo-form">
    <input v-model="newTodo" type="text" placeholder="输入新的待办事项" required class="todo-input">
    <button type="submit" class="add-btn">添加</button>
  </form>
</template>

<script>
export default {
  name: 'TodoForm',
  data() {
    return {
      newTodo: ''
    }
  },
  methods: {
    addTodo() {
      if (this.newTodo.trim()) {
        this.$emit('add-todo', this.newTodo)
        this.newTodo = ''
      }
    }
  }
}
</script>

<style scoped>
.todo-form {
  display: flex;
  margin-bottom: 20px;
}

.todo-input {
  flex-grow: 1;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ddd;
  border-radius: 4px 0 0 4px;
  outline: none;
}

.todo-input:focus {
  border-color: #4CAF50;
}

.add-btn {
  padding: 10px 20px;
  font-size: 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 0 4px 4px 0;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.add-btn:hover {
  background-color: #45a049;
}
</style>
