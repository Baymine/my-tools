<template>
  <div class="login-form">
    <h2>{{ isRegister ? '注册' : '登录' }}</h2>
    <form @submit.prevent="submitForm">
      <div v-if="isRegister" class="form-group">
        <label for="username">用户名</label>
        <input
          type="text"
          id="username"
          v-model="form.username"
          required
          placeholder="请输入用户名"
        />
        <div v-if="errors.username" class="error">{{ errors.username }}</div>
      </div>
      <div class="form-group">
        <label for="email">邮箱</label>
        <input
          type="email"
          id="email"
          v-model="form.email"
          required
          placeholder="请输入邮箱"
        />
        <div v-if="errors.email" class="error">{{ errors.email }}</div>
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <input
          type="password"
          id="password"
          v-model="form.password"
          required
          placeholder="请输入密码"
        />
        <div v-if="errors.password" class="error">{{ errors.password }}</div>
      </div>
      <div class="form-actions">
        <button type="submit" class="btn-primary">{{ isRegister ? '注册' : '登录' }}</button>
        <button type="button" class="btn-link" @click="toggleForm">
          {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
        </button>
      </div>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
    </form>
  </div>
</template>

<script>
import axios from 'axios'

const API_URL = 'http://localhost:8081'

export default {
  name: 'LoginForm',
  data() {
    return {
      isRegister: false,
      form: {
        username: '',
        email: '',
        password: ''
      },
      errors: {},
      errorMessage: ''
    }
  },
  methods: {
    async submitForm() {
      this.errors = {}
      this.errorMessage = ''
      
      try {
        let response
        if (this.isRegister) {
          response = await axios.post(`${API_URL}/register`, {
            username: this.form.username,
            email: this.form.email,
            password: this.form.password
          })
        } else {
          response = await axios.post(`${API_URL}/login`, {
            email: this.form.email,
            password: this.form.password
          })
        }
        
        // 保存用户信息和令牌
        localStorage.setItem('user', JSON.stringify(response.data))
        localStorage.setItem('token', response.data.token)
        
        // 触发登录成功事件
        this.$emit('auth-success')
      } catch (error) {
        console.error('Authentication error:', error)
        
        if (error.response) {
          if (error.response.status === 400 && typeof error.response.data === 'object') {
            // 表单验证错误
            this.errors = error.response.data
          } else if (error.response.status === 409) {
            // 用户名或邮箱已存在
            this.errorMessage = error.response.data
          } else if (error.response.status === 401) {
            // 登录失败
            this.errorMessage = '邮箱或密码错误'
          } else {
            // 其他错误
            this.errorMessage = '认证失败，请稍后再试'
          }
        } else {
          this.errorMessage = '无法连接到服务器，请检查网络连接'
        }
      }
    },
    toggleForm() {
      this.isRegister = !this.isRegister
      this.errors = {}
      this.errorMessage = ''
    }
  }
}
</script>

<style scoped>
.login-form {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #4CAF50;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
}

.btn-primary {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.btn-primary:hover {
  background-color: #45a049;
}

.btn-link {
  background: none;
  border: none;
  color: #4CAF50;
  cursor: pointer;
  font-size: 14px;
  text-decoration: underline;
}

.error {
  color: #f44336;
  font-size: 14px;
  margin-top: 5px;
}

.error-message {
  color: #f44336;
  text-align: center;
  margin-top: 15px;
  padding: 10px;
  background-color: #ffebee;
  border-radius: 4px;
}
</style>
