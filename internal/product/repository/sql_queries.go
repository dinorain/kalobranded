package repository

const (
	createProductQuery = `INSERT INTO products (name, description, price, brand_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING product_id, name, description, price, brand_id, created_at, updated_at`

	findByIdQuery = `SELECT product_id, name, description, price, brand_id, created_at, updated_at FROM products WHERE product_id = $1`

	findAllQuery = `SELECT product_id, name, description, price, brand_id, created_at, updated_at FROM products LIMIT $1 OFFSET $2`

	findAllByBrandIdQuery = `SELECT product_id, name, description, price, brand_id, created_at, updated_at FROM products WHERE brand_id = $1 LIMIT $2 OFFSET $3`

	updateByIdQuery = `UPDATE products SET name = $2, description = $3, price = $4, brand_id = $5 WHERE product_id = $1
		RETURNING product_id, name, description, price, brand_id, created_at, updated_at`

	deleteByIdQuery = `DELETE FROM products WHERE product_id = $1`
)
