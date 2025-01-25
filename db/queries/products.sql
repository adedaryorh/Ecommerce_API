-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: UpdateProduct :one
UPDATE products
SET name = $1, description = $2, price = $3, stock = $4, updated_at = $5
WHERE id = $6 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY id LIMIT $1 OFFSET $2;
