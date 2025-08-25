# ACG-FAKA v2.0 - Go + SQLite + Vue.js

This is a complete rewrite of the ACG-FAKA virtual card selling system using modern technologies:

- **Backend**: Go with Gin framework
- **Database**: SQLite for lightweight deployment
- **Frontend**: Vue.js 3 with Element Plus UI

## Features

- 🔐 JWT-based authentication system
- 📦 Product catalog and management
- 🛒 Order processing and management
- 👥 User and admin management
- 💳 Payment system integration ready
- 📱 Responsive mobile-friendly design
- 🚀 High performance with Go backend
- 📊 Admin dashboard with analytics

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Install Go dependencies:
```bash
go mod tidy
```

4. Run the backend server:
```bash
go run cmd/server/main.go
```

The backend will start on http://localhost:8080

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install npm dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will start on http://localhost:3000

## API Endpoints

### Public APIs
- `GET /api/v1/public/configs` - Get public configuration
- `GET /api/v1/public/categories` - Get product categories
- `GET /api/v1/public/commodities` - Get products with pagination
- `GET /api/v1/public/commodity/:id` - Get single product

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/admin/login` - Admin login

### User APIs (Requires Authentication)
- `GET /api/v1/user/profile` - Get user profile
- `PUT /api/v1/user/profile` - Update user profile
- `GET /api/v1/user/orders` - Get user orders
- `POST /api/v1/user/orders` - Create new order
- `GET /api/v1/user/balance` - Get user balance

### Admin APIs (Requires Admin Authentication)
- `GET /api/v1/admin/dashboard` - Dashboard statistics
- `GET /api/v1/admin/users` - Manage users
- `GET /api/v1/admin/orders` - Manage orders
- `GET /api/v1/admin/commodities` - Manage products
- `GET /api/v1/admin/categories` - Manage categories

## Database

The system uses SQLite with the following main tables:

- `users` - User accounts and profiles
- `categories` - Product categories
- `commodities` - Products/services
- `orders` - Customer orders
- `payments` - Payment methods
- `cards` - Virtual cards/keys
- `configs` - System configuration

The database is automatically initialized with demo data on first run.

## Configuration

### Backend Environment Variables

```env
ENVIRONMENT=development
PORT=8080
DATABASE_URL=./data/acg-faka.db
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

### Frontend Configuration

The frontend automatically proxies API requests to the backend during development. For production, update the `vite.config.js` file.

## Development

### Project Structure

```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── auth/           # Authentication logic
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   └── services/       # Business logic
├── migrations/         # Database migrations
└── data/              # SQLite database file

frontend/
├── src/
│   ├── api/           # API client
│   ├── components/    # Vue components
│   ├── router/        # Route configuration
│   ├── stores/        # Pinia state management
│   ├── views/         # Page components
│   └── utils/         # Utility functions
└── public/            # Static assets
```

### API Testing

You can test the API endpoints using curl:

```bash
# Get categories
curl http://localhost:8080/api/v1/public/categories

# Get products
curl http://localhost:8080/api/v1/public/commodities

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

## Production Deployment

### Backend

1. Build the binary:
```bash
cd backend
go build -o acg-faka cmd/server/main.go
```

2. Set production environment variables
3. Run the binary: `./acg-faka`

### Frontend

1. Build for production:
```bash
cd frontend
npm run build
```

2. Serve the `dist` folder with nginx or another web server
3. Configure proxy to backend API

## Migration from Original PHP Version

The new system maintains API compatibility where possible but uses a modern architecture:

- **Database**: All original table structures preserved in SQLite
- **Authentication**: Upgraded to JWT tokens instead of sessions
- **API**: RESTful design replacing PHP endpoints
- **Frontend**: Modern Vue.js SPA replacing jQuery/Bootstrap templates

## License

MIT License - Same as the original ACG-FAKA project.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Support

For issues and questions:
- Create an issue on GitHub
- Check the documentation
- Review the API endpoints above