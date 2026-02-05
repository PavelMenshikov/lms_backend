-- +goose Up
INSERT INTO users (id, first_name, last_name, email, password_hash, role)
VALUES ('00000000-0000-0000-0000-000000000001', 'Root', 'Admin', 'admin@capedu.kz', '', 'admin') 
ON CONFLICT (email) DO NOTHING;