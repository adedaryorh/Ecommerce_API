-- name: AddOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: RemoveOrderItem :exec
DELETE FROM order_items WHERE id = $1;
