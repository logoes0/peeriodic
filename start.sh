#!/bin/bash

# Peeriodic - Startup Script
# This script helps you set up and start the application

set -e

echo "🚀 Starting Peeriodic setup..."

# Check if required tools are installed
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed. Please install Go 1.24+"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is required but not installed. Please install Node.js 18+"; exit 1; }
command -v psql >/dev/null 2>&1 || { echo "❌ PostgreSQL is required but not installed. Please install PostgreSQL"; exit 1; }

# Function to check if database exists
check_database() {
    psql -lqt | cut -d \| -f 1 | grep -qw peeriodic
}

# Function to check if .env file exists
check_env() {
    [ -f backend/.env ]
}

echo "📦 Installing dependencies..."

# Install backend dependencies
echo "Installing backend dependencies..."
cd backend
go mod tidy
cd ..

# Install frontend dependencies
echo "Installing frontend dependencies..."
cd frontend/client
npm install
cd ../..

# Setup environment variables
if ! check_env; then
    echo "🔧 Setting up environment variables..."
    cd backend
    cat > .env << EOF
DB_USER=postgres
DB_PASSWORD=
DB_NAME=peeriodic
DB_HOST=localhost
DB_PORT=5432
DB_SSLMODE=disable
PORT=5000
HOST=localhost
EOF
    echo "✅ Created .env file. Please edit backend/.env with your database credentials."
    cd ..
else
    echo "✅ .env file already exists."
fi

# Setup database
if ! check_database; then
    echo "🗄️ Setting up database..."
    createdb peeriodic 2>/dev/null || echo "⚠️ Database 'peeriodic' already exists or could not be created."
    psql -d peeriodic -f backend/setup.sql 2>/dev/null || echo "⚠️ Could not run setup script. Please run manually: psql -d peeriodic -f backend/setup.sql"
else
    echo "✅ Database 'peeriodic' already exists."
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit backend/.env with your database credentials"
echo "2. Run 'make dev' to start both backend and frontend"
echo "3. Or run 'make run-be' and 'make run-fe' in separate terminals"
echo ""
echo "The application will be available at:"
echo "- Frontend: http://localhost:3000"
echo "- Backend API: http://localhost:5000"
