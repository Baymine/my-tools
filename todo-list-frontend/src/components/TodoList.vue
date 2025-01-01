<template>
  <div class="todo-list">
    <TodoForm @add-todo="addTodo" />
    <TodoFilter @filter-change="setFilter" />
    <ul>
      <TodoItem
        v-for="todo in sortedTodos"
        :key="todo.id"
        :todo="todo"
        @toggle-complete="toggleComplete(todo)"
        @update-priority="updatePriority(todo, $event)"
        @remove="removeTodo(todo)"
      />
    </ul>
  </div>
</template>

<script>
import TodoItem from './TodoItem.vue'
import TodoForm from './TodoForm.vue'
import TodoFilter from './TodoFilter.vue'
import axios from 'axios'

const API_URL = 'http://localhost:8081'

export default {
  name: 'TodoList',
  components: {
    TodoItem,
    TodoForm,
    TodoFilter
  },
  data() {
    return {
      todos: [],
      filter: 'all'
    }
  },
  computed: {
    filteredTodos() {
      if (this.filter === 'active') {
        return this.todos.filter(todo => !todo.completed)
      } else if (this.filter === 'completed') {
        return this.todos.filter(todo => todo.completed)
      }
      return this.todos
    },
    sortedTodos() {
      const priorityOrder = { high: 3, medium: 2, low: 1 }
      return this.filteredTodos.sort((a, b) => {
        if (a.completed === b.completed) {
          return priorityOrder[b.priority] - priorityOrder[a.priority]
        }
        return a.completed ? 1 : -1
      })
    }
  },
  methods: {
    async fetchTodos() {
      try {
        const response = await axios.get(`${API_URL}/todos`)
        this.todos = response.data
      } catch (error) {
        console.error('Error fetching todos:', error)
      }
    },
    async addTodo(title) {
      try {
        const response = await axios.post(`${API_URL}/todos`, {
          title: title,
          completed: false,
          priority: 'medium'
        })
        this.todos.push(response.data)
      } catch (error) {
        console.error('Error adding todo:', error)
      }
    },
    async toggleComplete(todo) {
      try {
        await axios.put(`${API_URL}/todos/${todo.id}`, {
          ...todo,
          completed: !todo.completed
        })
        todo.completed = !todo.completed
      } catch (error) {
        console.error('Error updating todo:', error)
      }
    },
    async updatePriority(todo, priority) {
      try {
        await axios.put(`${API_URL}/todos/${todo.id}`, {
          ...todo,
          priority: priority
        })
        todo.priority = priority
      } catch (error) {
        console.error('Error updating todo priority:', error)
      }
    },
    async removeTodo(todo) {
      try {
        await axios.delete(`${API_URL}/todos/${todo.id}`)
        this.todos = this.todos.filter(t => t.id !== todo.id)
      } catch (error) {
        console.error('Error removing todo:', error)
      }
    },
    setFilter(filter) {
      this.filter = filter
    }
  },
  mounted() {
    this.fetchTodos()
  }
}
</script>

<style scoped>
.todo-list {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

ul {
  list-style-type: none;
  padding: 0;
}
</style>
