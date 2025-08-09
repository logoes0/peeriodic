# Database Setup Guide

This guide will help you set up PostgreSQL for the Peeriodic application.

## Prerequisites

- PostgreSQL installed and running
- Access to PostgreSQL as a superuser (usually `postgres`)

## Step 1: Check PostgreSQL Status

First, ensure PostgreSQL is running:

```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# If not running, start it
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

## Step 2: Access PostgreSQL as Superuser

Connect to PostgreSQL as the superuser:

```bash
# Method 1: Switch to postgres user
sudo -u postgres psql

# Method 2: Direct connection (if configured)
psql -U postgres -h localhost
```

## Step 3: Create Database and User

Once connected to PostgreSQL, run these commands:

```sql
-- Create the database
CREATE DATABASE peeriodic;

-- Create a user for the application
CREATE USER logoes WITH PASSWORD 'your_secure_password_here';

-- Grant privileges to the user
GRANT ALL PRIVILEGES ON DATABASE peeriodic TO logoes;

-- Connect to the peeriodic database
\c peeriodic;

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO logoes;

-- Exit PostgreSQL
\q
```

## Step 4: Run the Setup Script

Now you can run the setup script with your new user:

```bash
# Run the setup script
psql -U logoes -d peeriodic -f backend/setup.sql
```

If you get a password prompt, enter the password you set in Step 3.

## Alternative: Using Existing User

If you prefer to use an existing PostgreSQL user:

### Option A: Use Default PostgreSQL User

```bash
# Connect as postgres user and run setup
sudo -u postgres psql -d peeriodic -f backend/setup.sql
```

### Option B: Use Your System User

If your system user has PostgreSQL access:

```bash
# Create database first (as postgres user)
sudo -u postgres createdb peeriodic

# Then run setup as your user
psql -d peeriodic -f backend/setup.sql
```

## Step 5: Configure Environment Variables

Create a `.env` file in the `backend` directory:

```env
DB_USER=logoes
DB_PASSWORD=your_secure_password_here
DB_NAME=peeriodic
DB_HOST=localhost
DB_PORT=5432
DB_SSLMODE=disable
```

## Step 6: Test the Connection

Test if everything is working:

```bash
# Test connection
psql -U logoes -d peeriodic -c "SELECT version();"
```

## Troubleshooting

### Error: "role 'logoes' does not exist"

**Solution**: You need to create the user first:

```bash
# Connect as postgres superuser
sudo -u postgres psql

# Create the user
CREATE USER logoes WITH PASSWORD 'your_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE peeriodic TO logoes;

# Exit
\q
```

### Error: "database 'peeriodic' does not exist"

**Solution**: Create the database first:

```bash
# Connect as postgres superuser
sudo -u postgres psql

# Create the database
CREATE DATABASE peeriodic;

# Exit
\q
```

### Error: "permission denied"

**Solution**: Grant proper permissions:

```bash
# Connect as postgres superuser
sudo -u postgres psql

# Grant permissions
GRANT ALL PRIVILEGES ON DATABASE peeriodic TO logoes;
\c peeriodic
GRANT ALL ON SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO logoes;

# Exit
\q
```

### Error: "connection to server failed"

**Solution**: Check if PostgreSQL is running:

```bash
# Check status
sudo systemctl status postgresql

# Start if needed
sudo systemctl start postgresql
```

## Quick Setup Script

Here's a quick setup script you can run:

```bash
#!/bin/bash
# Quick PostgreSQL setup for Peeriodic

echo "Setting up PostgreSQL for Peeriodic..."

# Check if PostgreSQL is running
if ! systemctl is-active --quiet postgresql; then
    echo "Starting PostgreSQL..."
    sudo systemctl start postgresql
fi

# Create database and user
sudo -u postgres psql << EOF
CREATE DATABASE peeriodic;
CREATE USER logoes WITH PASSWORD 'peeriodic123';
GRANT ALL PRIVILEGES ON DATABASE peeriodic TO logoes;
\c peeriodic
GRANT ALL ON SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO logoes;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO logoes;
\q
EOF

echo "Database setup complete!"
echo "Now run: psql -U logoes -d peeriodic -f backend/setup.sql"
```

## Verification

After setup, verify everything is working:

```bash
# Test database connection
psql -U logoes -d peeriodic -c "SELECT COUNT(*) FROM users;"

# Should return: count
# ---------
#       1
```

## Next Steps

1. **Run the application**: `cd backend && go run main.go`
2. **Test the API**: Visit `http://localhost:5000/api/rooms`
3. **Check logs**: Look for "✅ Database connection established successfully"

## Security Notes

- Change the default password in production
- Use environment variables for sensitive data
- Consider using SSL in production
- Regularly backup your database
