<template>
  <div>
    <LoginForm v-if="!isAuthenticated" @auth-success="onAuthSuccess" />
    <div v-else class="todo-list">
      <div class="user-info">
        <span>欢迎, {{ user.username }}</span>
        <button class="btn-logout" @click="logout">退出登录</button>
      </div>
      <TodoForm @add-todo="addTodo" />
      <TodoFilter @filter-change="setFilter" />
      
      <div v-if="loading" class="loading">加载中...</div>
      
      <ul v-else-if="sortedTodos.length > 0">
        <TodoItem
          v-for="todo in sortedTodos"
          :key="todo.id"
          :todo="todo"
          @toggle-complete="toggleComplete(todo)"
          @update-priority="updatePriority(todo, $event)"
          @remove="removeTodo(todo)"
        />
      </ul>
      
      <div v-else class="empty-state">
        没有待办事项，请添加一个新的待办事项
      </div>
      
      <Pagination
        v-if="usePagination && pagination.totalPages > 0"
        :current-page="pagination.page"
        :total-pages="pagination.totalPages"
        @page-change="changePage"
      />
      
      <div class="pagination-toggle">
        <label>
          <input type="checkbox" v-model="usePagination">
          启用分页 (每页 {{ pageSize }} 项)
        </label>
        <select v-if="usePagination" v-model="pageSize" @change="changePageSize">
          <option value="5">5</option>
          <option value="10">10</option>
          <option value="20">20</option>
          <option value="50">50</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script>
import TodoItem from './TodoItem.vue'
import TodoForm from './TodoForm.vue'
import TodoFilter from './TodoFilter.vue'
import LoginForm from './LoginForm.vue'
import Pagination from './Pagination.vue'
import axios from 'axios'

const API_URL = 'http://localhost:8081'

export default {
  name: 'TodoList',
  components: {
    TodoItem,
    TodoForm,
    TodoFilter,
    LoginForm,
    Pagination
  },
  data() {
    return {
      todos: [],
      filter: 'all',
      isAuthenticated: false,
      user: null,
      loading: false,
      usePagination: false,
      pageSize: 10,
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0,
        totalPages: 0
      }
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
    checkAuth() {
      const token = localStorage.getItem('token')
      const userStr = localStorage.getItem('user')
      
      if (token && userStr) {
        try {
          this.user = JSON.parse(userStr)
          this.isAuthenticated = true
          this.setupAxiosInterceptors()
          this.fetchTodos()
        } catch (e) {
          console.error('Error parsing user data:', e)
          this.logout()
        }
      }
    },
    setupAxiosInterceptors() {
      // 添加请求拦截器，自动添加认证头
      axios.interceptors.request.use(
        config => {
          const token = localStorage.getItem('token')
          if (token) {
            config.headers.Authorization = `Bearer ${token}`
          }
          return config
        },
        error => {
          return Promise.reject(error)
        }
      )
      
      // 添加响应拦截器，处理认证错误
      axios.interceptors.response.use(
        response => response,
        error => {
          if (error.response && error.response.status === 401) {
            // 认证失败，清除用户信息并重定向到登录页
            this.logout()
          }
          return Promise.reject(error)
        }
      )
    },
    onAuthSuccess() {
      this.checkAuth()
    },
    logout() {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      this.isAuthenticated = false
      this.user = null
      this.todos = []
    },
    async fetchTodos() {
      this.loading = true
      console.log('Fetching todos...')
      try {
        let url = `${API_URL}/todos`
        
        if (this.usePagination) {
          url = `${url}?page=${this.pagination.page}&pageSize=${this.pageSize}`
        }
        
        const response = await axios.get(url)
        console.log('Todos fetched successfully:', response.data)
        
        if (this.usePagination) {
          this.todos = response.data.todos
          this.pagination = response.data.pagination
        } else {
          this.todos = response.data
        }
        
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
      } finally {
        this.loading = false
      }
    },
    changePage(page) {
      this.pagination.page = page
      this.fetchTodos()
    },
    changePageSize() {
      this.pagination.page = 1
      this.pagination.pageSize = this.pageSize
      this.fetchTodos()
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
        // 重新获取待办事项列表以获取完整的待办事项对象
        this.fetchTodos()
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
    this.checkAuth()
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

.user-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.btn-logout {
  background-color: #f44336;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}

.btn-logout:hover {
  background-color: #d32f2f;
}

ul {
  list-style-type: none;
  padding: 0;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 30px;
  color: #666;
  background-color: #f9f9f9;
  border-radius: 4px;
  margin: 20px 0;
}

.pagination-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 20px;
  font-size: 14px;
  color: #666;
}

.pagination-toggle select {
  margin-left: 10px;
  padding: 3px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.pagination-toggle input[type="checkbox"] {
  margin-right: 5px;
}
</style>



