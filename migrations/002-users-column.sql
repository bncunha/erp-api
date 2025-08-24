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
ALTER TABLE inventories ADD COLUMN deleted_at TIMESTAMP NULL;

CREATE TABLE inventory_items (
  id BIGSERIAL PRIMARY KEY,
  inventory_id BIGINT NOT NULL,
  sku_id BIGINT NOT NULL,
  quantity FLOAT NOT NULL,
  tenant_id BIGINT NOT NULL,
  CONSTRAINT InventoryItems_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES inventories(id),
  CONSTRAINT InventoryItems_sku_id_fkey FOREIGN KEY (sku_id) REFERENCES skus(id),
  CONSTRAINT InventoryItems_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

ALTER TABLE inventory_items ADD COLUMN created_at TIMESTAMP DEFAULT NOW() NOT NULL;
ALTER TABLE inventory_items ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE UNIQUE INDEX inventory_items_unique ON inventory_items (inventory_id, sku_id, tenant_id);

CREATE TABLE inventory_transactions (
  id BIGSERIAL PRIMARY KEY,
  quantity FLOAT NOT NULL,
  type VARCHAR(50) NOT NULL,
  date TIMESTAMP NOT NULL,
  inventory_in_id BIGINT,
  inventory_out_id BIGINT,
  inventory_item_id BIGINT NOT NULL,
  tenant_id BIGINT NOT NULL,
  justification VARCHAR(500) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMP NULL,
  CONSTRAINT InventoryTransactions_inventory_in_id_fkey FOREIGN KEY (inventory_in_id) REFERENCES inventories(id),
  CONSTRAINT InventoryTransactions_inventory_out_id_fkey FOREIGN KEY (inventory_out_id) REFERENCES inventories(id),
  CONSTRAINT InventoryTransactions_inventory_item_id_fkey FOREIGN KEY (inventory_item_id) REFERENCES inventory_items(id),
  CONSTRAINT InventoryTransactions_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);