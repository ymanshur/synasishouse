-- name: CreateOrder :one
INSERT INTO orders (
    order_no,
    user_id,
    status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserOrder :one
SELECT * FROM orders
WHERE
    order_no = $1
    AND user_id = $2
LIMIT 1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;
