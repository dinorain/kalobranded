package repository

const (
	createSellerQuery = `INSERT INTO sellers (first_name, last_name, email, password, avatar, pickup_address) 
		VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), null), $6) 
		RETURNING seller_id, first_name, last_name, email, password, avatar, pickup_address, created_at, updated_at`

	findByEmailQuery = `SELECT seller_id, email, first_name, last_name, avatar, password, pickup_address, created_at, updated_at FROM sellers WHERE email = $1`

	findByIdQuery = `SELECT seller_id, email, first_name, last_name, avatar, password, pickup_address, created_at, updated_at FROM sellers WHERE seller_id = $1`

	findAllQuery = `SELECT seller_id, email, first_name, last_name, avatar, password, pickup_address, created_at, updated_at FROM sellers LIMIT $1 OFFSET $2`

	updateByIdQuery = `UPDATE sellers SET first_name = $2, last_name = $3, email = $4, password = $5, avatar = $6, pickup_address = $7 WHERE seller_id = $1
		RETURNING seller_id, first_name, last_name, email, password, avatar, pickup_address, created_at, updated_at`

	deleteByIdQuery = `DELETE FROM sellers WHERE seller_id = $1`
)
