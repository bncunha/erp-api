-- Converte colunas sem timezone para TIMESTAMPTZ (fuso: America/Sao_Paulo)

-- customers
ALTER TABLE public.customers
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo',
  ALTER COLUMN deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'America/Sao_Paulo';

-- inventories
ALTER TABLE public.inventories
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo',
  ALTER COLUMN deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'America/Sao_Paulo';

-- inventory_items
ALTER TABLE public.inventory_items
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo',
  ALTER COLUMN deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'America/Sao_Paulo';

-- inventory_transactions
ALTER TABLE public.inventory_transactions
  ALTER COLUMN date TYPE timestamptz USING date AT TIME ZONE 'America/Sao_Paulo',
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo',
  ALTER COLUMN deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'America/Sao_Paulo';

-- payment_dates
ALTER TABLE public.payment_dates
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo';

-- payments
ALTER TABLE public.payments
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo';

-- sales
ALTER TABLE public.sales
  ALTER COLUMN date TYPE timestamptz USING date AT TIME ZONE 'America/Sao_Paulo',
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo';

-- sales_items
ALTER TABLE public.sales_items
  ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'America/Sao_Paulo';

-- users
ALTER TABLE public.users
  ALTER COLUMN deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'America/Sao_Paulo';

-- ✅ Opcional: se quiser garantir que futuros registros usem o fuso correto
ALTER DATABASE current_database() SET timezone TO 'America/Sao_Paulo';
