package repository

const (
	createOrderQuery = `INSERT INTO orders (user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING order_id, user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address, created_at, updated_at`

	findByIdQuery = `SELECT order_id, user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address, created_at, updated_at FROM orders WHERE order_id = $1`

	findAllQuery = `SELECT order_id, user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address, created_at, updated_at FROM orders LIMIT $1 OFFSET $2`

	findByUserIdQuery = `SELECT order_id, user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address, created_at, updated_at FROM orders WHERE user_id = $1 LIMIT $2 OFFSET $3`

	findAllBySellerIdQuery = `SELECT order_id, user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address, created_at, updated_at FROM orders WHERE seller_id = $1 LIMIT $2 OFFSET $3`

	findAllByUserIdSellerIDQuery = `SELECT order_id, user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address, created_at, updated_at FROM orders WHERE user_id = $1 AND seller_id = $2 LIMIT $3 OFFSET $4`

	updateByIdQuery = `UPDATE orders SET user_id = $2, seller_id = $3, item = $4, quantity = $5, total_price = $6, status = $7, delivery_source_address = $8, delivery_destination_address = $9 WHERE order_id = $1
		RETURNING order_id, user_id, seller_id, item, quantity, total_price, status, delivery_source_address, delivery_destination_address, created_at, updated_at`

	deleteByIdQuery = `DELETE FROM orders WHERE order_id = $1`
)
