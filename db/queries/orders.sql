-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, status)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: ListUserOrders :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :one
UPDATE orders SET status = $1 WHERE id = $2 RETURNING *;

-- name: CancelOrder :exec
UPDATE orders SET status = 'Cancelled' WHERE id = $1 AND status = 'Pending';
