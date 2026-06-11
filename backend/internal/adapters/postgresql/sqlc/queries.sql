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

-- name: GetCategories :many
SELECT * FROM categories ORDER BY name;

-- name: GetProducts :many
SELECT 
    p.id as product_id, p.category_id, p.name as product_name, p.description,
    s.id as sku_id, s.sku_code, s.price, s.stock, s.attributes
FROM products p
LEFT JOIN skus s ON p.id = s.product_id
WHERE p.deleted_at IS NULL AND s.deleted_at IS NULL;

-- name: GetProduct :many
SELECT 
    p.id as product_id, p.category_id, p.name as product_name, p.description,
    s.id as sku_id, s.sku_code, s.price, s.stock, s.attributes
FROM products p
LEFT JOIN skus s ON p.id = s.product_id
WHERE p.id = $1 AND p.deleted_at IS NULL AND s.deleted_at IS NULL;
-- name: UpdateOrderStatus :exec
UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2;

-- name: GetUserOrders :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetAllOrders :many
SELECT * FROM orders ORDER BY created_at DESC;

-- name: GetUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUserRole :exec
UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2;

-- name: CreateProduct :one
INSERT INTO products (category_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE products SET category_id = $1, name = $2, description = $3, updated_at = NOW()
WHERE id = $4 AND deleted_at IS NULL;

-- name: SoftDeleteProduct :exec
UPDATE products SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: CreateSku :one
INSERT INTO skus (product_id, sku_code, attributes, price, stock)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateSku :exec
UPDATE skus SET sku_code = $1, attributes = $2, price = $3, stock = $4, updated_at = NOW()
WHERE id = $5 AND deleted_at IS NULL;

-- name: SoftDeleteSku :exec
UPDATE skus SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1;
