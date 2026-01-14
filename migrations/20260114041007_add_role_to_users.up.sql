CREATE TYPE role AS ENUM ('admin', 'user');

ALTER TABLE users ADD COLUMN IF NOT EXISTS role role DEFAULT 'user';