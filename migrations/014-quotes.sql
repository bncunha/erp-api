ALTER TABLE companies ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE companies ADD COLUMN IF NOT EXISTS logo_url VARCHAR(1024) NULL;

CREATE TABLE IF NOT EXISTS quote_number_sequences (
  tenant_id BIGINT PRIMARY KEY,
  last_number BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT QuoteNumberSequences_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id)
);

CREATE TABLE IF NOT EXISTS quotes (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  customer_id BIGINT NOT NULL,
  quote_number VARCHAR(20) NOT NULL,
  status VARCHAR(20) NOT NULL,
  valid_until DATE NOT NULL,
  down_payment_percentage FLOAT NOT NULL,
  notes TEXT NULL,
  shipping_type VARCHAR(20) NOT NULL,
  shipping_region VARCHAR(255) NULL,
  shipping_min_value FLOAT NULL,
  shipping_description VARCHAR(255) NOT NULL,
  subtotal_amount FLOAT NOT NULL DEFAULT 0,
  total_amount FLOAT NOT NULL DEFAULT 0,
  created_by_user_id BIGINT NOT NULL,
  production_order_id BIGINT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT Quotes_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES companies(id),
  CONSTRAINT Quotes_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customers(id),
  CONSTRAINT Quotes_created_by_user_id_fkey FOREIGN KEY (created_by_user_id) REFERENCES users(id),
  CONSTRAINT Quotes_down_payment_percentage_check CHECK (down_payment_percentage >= 0 AND down_payment_percentage <= 100),
  CONSTRAINT Quotes_shipping_min_value_check CHECK (shipping_min_value IS NULL OR shipping_min_value > 0),
  CONSTRAINT Quotes_unique_number_per_tenant UNIQUE (tenant_id, quote_number)
);

CREATE INDEX IF NOT EXISTS idx_quotes_tenant_id ON quotes (tenant_id);
CREATE INDEX IF NOT EXISTS idx_quotes_customer_id ON quotes (customer_id);
CREATE INDEX IF NOT EXISTS idx_quotes_status ON quotes (status);
CREATE INDEX IF NOT EXISTS idx_quotes_valid_until ON quotes (valid_until);
CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes (created_at);

CREATE TABLE IF NOT EXISTS quote_items (
  id BIGSERIAL PRIMARY KEY,
  quote_id BIGINT NOT NULL,
  sku_id BIGINT NOT NULL,
  product_description_snapshot TEXT NOT NULL,
  quantity FLOAT NOT NULL,
  unit_price FLOAT NOT NULL,
  total_price FLOAT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT QuoteItems_quote_id_fkey FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE,
  CONSTRAINT QuoteItems_sku_id_fkey FOREIGN KEY (sku_id) REFERENCES skus(id),
  CONSTRAINT QuoteItems_quantity_check CHECK (quantity > 0),
  CONSTRAINT QuoteItems_unit_price_check CHECK (unit_price >= 0),
  CONSTRAINT QuoteItems_total_price_check CHECK (total_price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_quote_items_quote_id ON quote_items (quote_id);
CREATE INDEX IF NOT EXISTS idx_quote_items_sku_id ON quote_items (sku_id);

