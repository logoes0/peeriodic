-- Database setup script for Peeriodic
-- Run this script as a PostgreSQL superuser (like 'postgres')

-- Step 1: Create the database (run this as superuser)
-- Uncomment the line below if the database doesn't exist
-- CREATE DATABASE peeriodic;

-- Step 2: Create a user for the application (run this as superuser)
-- Uncomment and modify the lines below if you want to create a new user
-- CREATE USER logoes WITH PASSWORD 'your_password_here';
-- GRANT ALL PRIVILEGES ON DATABASE peeriodic TO logoes;

-- Step 3: Connect to the peeriodic database and run the following:
-- \c peeriodic;

-- Drop existing tables if they exist
DROP TABLE IF EXISTS rooms CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create rooms table with proper schema
CREATE TABLE rooms (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL DEFAULT 'Untitled Room',
    content TEXT DEFAULT '',
    user_uid VARCHAR(255) REFERENCES users(uid) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_rooms_user_uid ON rooms(user_uid);
CREATE INDEX idx_rooms_updated_at ON rooms(updated_at);
CREATE INDEX idx_users_uid ON users(uid);
CREATE INDEX idx_users_email ON users(email);

-- Optional: Create a function to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
DROP TRIGGER IF EXISTS update_rooms_updated_at ON rooms;
CREATE TRIGGER update_rooms_updated_at 
    BEFORE UPDATE ON rooms 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert test data
INSERT INTO users (uid, email, name) 
VALUES ('test-user', 'test@example.com', 'Test User')
ON CONFLICT (uid) DO NOTHING;

INSERT INTO rooms (id, title, content, user_uid) 
VALUES ('test-room-123', 'Test Room', 'Hello World!', 'test-user')
ON CONFLICT (id) DO NOTHING;

