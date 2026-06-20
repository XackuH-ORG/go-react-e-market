-- name: CreateUser :one
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: AddToCart :one
INSERT INTO cart_items (user_id, sku_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, sku_id) 
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = NOW()
RETURNING *;

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

-- name: CreateSession :one
INSERT INTO sessions (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByToken :one
SELECT s.id, s.user_id, s.token, s.expires_at, s.created_at, u.role
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token = $1 AND s.expires_at > NOW()
LIMIT 1;

-- name: GetUsers :many
SELECT id, email, role, created_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, email, role, updated_at;

-- name: CreateProduct :one
INSERT INTO products (name, description, image_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateSku :one
INSERT INTO skus (product_id, sku_code, price, stock)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListProducts :many
SELECT * FROM products
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetSkusByProductID :many
SELECT * FROM skus
WHERE product_id = $1 AND deleted_at IS NULL;
-- name: UpdateProduct :one
UPDATE products
SET name = COALESCE(NULLIF($2, ''), name),
    description = COALESCE(NULLIF($3, ''), description),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateProductImage :one
UPDATE products
SET image_url = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SearchProducts :many
SELECT DISTINCT p.* 
FROM products p
LEFT JOIN skus s ON p.id = s.product_id AND s.deleted_at IS NULL
WHERE p.deleted_at IS NULL
  AND (@search_query::text = '' OR p.name ILIKE '%' || @search_query || '%')
  AND (@min_price::int = 0 OR s.price >= @min_price)
  AND (@max_price::int = 0 OR s.price <= @max_price)
  AND (@in_stock::boolean = false OR s.stock > 0)
ORDER BY p.created_at DESC
LIMIT sqlc.arg('limit_val')::int OFFSET sqlc.arg('offset_val')::int;
