// Code generated by sqlc. DO NOT EDIT.
// source: entery_Q.sql

package db

import (
	"context"
)

const createEnteries = `-- name: CreateEnteries :one
INSERT INTO enteries (
  account_id,
  amount
) VALUES (
  $1, $2
) RETURNING id, account_id, amount, created_at
`

type CreateEnteriesParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEnteries(ctx context.Context, arg CreateEnteriesParams) (Entery, error) {
	row := q.db.QueryRowContext(ctx, createEnteries, arg.AccountID, arg.Amount)
	var i Entery
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getEnteries = `-- name: GetEnteries :one
SELECT id, account_id, amount, created_at FROM enteries
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEnteries(ctx context.Context, id int64) (Entery, error) {
	row := q.db.QueryRowContext(ctx, getEnteries, id)
	var i Entery
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listEnteries = `-- name: ListEnteries :many
SELECT id, account_id, amount, created_at FROM enteries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListEnteriesParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEnteries(ctx context.Context, arg ListEnteriesParams) ([]Entery, error) {
	rows, err := q.db.QueryContext(ctx, listEnteries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entery
	for rows.Next() {
		var i Entery
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
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