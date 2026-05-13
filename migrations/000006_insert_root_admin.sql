-- +goose Up
INSERT INTO users (id, first_name, last_name, email, password_hash, role)
VALUES ('00000000-0000-0000-0000-000000000001', 'Root', 'Admin', 'admin@capedu.kz', '$2a$12$yfSm99Ns/GwILUa7x0o96OqbFCZBipCVDM.e/P8BFfP4ISZvA.sjG', 'admin')
ON CONFLICT (email) DO NOTHING;