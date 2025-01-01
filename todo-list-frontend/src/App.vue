<template>
  <div id="app">
    <header>
      <h1>待办事项管理系统</h1>
    </header>
    <main>
      <TodoList />
    </main>
    <footer>
      <p>&copy; 2023 待办事项管理系统</p>
    </footer>
  </div>
</template>

<script>
import TodoList from './components/TodoList.vue'

export default {
  name: 'App',
  components: {
    TodoList
  }
}
</script>

<style>
body {
  font-family: 'Arial', sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #f4f4f4;
  margin: 0;
  padding: 0;
}

#app {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #fff;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

header {
  background-color: #4CAF50;
  color: white;
  text-align: center;
  padding: 1rem;
  margin-bottom: 2rem;
}

h1 {
  margin: 0;
}

main {
  padding: 0 1rem;
}

footer {
  text-align: center;
  margin-top: 2rem;
  padding: 1rem;
  background-color: #f8f8f8;
  border-top: 1px solid #e7e7e7;
}
</style>
