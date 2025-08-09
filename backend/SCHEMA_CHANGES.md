# Database Schema Changes

This document outlines the changes made to the database schema to implement proper persistence and user management.

## Overview

The original schema only had a `rooms` table with basic fields. The new schema adds:

1. **Users table** - For user management and authentication
2. **Improved rooms table** - With proper foreign key relationships
3. **Better data integrity** - With constraints and triggers

## Schema Changes

### 1. New Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Purpose**: Store user information for authentication and room ownership.

### 2. Updated Rooms Table

```sql
CREATE TABLE rooms (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL DEFAULT 'Untitled Room',
    content TEXT DEFAULT '',
    user_uid VARCHAR(255) REFERENCES users(uid) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Changes**:
- `user_uid` now references the `users` table
- Added foreign key constraint with `ON DELETE SET NULL`
- Changed from `DEFAULT ''` to `NULL` for optional user association

### 3. Indexes

```sql
-- Users table indexes
CREATE INDEX idx_users_uid ON users(uid);
CREATE INDEX idx_users_email ON users(email);

-- Rooms table indexes (existing)
CREATE INDEX idx_rooms_user_uid ON rooms(user_uid);
CREATE INDEX idx_rooms_updated_at ON rooms(updated_at);
```

### 4. Triggers

Both tables now have automatic `updated_at` timestamp updates:

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Applied to both users and rooms tables
```

## Code Changes

### 1. Models (`backend/models/model.go`)

- Added `User` struct
- Updated `Room` struct to use `*string` for `UserUID` (nullable)
- Updated `NewRoom` function signature

### 2. Database Service (`backend/services/database.go`)

- Added `User` struct and related methods
- Added `CreateUser`, `GetUserByUID`, `EnsureUserExists` methods
- Updated `Room` struct to handle nullable `UserUID`
- Updated `CreateRoom` and `EnsureRoomExists` methods

### 3. Room Handler (`backend/handlers/room_handler.go`)

- Updated to handle optional user creation
- Added user management in room creation
- Updated response structures to handle nullable `user_uid`

## Migration Steps

### Option 1: Fresh Installation

1. Run the updated `setup.sql`:
   ```bash
   psql -d peeriodic -f backend/setup.sql
   ```

### Option 2: Migration from Existing Database

1. Run the migration script:
   ```bash
   psql -d peeriodic -f backend/migrations/001_add_users_table.sql
   ```

2. Update existing data:
   ```sql
   -- Convert empty user_uid strings to NULL
   UPDATE rooms SET user_uid = NULL WHERE user_uid = '';
   ```

## API Changes

### Room Creation

**Before**:
```json
{
  "title": "My Room",
  "uid": "user123"
}
```

**After**:
```json
{
  "title": "My Room",
  "uid": "user123",
  "email": "user@example.com",
  "name": "John Doe"
}
```

### Room Response

**Before**:
```json
{
  "id": "room123",
  "title": "My Room",
  "content": "Hello World"
}
```

**After**:
```json
{
  "id": "room123",
  "title": "My Room",
  "content": "Hello World",
  "user_uid": "user123"
}
```

## Benefits

1. **Proper User Management**: Users are now tracked and can be associated with rooms
2. **Data Integrity**: Foreign key constraints ensure referential integrity
3. **Scalability**: Better indexing for performance
4. **Flexibility**: Rooms can exist without users (anonymous rooms)
5. **Audit Trail**: Automatic timestamp updates for all records

## Backward Compatibility

- Existing rooms without users will have `user_uid` set to `NULL`
- API responses include `user_uid` when available
- All existing functionality continues to work

## Testing

After implementing these changes, test:

1. **Room Creation**: Create rooms with and without users
2. **User Management**: Create users and associate with rooms
3. **Data Persistence**: Verify content is saved and retrieved correctly
4. **WebSocket**: Ensure real-time collaboration still works
5. **API Endpoints**: Test all room and user endpoints


