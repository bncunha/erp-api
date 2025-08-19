ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP NULL;

CREATE TABLE inventories (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NULL,
  tenant_id BIGINT NOT NULL,
  type VARCHAR(50) NOT NULL,
  CONSTRAINT Inventories_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT Inventories_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

ALTER TABLE inventories ADD COLUMN created_at TIMESTAMP DEFAULT NOW() NOT NULL;