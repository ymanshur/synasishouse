-- name: CreateProduct :one
INSERT INTO products (
    code,
    total,
    hold
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: UpdateProduct :one
UPDATE products
SET
    code = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: AddStock :one
UPDATE products
SET
    total = total + sqlc.arg(total),
    hold = hold + sqlc.arg(hold),
    updated_at = NOW()
WHERE code = sqlc.arg(code)
RETURNING *;
