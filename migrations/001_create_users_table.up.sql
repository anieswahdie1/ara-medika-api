-- migrations/001_create_users_table.up.sql
CREATE TYPE role AS ENUM ('super_admin', 'admin', 'user');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role role NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Insert initial super admin
INSERT INTO users (name, email, password, role) 
VALUES ('Super Admin', 'superadmin@aramedika.com', '$2a$10$somehashedpassword', 'super_admin');