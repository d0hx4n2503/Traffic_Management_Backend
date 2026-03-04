package repository

const (
	createGovAgencyQuery = `
	INSERT INTO gov_agencies (
		id, name, user_address, address, city, type, phone, email, status, 
		version, created_at, updated_at, active
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	) RETURNING 
		id, name, user_address, address, city, type, phone, email, status, 
		version, created_at, updated_at, active
	`

	updateGovAgencyQuery = `
	UPDATE gov_agencies
	SET
		name = COALESCE(NULLIF($1, ''), name),
		user_address = COALESCE(NULLIF($2, ''), user_address),
		address = COALESCE(NULLIF($3, ''), address),
		city = COALESCE(NULLIF($4, ''), city),
		type = COALESCE(NULLIF($5, ''), type),
		phone = COALESCE(NULLIF($6, ''), phone),
		email = COALESCE(NULLIF($7, ''), email),
		status = COALESCE(NULLIF($8, ''), status),
		version = version + 1,
		updated_at = $9
	WHERE id = $10
	RETURNING *
	`

	deleteGovAgencyQuery = `
	UPDATE gov_agencies
	SET
		active = false,
		version = version + 1,
		updated_at = $1
	WHERE id = $2
	RETURNING *
	`

	revokeGovAgencyQuery = `
	UPDATE gov_agencies
	SET
		status = 'revoked',
		active = false,
		version = version + 1,
		updated_at = $1
	WHERE id = $2
	RETURNING *
	`

	getGovAgencyQuery = `
	SELECT *
	FROM gov_agencies
	WHERE id = $1 AND active = true
	`

	getTotalGovAgencyCount = `
	SELECT COUNT(id)
	FROM gov_agencies
	WHERE active = true
	`

	searchGovAgencyByNameCount = `
	SELECT COUNT(*)
	FROM gov_agencies
	WHERE active = true
	AND name ILIKE '%' || $1 || '%'
	`

	searchGovAgencyByName = `
	SELECT * 
	FROM gov_agencies
	WHERE name ILIKE '%' || $1 || '%' AND active = true	
	ORDER BY name
	OFFSET $2 LIMIT $3
	`

	getAllGovAgency = `
	SELECT id, user_address, name, address, city, type, phone, email, status, 
		version, updated_at, created_at, active
	FROM gov_agencies
	WHERE active = true
	ORDER BY updated_at, created_at OFFSET $1 LIMIT $2
	`

	findGovAgencyByName = `
	SELECT *
	FROM gov_agencies
	WHERE name = $1 AND active = true
	`

	findAgencyByUserAddress = `
	SELECT id, user_address, email
	FROM gov_agencies
	WHERE user_address = $1 AND active = true 
	`
)
