ALTER TABLE IF EXISTS payments RENAME TO sales_payments;
ALTER TABLE IF EXISTS payment_dates RENAME TO sales_payment_dates;

CREATE TABLE plans (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(30) UNIQUE NOT NULL,          -- TRIAL, BASIC, PREMIUM
  price DECIMAL(10,2) NOT NULL,              
  status VARCHAR(10) NOT NULL DEFAULT 'ACTIVE' -- ACTIVE, INACTIVE
);

CREATE TABLE subscriptions (
  id BIGSERIAL PRIMARY KEY,
  company_id BIGINT NOT NULL UNIQUE REFERENCES companies(id) ON DELETE CASCADE,
  plan_id BIGINT NOT NULL REFERENCES plans(id),

  status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',      -- ACTIVE, PAST_DUE, CANCELED
  current_period_end TIMESTAMPTZ NOT NULL,           -- validade do acesso

  -- Mercado Pago (futuro)
  provider_name VARCHAR(20),                             -- MERCADOPAGO
  provider_subscription_id VARCHAR(120),                   -- preapproval/assinatura
  provider_customer_id VARCHAR(120),

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE subscriptions_payments (
  id BIGSERIAL PRIMARY KEY,
  company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
  subscription_id BIGINT REFERENCES subscriptions(id) ON DELETE SET NULL,
  plan_id BIGINT NOT NULL REFERENCES plans(id),

  provider VARCHAR(20) NOT NULL DEFAULT 'MANUAL', -- MANUAL, MERCADOPAGO
  status VARCHAR(20) NOT NULL DEFAULT 'PAID',     -- PAID, PENDING, FAILED, REFUNDED

  amount DECIMAL(10,2) NOT NULL,
  paid_at TIMESTAMPTZ,

  -- Mercado Pago (futuro)
  provider_payment_id VARCHAR(120),
  provider_external_reference VARCHAR(120),

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_payments_company ON subscriptions_payments(company_id);

CREATE TABLE payment_events (
  id BIGSERIAL PRIMARY KEY,
  company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
  provider VARCHAR(20) NOT NULL, -- MERCADOPAGO
  provider_event_id VARCHAR(120) NOT NULL, -- id do evento no MP (ou assinatura do webhook)
  event_type VARCHAR(80) NOT NULL,
  payload JSONB NOT NULL,
  processed_at TIMESTAMPTZ,
  error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (provider, provider_event_id)
);

CREATE INDEX idx_payment_events_processed ON payment_events(processed_at);

INSERT INTO plans (name, price, status)
VALUES
  ('TRIAL', 0.00, 'ACTIVE'),
  ('PREMIUM', 79.90, 'ACTIVE');