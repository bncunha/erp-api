CREATE TABLE customers (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  phone_number VARCHAR(20) NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMP NULL,
  CONSTRAINT Customers_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

CREATE TABLE sales (
  id BIGSERIAL PRIMARY KEY,
  date TIMESTAMP NOT NULL,
  customer_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT Sales_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customers(id),
  CONSTRAINT Sales_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT Sales_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

CREATE TABLE sales_items (
  id BIGSERIAL PRIMARY KEY,
  quantity FLOAT NOT NULL,
  unit_price FLOAT NOT NULL,
  sku_id BIGINT NOT NULL,
  sales_id BIGINT NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT SalesItems_sales_id_fkey FOREIGN KEY (sales_id) REFERENCES sales(id),
  CONSTRAINT SalesItems_sku_id_fkey FOREIGN KEY (sku_id) REFERENCES skus(id),
  CONSTRAINT SalesItems_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

CREATE TABLE payments (
  id BIGSERIAL PRIMARY KEY,
  payment_type VARCHAR(50) NOT NULL,
  sales_id BIGINT NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT Payments_sales_id_fkey FOREIGN KEY (sales_id) REFERENCES sales(id),
  CONSTRAINT Payments_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

CREATE TABLE payment_dates (
  id BIGSERIAL PRIMARY KEY,
  due_date DATE NOT NULL,
  paid_date DATE NULL,
  installment_number INT NOT NULL,
  installment_value FLOAT NOT NULL,
  status VARCHAR(50) NOT NULL,
  payment_id BIGINT NOT NULL,
  tenant_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT PaymentDates_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES payments(id),
  CONSTRAINT PaymentDates_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

ALTER TABLE inventory_transactions ADD COLUMN sales_id BIGINT NULL;
ALTER TABLE inventory_transactions ADD CONSTRAINT InventoryTransactions_sales_id_fkey FOREIGN KEY (sales_id) REFERENCES sales(id);
ALTER TABLE skus ALTER COLUMN price SET NOT NULL;

ALTER TABLE sales ADD COLUMN code VARCHAR(50) NOT NULL DEFAULT CONCAT('V-', NOW());