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
      console.log('Fetching todos...')
      try {
        const response = await axios.get(`${API_URL}/todos`)
        console.log('Todos fetched successfully:', response.data)
        this.todos = response.data
        console.log('Updated todos:', this.todos)
      } catch (error) {
        console.error('Error fetching todos:', error)
        if (error.response) {
          console.error('Response data:', error.response.data)
          console.error('Response status:', error.response.status)
          console.error('Response headers:', error.response.headers)
        } else if (error.request) {
          console.error('No response received:', error.request)
        } else {
          console.error('Error setting up request:', error.message)
        }
      }
    },
    async addTodo(title) {
      console.log('Adding todo:', title)
      try {
        const response = await axios.post(`${API_URL}/todos`, {
          title: title,
          completed: false,
          priority: 'medium'
        })
        console.log('Todo added successfully:', response.data)
        this.todos.push(response.data)
      } catch (error) {
        console.error('Error adding todo:', error)
      }
    },
    async toggleComplete(todo) {
      console.log('Toggling todo completion:', todo)
      try {
        await axios.put(`${API_URL}/todos/${todo.id}`, {
          ...todo,
          completed: !todo.completed
        })
        todo.completed = !todo.completed
        console.log('Todo completion toggled:', todo)
      } catch (error) {
        console.error('Error updating todo:', error)
      }
    },
    async updatePriority(todo, priority) {
      console.log('Updating todo priority:', todo, priority)
      try {
        await axios.put(`${API_URL}/todos/${todo.id}`, {
          ...todo,
          priority: priority
        })
        todo.priority = priority
        console.log('Todo priority updated:', todo)
      } catch (error) {
        console.error('Error updating todo priority:', error)
      }
    },
    async removeTodo(todo) {
      console.log('Removing todo:', todo)
      try {
        await axios.delete(`${API_URL}/todos/${todo.id}`)
        this.todos = this.todos.filter(t => t.id !== todo.id)
        console.log('Todo removed:', todo)
      } catch (error) {
        console.error('Error removing todo:', error)
      }
    },
    setFilter(filter) {
      this.filter = filter
      console.log('Filter set to:', filter)
    }
  },
  mounted() {
    console.log('TodoList component mounted')
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
