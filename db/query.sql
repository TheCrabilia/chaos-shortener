-- name: GetURLById :one
SELECT url FROM urls WHERE id=$1;

-- name: CreateURL :exec
INSERT INTO urls(id, url) VALUES ($1, $2);
