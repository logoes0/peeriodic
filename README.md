# Peeriodic - Real-Time Collaborative Text Editor

A modern, real-time collaborative text editor built with Go (backend) and React/TypeScript (frontend). Multiple users can edit the same document simultaneously with live synchronization.

## 🚀 Features

- **Real-time collaboration**: Multiple users can edit simultaneously
- **Room-based system**: Each document is a "room" with unique ID
- **Live synchronization**: Changes appear instantly for all users
- **Auto-save**: Documents are automatically saved every 30 seconds
- **Shareable links**: Share room links with others to collaborate
- **Modern UI**: Clean, responsive interface with smooth animations
- **TypeScript**: Full type safety for better development experience

## 🏗️ Architecture

### Backend (Go)
- **Modular design**: Clean separation of concerns with services, handlers, and middleware
- **WebSocket support**: Real-time communication using Gorilla WebSocket
- **PostgreSQL**: Persistent storage for documents and room data
- **Configuration management**: Environment-based configuration
- **Graceful shutdown**: Proper server shutdown handling
- **Error handling**: Comprehensive error handling and logging

### Frontend (React/TypeScript)
- **Component-based**: Reusable, maintainable components
- **Type safety**: Full TypeScript support for better development
- **Service layer**: Centralized API and WebSocket services
- **State management**: React hooks for local state
- **Responsive design**: Mobile-friendly interface
- **Modern styling**: CSS with smooth animations and transitions

## 📁 Project Structure

```
peeriodic/
├── backend/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware (CORS, logging)
│   ├── models/          # Data models
│   ├── routers/         # Route definitions
│   ├── services/        # Business logic services
│   ├── utils/           # Utility functions
│   └── main.go          # Application entry point
├── frontend/
│   └── client/
│       ├── src/
│       │   ├── components/  # React components
│       │   ├── services/    # API and WebSocket services
│       │   ├── types/       # TypeScript type definitions
│       │   ├── utils/       # Utility functions
│       │   └── App.tsx      # Main application component
│       └── package.json
└── README.md
```

## 🛠️ Setup Instructions

### Prerequisites
- Go 1.24+ 
- Node.js 18+
- PostgreSQL 12+
- Git

### Quick Setup (Recommended)
```bash
# Run the automated setup script
./start.sh
```

This script will:
- Install all dependencies
- Create the database
- Set up environment variables
- Provide next steps

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd peeriodic
   ```

2. **Set up environment variables**
   ```bash
   cd backend
   # Create .env file with your database credentials
   echo "DB_USER=your_postgres_username" > .env
   echo "DB_PASSWORD=your_postgres_password" >> .env
   echo "DB_NAME=peeriodic" >> .env
   echo "DB_HOST=localhost" >> .env
   echo "DB_PORT=5432" >> .env
   echo "DB_SSLMODE=disable" >> .env
   echo "PORT=5000" >> .env
   echo "HOST=localhost" >> .env
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Set up database**
   ```bash
   # Create database
   createdb peeriodic
   
   # Run setup script
   psql -d peeriodic -f setup.sql
   ```

5. **Run the backend**
   ```bash
   go run main.go
   ```

### Frontend Setup

1. **Install dependencies**
   ```bash
   cd frontend/client
   npm install
   ```

2. **Start the development server**
   ```bash
   npm start
   ```

3. **Build for production**
   ```bash
   npm run build
   ```

### Using Makefile

```bash
# Start backend with live reload
make run-be

# Start frontend
make run-fe

# Tidy Go modules
make mod
```

## 🔧 Configuration

### Backend Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_USER` | Database username | Required |
| `DB_PASSWORD` | Database password | "" |
| `DB_NAME` | Database name | "peeriodic" |
| `DB_HOST` | Database host | "localhost" |
| `DB_PORT` | Database port | "5432" |
| `DB_SSLMODE` | SSL mode | "disable" |
| `PORT` | Server port | "5000" |
| `HOST` | Server host | "localhost" |

### Frontend Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `REACT_APP_API_URL` | Backend API URL | "http://localhost:5000" |

## 🚀 Usage

1. **Create a room**: Click "Create New Room" on the home page
2. **Share the room**: Click the share button to copy the room link
3. **Collaborate**: Multiple users can join via the shared link
4. **Real-time editing**: See changes as others type
5. **Auto-save**: Documents are saved automatically

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend/client
npm test
```

## 📝 API Documentation

### WebSocket Endpoints

- `GET /ws?room={roomId}` - Connect to a room for real-time collaboration

### HTTP Endpoints

- `GET /api/rooms?uid={userId}` - Get user's rooms
- `POST /api/rooms` - Create a new room
- `GET /api/rooms/{id}` - Get room details
- `DELETE /api/rooms/{id}` - Delete a room
- `POST /api/save?room={roomId}` - Save document content

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Gorilla WebSocket](https://github.com/gorilla/websocket) for WebSocket support
- [React](https://reactjs.org/) for the frontend framework
- [TypeScript](https://www.typescriptlang.org/) for type safety