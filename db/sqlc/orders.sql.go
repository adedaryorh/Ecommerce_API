// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: orders.sql

package db

import (
	"context"
)

const cancelOrder = `-- name: CancelOrder :exec
UPDATE orders SET status = 'Cancelled' WHERE id = $1 AND status = 'Pending'
`

func (q *Queries) CancelOrder(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, cancelOrder, id)
	return err
}

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, status)
VALUES ($1, $2, $3) RETURNING id, user_id, status, total_amount, created_at, updated_at
`

type CreateOrderParams struct {
	UserID      int64  `json:"user_id"`
	TotalAmount string `json:"total_amount"`
	Status      string `json:"status"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder, arg.UserID, arg.TotalAmount, arg.Status)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.TotalAmount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrderByID = `-- name: GetOrderByID :one
SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE id = $1
`

func (q *Queries) GetOrderByID(ctx context.Context, id int64) (Order, error) {
	row := q.db.QueryRowContext(ctx, getOrderByID, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.TotalAmount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUserOrders = `-- name: ListUserOrders :many
SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3
`

type ListUserOrdersParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUserOrders(ctx context.Context, arg ListUserOrdersParams) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, listUserOrders, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Status,
			&i.TotalAmount,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrderStatus = `-- name: UpdateOrderStatus :one
UPDATE orders SET status = $1 WHERE id = $2 RETURNING id, user_id, status, total_amount, created_at, updated_at
`

type UpdateOrderStatusParams struct {
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

func (q *Queries) UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrderStatus, arg.Status, arg.ID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.TotalAmount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
