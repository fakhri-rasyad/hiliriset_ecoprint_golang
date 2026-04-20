ALTER TABLE users DROP CONSTRAINT user_public_id_unique;
DROP TABLE IF EXISTS users;

DROP TYPE users_role_enum;