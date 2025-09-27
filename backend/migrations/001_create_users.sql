-- Initial database schema for ABC User Management
-- File: 001_initial_schema.sql

-- Create users table with all required fields
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    age INT NOT NULL CHECK (age >= 18),
    password VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin', 'moderator') DEFAULT 'user',
    active BOOLEAN DEFAULT TRUE,
    permissions JSON DEFAULT NULL,
    last_login_at TIMESTAMP NULL,
    login_attempts INT UNSIGNED DEFAULT 0,
    locked_until TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes for better performance
    -- INDEX idx_email (email),
    -- INDEX idx_role (role),
    -- INDEX idx_active (active),
    -- INDEX idx_deleted_at (deleted_at),
    -- INDEX idx_created_at (created_at),
    -- INDEX idx_users_search (name, email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create admins table for admin-specific data
CREATE TABLE IF NOT EXISTS admins (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    admin_level ENUM('super_admin', 'admin', 'moderator') DEFAULT 'admin',
    permissions JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_admin_level (admin_level),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default admin user (password will be hashed by application)
-- Using a placeholder password that should be changed on first login
INSERT INTO users (name, email, age, password, role, active) 
VALUES (
    'System Administrator', 
    'admin@abc.com', 
    30, 
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewWrG3k9mHkKM8wG', -- placeholder: 'admin123'
    'admin', 
    TRUE
) ON DUPLICATE KEY UPDATE 
    name = VALUES(name),
    role = VALUES(role),
    active = VALUES(active);

-- Create corresponding admin record
INSERT INTO admins (user_id, admin_level, permissions)
SELECT 
    u.id, 
    'super_admin',
    JSON_OBJECT(
        'users', JSON_ARRAY('create', 'read', 'update', 'delete'),
        'admins', JSON_ARRAY('create', 'read', 'update', 'delete'),
        'system', JSON_ARRAY('manage', 'configure')
    )
FROM users u 
WHERE u.email = 'admin@abc.com' 
ON DUPLICATE KEY UPDATE 
    admin_level = VALUES(admin_level),
    permissions = VALUES(permissions);