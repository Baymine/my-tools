# Todo List Application

![Todo List App](https://via.placeholder.com/800x400?text=Todo+List+App)

## Project Overview

This is a powerful todo list application with a separate frontend and backend architecture, providing an intuitive user interface and a reliable backend API. The application allows users to easily manage their tasks, including adding new todos, updating task status, and setting priorities.

### Key Features

- 📝 Create and manage todo items
- ✅ Mark tasks as completed or incomplete
- 🔝 Set task priorities
- 🔄 Real-time synchronization between frontend and backend
- 📱 Responsive design for various devices

## Technology Stack

### Frontend
- Vue.js 3
- Axios for API requests
- CSS3 for styling

### Backend
- Go (Golang)
- Standard library `net/http` for HTTP server
- Custom middleware for CORS and request logging

## Quick Start

### Prerequisites

- Node.js (v14+)
- Go (v1.16+)

### Installation and Running

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/todo-list-app.git
   cd todo-list-app
   ```

2. Start the backend:
   ```
   cd todo-list-backend
   go run cmd/api/main.go
   ```

3. Start the frontend:
   ```
   cd todo-list-frontend
   npm install
   npm run serve
   ```

4. Open your browser and visit `http://localhost:8083`

## Project Structure

```
todo-list-app/
├── todo-list-backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── database/
│   │   ├── middleware/
│   │   ├── models/
│   │   └── validator/
│   └── test_api.sh
├── todo-list-frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── App.vue
│   │   └── main.js
│   └── package.json
└── Documentation/
    ├── project-status.md
    ├── DevelopmentGuidelines.md
    └── TroubleshootingGuide.md
```

## Contributing

Contributions are welcome! Please check out [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Contact

If you have any questions or suggestions, please contact the project maintainer:

- Email: kaihuacao04@gmail.com
- GitHub: [@Baymine](https://github.com/Baymine)

---

⭐️ If you like this project, please give it a star!
