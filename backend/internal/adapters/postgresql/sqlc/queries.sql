-- name: CreateUser :one
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: AddToCart :exec
INSERT INTO cart_items (user_id, sku_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, sku_id) 
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = NOW();

-- name: GetCartItems :many
SELECT c.quantity, s.id as sku_id, s.price, s.stock, p.name 
FROM cart_items c
JOIN skus s ON c.sku_id = s.id
JOIN products p ON s.product_id = p.id
WHERE c.user_id = $1 AND s.deleted_at IS NULL;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE user_id = $1;

-- name: DecrementSkuStock :one
-- Атомарное списание остатков. Если вернет sql.ErrNoRows - значит товара не хватает (Race Condition предотвращен)
UPDATE skus 
SET stock = stock - $1, updated_at = NOW()
WHERE id = $2 AND stock >= $1 AND deleted_at IS NULL
RETURNING *;

-- name: CreateOrder :one
INSERT INTO orders (
    public_number, user_id, delivery_type, delivery_address, 
    pickup_point_id, total_amount, promo_code_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, sku_id, quantity, price_at_purchase)
VALUES ($1, $2, $3, $4);

-- name: SearchOrders :many
-- Быстрый поиск по последним 4 символам публичного номера
SELECT * FROM orders 
WHERE search_index = $1 
ORDER BY created_at DESC;