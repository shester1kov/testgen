-- Add back the role column
ALTER TABLE users ADD COLUMN role VARCHAR(50);

-- Migrate role_id back to role
UPDATE users SET role = (SELECT name FROM roles WHERE id = users.role_id);

-- Make role NOT NULL
ALTER TABLE users ALTER COLUMN role SET NOT NULL;

-- Add CHECK constraint back
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'teacher', 'student'));

-- Drop role_id column
ALTER TABLE users DROP COLUMN role_id;

-- Drop roles table
DROP TABLE IF EXISTS roles;
