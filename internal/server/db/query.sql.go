// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"
)

const createURL = `-- name: CreateURL :exec
INSERT INTO urls(id, url) VALUES ($1, $2)
`

type CreateURLParams struct {
	ID  string
	Url string
}

func (q *Queries) CreateURL(ctx context.Context, arg CreateURLParams) error {
	_, err := q.db.Exec(ctx, createURL, arg.ID, arg.Url)
	return err
}

const getURLById = `-- name: GetURLById :one
SELECT url FROM urls WHERE id=$1
`

func (q *Queries) GetURLById(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRow(ctx, getURLById, id)
	var url string
	err := row.Scan(&url)
	return url, err
}
