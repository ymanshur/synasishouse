-- name: CreateOrder :one
INSERT INTO orders (
    order_no,
    user_id,
    status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CreateOrderDetail :one
INSERT INTO order_details (
    order_id,
    product_code,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;
