CREATE TABLE user_tokens (
  id           BIGSERIAL PRIMARY KEY,
  uuid         TEXT NOT NULL UNIQUE,
  user_id      BIGINT NOT NULL REFERENCES users(id),
  tenant_id    BIGINT NULL REFERENCES companies(id),

  type         TEXT NOT NULL,  -- 'invite' ou 'reset_password'
  code_hash    TEXT NOT NULL,
  expires_at   TIMESTAMPTZ NOT NULL,
  used_at      TIMESTAMPTZ NULL,

  created_by   BIGINT NULL REFERENCES users(id), -- quem gerou (admin no invite, próprio user no reset, etc.)
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX user_tokens_uuid_idx ON user_tokens(uuid);
CREATE INDEX user_tokens_user_type_idx ON user_tokens(user_id, type);
CREATE UNIQUE INDEX user_tokens_code_hash_idx ON user_tokens(code_hash);


ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;
UPDATE users SET email = CONCAT('user', id, '@example.com') WHERE email IS NULL;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
CREATE UNIQUE INDEX usuarios_phonenumber_unique_notnull ON users (phone_number) WHERE email IS NOT NULL;
ALTER TABLE customers ADD CONSTRAINT customers_phonenumber_unique_key UNIQUE (phone_number);
