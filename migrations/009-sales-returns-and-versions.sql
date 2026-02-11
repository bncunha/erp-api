CREATE TABLE sales_versions (
  id BIGSERIAL PRIMARY KEY,
  sales_id BIGINT NOT NULL,
  version INT NOT NULL,
  date TIMESTAMP NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT SalesVersions_sales_id_fkey FOREIGN KEY (sales_id) REFERENCES sales(id),
  CONSTRAINT SalesVersions_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id),
  CONSTRAINT SalesVersions_unique UNIQUE (sales_id, version)
);

ALTER TABLE sales
  ADD COLUMN last_version INT NOT NULL DEFAULT 1;

INSERT INTO sales_versions (sales_id, version, date, tenant_id)
SELECT s.id, 1, s.date, s.tenant_id
FROM sales s;

UPDATE sales SET last_version = 1;

ALTER TABLE sales_items
  ADD COLUMN sales_version_id BIGINT;

UPDATE sales_items si
SET sales_version_id = sv.id
FROM sales_versions sv
WHERE sv.sales_id = si.sales_id
  AND sv.version = 1
  AND sv.tenant_id = si.tenant_id;

ALTER TABLE sales_items
  ALTER COLUMN sales_version_id SET NOT NULL;

ALTER TABLE sales_items
  ADD CONSTRAINT SalesItems_sales_version_id_fkey FOREIGN KEY (sales_version_id) REFERENCES sales_versions(id);

ALTER TABLE sales_items
  DROP CONSTRAINT IF EXISTS SalesItems_sales_id_fkey;

ALTER TABLE payments
  ADD COLUMN sales_version_id BIGINT;

UPDATE payments p
SET sales_version_id = sv.id
FROM sales_versions sv
WHERE sv.sales_id = p.sales_id
  AND sv.version = 1
  AND sv.tenant_id = p.tenant_id;

ALTER TABLE payments
  ALTER COLUMN sales_version_id SET NOT NULL;

ALTER TABLE payments
  ADD CONSTRAINT Payments_sales_version_id_fkey FOREIGN KEY (sales_version_id) REFERENCES sales_versions(id);

ALTER TABLE payments
  DROP CONSTRAINT IF EXISTS Payments_sales_id_fkey;

ALTER TABLE inventory_transactions
  ADD COLUMN sales_version_id BIGINT NULL;

UPDATE inventory_transactions it
SET sales_version_id = sv.id
FROM sales_versions sv
WHERE sv.sales_id = it.sales_id
  AND sv.version = 1
  AND sv.tenant_id = it.tenant_id;

ALTER TABLE inventory_transactions
  ADD CONSTRAINT InventoryTransactions_sales_version_id_fkey FOREIGN KEY (sales_version_id) REFERENCES sales_versions(id);

CREATE TABLE sales_returns (
  id BIGSERIAL PRIMARY KEY,
  sales_id BIGINT NOT NULL,
  from_sales_version_id BIGINT NOT NULL,
  to_sales_version_id BIGINT NOT NULL,
  return_date TIMESTAMP NOT NULL,
  returner_name VARCHAR(255) NOT NULL,
  reason VARCHAR(2000) NOT NULL,
  created_by_user_id BIGINT NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT SalesReturns_sales_id_fkey FOREIGN KEY (sales_id) REFERENCES sales(id),
  CONSTRAINT SalesReturns_from_sales_version_id_fkey FOREIGN KEY (from_sales_version_id) REFERENCES sales_versions(id),
  CONSTRAINT SalesReturns_to_sales_version_id_fkey FOREIGN KEY (to_sales_version_id) REFERENCES sales_versions(id),
  CONSTRAINT SalesReturns_created_by_user_id_fkey FOREIGN KEY (created_by_user_id) REFERENCES users(id),
  CONSTRAINT SalesReturns_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id),
  CONSTRAINT SalesReturns_reason_length CHECK (char_length(reason) <= 2000)
);

CREATE TABLE sales_return_items (
  id BIGSERIAL PRIMARY KEY,
  sales_return_id BIGINT NOT NULL,
  sku_id BIGINT NOT NULL,
  quantity FLOAT NOT NULL,
  unit_price FLOAT NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT SalesReturnItems_sales_return_id_fkey FOREIGN KEY (sales_return_id) REFERENCES sales_returns(id),
  CONSTRAINT SalesReturnItems_sku_id_fkey FOREIGN KEY (sku_id) REFERENCES skus(id),
  CONSTRAINT SalesReturnItems_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id),
  CONSTRAINT SalesReturnItems_quantity_positive CHECK (quantity > 0)
);
