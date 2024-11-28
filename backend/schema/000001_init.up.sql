CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
  id            UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
  email         CITEXT       NOT NULL UNIQUE,
  username      VARCHAR(255) NOT NULL UNIQUE DEFAULT substring(md5(random()::text), 0, 24),
  password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_sessions
(
  user_id       UUID          REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  refresh_token VARCHAR(255) UNIQUE,
  expires_at    TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS codes
(
  user_email         CITEXT       REFERENCES users (email) ON DELETE CASCADE NOT NULL,
  code               VARCHAR(255) PRIMARY KEY DEFAULT substring(md5(random()::text), 0, 32),
  code_type          VARCHAR(255) NOT NULL,
  expires_at         TIMESTAMP 
);