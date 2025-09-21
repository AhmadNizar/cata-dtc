SELECT
  customer_id,
  COUNT(*) as total_orders,
  SUM(amount) as total_spend
FROM orders
GROUP BY customer_id
ORDER BY total_spend DESC
LIMIT 5;