-- Create roles table
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL CHECK (name IN ('admin', 'teacher', 'student')),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator with full system access'),
    ('teacher', 'Teacher who can create and manage tests'),
    ('student', 'Student who can take tests');

-- Add role_id column to users table
ALTER TABLE users ADD COLUMN role_id UUID REFERENCES roles(id) ON DELETE RESTRICT;

-- Migrate existing role data to role_id
UPDATE users SET role_id = (SELECT id FROM roles WHERE name = users.role);

-- Make role_id NOT NULL after migration
ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;

-- Drop the old role column
ALTER TABLE users DROP COLUMN role;

-- Add index for role_id
CREATE INDEX idx_users_role_id ON users(role_id);

-- Apply updated_at trigger to roles table
CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
