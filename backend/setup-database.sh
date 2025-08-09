#!/bin/bash

# Quick PostgreSQL setup script for Peeriodic
# Run this script as a regular user (it will use sudo when needed)

set -e

echo "🚀 Setting up PostgreSQL for Peeriodic..."

# Check if PostgreSQL is running
if ! systemctl is-active --quiet postgresql; then
    echo "📦 Starting PostgreSQL..."
    sudo systemctl start postgresql
    sudo systemctl enable postgresql
fi

echo "🔐 Creating database and user..."

# Create database and user
sudo -u postgres psql << EOF
-- Create database if it doesn't exist
SELECT 'CREATE DATABASE peeriodic' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'peeriodic')\gexec

-- Create user if it doesn't exist
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'logoes') THEN
        CREATE USER logoes WITH PASSWORD 'peeriodic123';
    END IF;
END
\$\$;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE peeriodic TO logoes;

-- Connect to peeriodic database
\c peeriodic

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO logoes;

-- Exit
\q
EOF

echo "✅ Database and user created successfully!"

echo "📝 Running setup script..."

# Run the setup script
sudo -u postgres psql -d peeriodic -f setup.sql

echo "🎉 Database setup complete!"
echo ""
echo "📋 Next steps:"
echo "1. Create a .env file in the backend directory with:"
echo "   DB_USER=logoes"
echo "   DB_PASSWORD=peeriodic123"
echo "   DB_NAME=peeriodic"
echo "   DB_HOST=localhost"
echo "   DB_PORT=5432"
echo "   DB_SSLMODE=disable"
echo ""
echo "2. Test the connection:"
echo "   psql -U logoes -d peeriodic -c \"SELECT COUNT(*) FROM users;\""
echo ""
echo "3. Run the application:"
echo "   cd backend && go run main.go"
echo ""
echo "🔒 Security note: Change the default password in production!"
