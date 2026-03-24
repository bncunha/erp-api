ALTER TABLE quotes
  ADD COLUMN IF NOT EXISTS discount_percentage FLOAT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS shipping_cost FLOAT NOT NULL DEFAULT 0;

ALTER TABLE quotes
  DROP CONSTRAINT IF EXISTS quotes_discount_percentage_check,
  ADD CONSTRAINT quotes_discount_percentage_check CHECK (discount_percentage >= 0 AND discount_percentage <= 100);

ALTER TABLE quotes
  DROP CONSTRAINT IF EXISTS quotes_shipping_cost_check,
  ADD CONSTRAINT quotes_shipping_cost_check CHECK (shipping_cost >= 0);