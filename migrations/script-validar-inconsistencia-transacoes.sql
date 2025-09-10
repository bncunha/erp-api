-- Params:
--   :tenant_id      BIGINT (obrigatório)
--   :inventory_id   BIGINT (opcional; passe NULL para considerar todos)

WITH tx_by_sku AS (
  SELECT
    ii.sku_id,
    SUM(CASE t.type
          WHEN 'IN'        THEN t.quantity
          WHEN 'OUT'       THEN -t.quantity
          WHEN 'TRANSFER'  THEN 0
          ELSE 0
        END) AS tx_net,
    SUM(CASE WHEN t.type = 'IN'       THEN t.quantity ELSE 0 END) AS tx_in,
    SUM(CASE WHEN t.type = 'OUT'      THEN t.quantity ELSE 0 END) AS tx_out,
    SUM(CASE WHEN t.type = 'TRANSFER' THEN t.quantity ELSE 0 END) AS tx_transfer
  FROM inventory_transactions t
  JOIN inventory_items       ii ON ii.id = t.inventory_item_id
  JOIN skus                  s  ON s.id  = ii.sku_id
  WHERE s.tenant_id = 3
    --AND (NULL IS NULL OR ii.inventory_id = NULL)
  GROUP BY ii.sku_id
),
stock_by_sku AS (
  SELECT
    ii.sku_id,
    SUM(ii.quantity) AS stock_qty
  FROM inventory_items ii
  JOIN skus           s ON s.id = ii.sku_id
  WHERE s.tenant_id = 3
    --AND (:inventory_id IS NULL OR ii.inventory_id = :inventory_id)
    AND ii.deleted_at IS NULL
  GROUP BY ii.sku_id
)
SELECT
  s.id                    AS sku_id,
  s.code,
  COALESCE(tx.tx_in, 0)   AS tx_in,
  COALESCE(tx.tx_out, 0)  AS tx_out,
  COALESCE(tx.tx_transfer, 0) AS tx_transfer,
  COALESCE(tx.tx_net, 0)  AS tx_net,     -- IN - OUT (TRANSFER = 0)
  COALESCE(st.stock_qty, 0) AS stock_qty -- soma atual em inventory_items
FROM skus s
LEFT JOIN tx_by_sku   tx ON tx.sku_id = s.id
LEFT JOIN stock_by_sku st ON st.sku_id = s.id
WHERE s.tenant_id = 3 AND s.deleted_at IS NULL AND tx_net != stock_qty
ORDER BY s.id;



