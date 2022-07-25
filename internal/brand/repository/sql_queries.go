package repository

const (
	createBrandQuery = `INSERT INTO brands (brand_name, logo, pickup_address) 
		VALUES ($1, COALESCE(NULLIF($2, ''), null), $3) 
		RETURNING brand_id, brand_name, logo, pickup_address, created_at, updated_at`

	findByIdQuery = `SELECT brand_id, email, brand_name, logo, pickup_address, created_at, updated_at FROM brands WHERE brand_id = $1`

	findAllQuery = `SELECT brand_id, email, brand_name, logo, pickup_address, created_at, updated_at FROM brands LIMIT $1 OFFSET $2`

	updateByIdQuery = `UPDATE brands SET brand_name = $2, logo = $3, pickup_address = $4 WHERE brand_id = $1
		RETURNING brand_id, brand_name, logo, pickup_address, created_at, updated_at`

	deleteByIdQuery = `DELETE FROM brands WHERE brand_id = $1`
)
