package repository

const (
	createLicenseQuery = `
	INSERT INTO vehicle_registration (
		id, owner_id, brand, type_vehicle, vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,
		owner_name, seats, issue_date, issuer, registration_code, registration_date, expiry_date, registration_place, on_blockchain, blockchain_txhash, status, 
		version, creator_id, modifier_id, created_at, updated_at
	)VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
	)RETURNING 
		id, owner_id, brand, type_vehicle, vehicle_no, color_plate, chassis_no, engine_no, color_vehicle,
		owner_name,seats, issue_date,issuer, registration_code, registration_date, expiry_date, registration_place, on_blockchain, blockchain_txhash,
		status, version, creator_id, modifier_id, created_at, updated_at
	`

	updateLicenseQuery = `
    UPDATE vehicle_registration
    SET
        owner_id = COALESCE($1, owner_id),
        brand = COALESCE(NULLIF($2, ''), brand),
        type_vehicle = COALESCE(NULLIF($3, ''), type_vehicle),
        vehicle_no = COALESCE(NULLIF($4, ''), vehicle_no),
        color_plate = COALESCE(NULLIF($5, ''), color_plate),
        chassis_no = COALESCE(NULLIF($6, ''), chassis_no),
        engine_no = COALESCE(NULLIF($7, ''), engine_no),
        color_vehicle = COALESCE(NULLIF($8, ''), color_vehicle),
        owner_name = COALESCE(NULLIF($9, ''), owner_name),
        seats = COALESCE($10, seats),
        issue_date = COALESCE(NULLIF($11, '')::date, issue_date),
		issuer = COALESCE(NULLIF($12, ''), issuer),
        registration_code = COALESCE(NULLIF($13, ''), registration_code),
		registration_date = COALESCE(NULLIF($14, '')::date, registration_date),
        expiry_date = COALESCE(NULLIF($15, '')::date, expiry_date),
		registration_place = COALESCE(NULLIF($16, ''), registration_place),
        status = COALESCE(NULLIF($17, ''), status),
        modifier_id = COALESCE($18, modifier_id),
        version = version + 1,
        updated_at = now(),
        active = COALESCE($19, active)
    WHERE id = $20
    RETURNING *
	`

	updateBlockchainConfirmationQuery = `
    UPDATE vehicle_registration 
    SET
        blockchain_txhash = COALESCE(NULLIF($1, ''), blockchain_txhash),
        on_blockchain = $2,
        modifier_id = COALESCE($3, modifier_id),
        version = version + 1,
        updated_at = $4
    WHERE id = $5
    RETURNING *
    `

	deleteLicenseQuery = `
	UPDATE vehicle_registration
	SET
		active = false,
		version = version + 1,
		modifier_id = $1,
		updated_at = $2
	WHERE id = $3
	RETURNING *
	`

	getLicenseQuery = `
	SELECT
		vr.*,
		COALESCE(u.user_address, dl.wallet_address) AS user_address
	FROM vehicle_registration vr
	LEFT JOIN driver_licenses dl ON vr.owner_id = dl.id AND dl.active = true
	LEFT JOIN users u ON dl.identity_no = u.identity_no AND u.active = true
	WHERE vr.id = $1 AND vr.active = true
	`

	getTotalCount = `
	SELECT COUNT(id)
	FROM vehicle_registration
	WHERE active = true
	`

	findByVehiclePlateNOCount = `
		SELECT COUNT(*)
		FROM vehicle_registration
		WHERE active = true
		AND vehicle_no ILIKE '%' || $1 || '%'
	`

	searchByVehiclePlateNO = `
    SELECT * 
    FROM vehicle_registration
    WHERE vehicle_no ILIKE '%' || $1 || '%' AND active = true	
    ORDER BY vehicle_no
    OFFSET $2 LIMIT $3
	`

	getVehicleDocuments = `
    SELECT 
        vr.id, 
        vr.owner_id, 
        vr.brand, 
        vr.type_vehicle, 
        vr.vehicle_no, 
        vr.color_plate, 
        vr.chassis_no, 
        vr.engine_no, 
        vr.color_vehicle,
        COALESCE(dl.full_name, vr.owner_name) AS owner_name,
        vr.seats,
        vr.issue_date, 
        vr.expiry_date, 
        vr.issuer,
        vr.registration_date,
        vr.registration_place,
        vr.status, 
        vr.version, 
        vr.creator_id, 
        vr.modifier_id, 
        vr.updated_at, 
        vr.created_at,
        vr.on_blockchain,
        vr.blockchain_txhash,
		COALESCE(u.user_address, dl.wallet_address) AS user_address
    FROM vehicle_registration vr
    LEFT JOIN driver_licenses dl ON vr.owner_id = dl.id AND dl.active = true
	LEFT JOIN users u ON dl.identity_no = u.identity_no AND u.active = true
    WHERE vr.active = true
    ORDER BY vr.updated_at DESC, vr.created_at DESC
    OFFSET $1 LIMIT $2
    `

	findVehiclePlateNO = `
	SELECT *
	FROM vehicle_registration
	WHERE vehicle_no = $1 AND active = true
	`

	// Query for count by type_vehicle
	getCountByType = `
    SELECT type_vehicle, COUNT(*) as count
    FROM vehicle_registration
    WHERE active = true
    GROUP BY type_vehicle
    ORDER BY count DESC
    `

	getRegistrationStatusStats = `
    SELECT 
        COUNT(*) FILTER (WHERE expiry_date >= CURRENT_DATE AND expiry_date IS NOT NULL) AS valid_count,
        COUNT(*) FILTER (WHERE expiry_date < CURRENT_DATE AND expiry_date IS NOT NULL) AS expired_count,
        COUNT(*) FILTER (WHERE expiry_date IS NULL OR registration_date IS NULL) AS pending_count
    FROM vehicle_registration
    WHERE active = true
      AND type_vehicle NOT ILIKE ANY (ARRAY[
        '%xe máy%', '%xe mô tô%', '%xe gắn máy%', 
        '%xe đạp%', '%xe đạp điện%', '%xe máy điện%'
      ])
    `

	// Query for top 5 brands
	getTopBrands = `
    SELECT brand, COUNT(*) as count
    FROM vehicle_registration
    WHERE active = true
    GROUP BY brand
    ORDER BY count DESC
    LIMIT 5
    `

	// Total active vehicles for others calculation
	getTotalActiveVehicles = `
    SELECT COUNT(*)
    FROM vehicle_registration
    WHERE active = true
    `

	// User - Owner ID
	getVehiclesByOwnerID = `
    SELECT 
        vr.id,
        vr.owner_id,
        vr.brand,
        vr.type_vehicle,
        vr.vehicle_no,
        vr.color_plate,
        vr.chassis_no,
        vr.engine_no,
        vr.color_vehicle,
        COALESCE(dl.full_name, vr.owner_name) AS owner_name,
        vr.seats,
        vr.issue_date,
        vr.expiry_date,
        vr.issuer,
        vr.registration_code,
        vr.registration_date,
        vr.registration_place,
        vr.status,
        vr.version,
        vr.creator_id,
        vr.modifier_id,
        vr.updated_at,
        vr.created_at,
        vr.on_blockchain,
        vr.blockchain_txhash,
		COALESCE(u.user_address, dl.wallet_address) AS user_address,
        vr.active
    FROM vehicle_registration vr
    INNER JOIN driver_licenses dl ON vr.owner_id = dl.id AND dl.active = true
    INNER JOIN users u ON dl.identity_no = u.identity_no AND u.active = true
    WHERE u.id = $1
      AND vr.active = true
    ORDER BY vr.updated_at DESC, vr.created_at DESC
    OFFSET $2 LIMIT $3
    `

	getTotalCountByOwnerID = `
    SELECT COUNT(*)
    FROM vehicle_registration vr
    INNER JOIN driver_licenses dl ON vr.owner_id = dl.id AND dl.active = true
    INNER JOIN users u ON dl.identity_no = u.identity_no AND u.active = true
    WHERE u.id = $1 AND vr.active = true
    `

	getVehicleByIDAndOwner = `
        SELECT 
			vr.*,
			COALESCE(u.user_address, dl.wallet_address) AS user_address
        FROM vehicle_registration vr
		LEFT JOIN driver_licenses dl ON vr.owner_id = dl.id AND dl.active = true
		LEFT JOIN users u ON dl.identity_no = u.identity_no AND u.active = true
        WHERE vr.id = $1 AND vr.owner_id = $2 AND vr.active = true
    `

	getInspections = `
    SELECT 
        vr.id, 
        vr.owner_id, 
        vr.brand, 
        vr.type_vehicle, 
        vr.vehicle_no, 
        vr.color_plate, 
        vr.chassis_no, 
        vr.engine_no, 
        vr.color_vehicle,
        COALESCE(dl.full_name, vr.owner_name) AS owner_name,
        vr.seats,
        vr.issue_date, 
        vr.expiry_date, 
        vr.issuer,
        vr.registration_code,
        vr.registration_date,
        vr.registration_place,
        vr.status, 
        vr.version, 
        vr.creator_id, 
        vr.modifier_id, 
        vr.updated_at, 
        vr.created_at,
        vr.on_blockchain,
        vr.blockchain_txhash,
		COALESCE(u.user_address, dl.wallet_address) AS user_address
    FROM vehicle_registration vr
    LEFT JOIN driver_licenses dl ON vr.owner_id = dl.id AND dl.active = true
	LEFT JOIN users u ON dl.identity_no = u.identity_no AND u.active = true
    WHERE vr.registration_code IS NOT NULL AND vr.active = true
    ORDER BY vr.updated_at DESC, vr.created_at DESC
    OFFSET $1 LIMIT $2
    `

	getInspectionsCount = `
    SELECT COUNT(*)
    FROM vehicle_registration vr
    WHERE vr.registration_code IS NOT NULL AND vr.active = true
    `

	getByRegistrationCode = `
    SELECT 
        vr.id, 
        vr.owner_id, 
        vr.brand, 
        vr.type_vehicle, 
        vr.vehicle_no, 
        vr.color_plate, 
        vr.chassis_no, 
        vr.engine_no, 
        vr.color_vehicle,
        COALESCE(dl.full_name, vr.owner_name) AS owner_name,
        vr.seats,
        vr.issue_date, 
        vr.expiry_date, 
        vr.issuer,
        vr.registration_code,
        vr.registration_date,
        vr.registration_place,
        vr.status, 
        vr.version, 
        vr.creator_id, 
        vr.modifier_id, 
        vr.updated_at, 
        vr.created_at,
        vr.on_blockchain,
        vr.blockchain_txhash,
		COALESCE(u.user_address, dl.wallet_address) AS user_address
    FROM vehicle_registration vr
    LEFT JOIN driver_licenses dl ON vr.owner_id = dl.id AND dl.active = true
	LEFT JOIN users u ON dl.identity_no = u.identity_no AND u.active = true
    WHERE vr.registration_code = $1 AND vr.active = true
    `
)

var excludedVehicleTypes = []string{
	"%xe máy%",
	"%xe mô tô%",
	"%xe gắn máy%",
	"%xe đạp%",
	"%xe đạp điện%",
	"%xe máy điện%",
}
