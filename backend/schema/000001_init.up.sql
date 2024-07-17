CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
  id            UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
  name          VARCHAR(255)      NOT NULL,
  username      VARCHAR(255)      NOT NULL UNIQUE,
  password_hash VARCHAR(255)      NOT NULL
);

CREATE TABLE IF NOT EXISTS users_sessions
(
  user_id       UUID          REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  refresh_token VARCHAR(255) UNIQUE,
  expires_at    TIMESTAMP 
);