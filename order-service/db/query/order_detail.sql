-- name: CreateOrderDetail :one
INSERT INTO order_details (
    order_id,
    product_code,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: ListOrderDetails :many
SELECT * FROM order_details
WHERE order_id = $1;
