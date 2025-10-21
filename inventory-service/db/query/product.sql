-- name: CreateProduct :one
INSERT INTO products (
    code,
    total,
    reserved
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: UpdateProduct :one
UPDATE products
SET code = $1
WHERE id = $2
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: UpdateStock :one
UPDATE products
SET
    total = COALESCE(sqlc.narg(total), total),
    reserved = COALESCE(sqlc.narg(reserved), reserved)
WHERE code = sqlc.arg(code)
RETURNING *;
