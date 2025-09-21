SELECT
  status,
  COUNT(*) as order_count,
  SUM(amount) as total_amount
FROM orders
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY status
ORDER BY status;