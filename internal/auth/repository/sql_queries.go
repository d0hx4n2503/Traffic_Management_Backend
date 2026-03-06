package repository

const (
	createUserQuery = `
        INSERT INTO users (
            id, user_address, identity_no, 
            full_name, date_of_birth, gender, nationality, place_of_origin, place_of_residence,
            active, role, version, creator_id, modifier_id, 
            created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
        ) RETURNING id, identity_no, user_address, full_name, date_of_birth, gender, nationality, place_of_origin, place_of_residence, active, role, version, creator_id, modifier_id, created_at, updated_at`

	updateUserQuery = `
        UPDATE users 
        SET
            user_address = COALESCE(NULLIF($1, ''), user_address),
            identity_no = COALESCE(NULLIF($2, ''), identity_no),
            full_name = COALESCE(NULLIF($3, ''), full_name),
            date_of_birth= COALESCE(NULLIF($4, ''), date_of_birth),
            gender= COALESCE(NULLIF($5, ''), gender),
            nationality= COALESCE(NULLIF($6, ''), nationality),
            place_of_origin= COALESCE(NULLIF($7, ''), place_of_origin),
            place_of_residence= COALESCE(NULLIF($8, ''), place_of_residence),
            active = COALESCE($9, active),
            role = COALESCE(NULLIF($10, ''), role),
            creator_id = COALESCE($11, creator_id),
            modifier_id = COALESCE($12, modifier_id),
            version = version + 1,
            updated_at = now()
        WHERE id = $13 AND version = $14
        RETURNING id, identity_no, user_address, active, role, version, creator_id, modifier_id, created_at, updated_at`

	deleteUserQuery = `
        UPDATE users
        SET
            active = false,
            version = version + 1,
            modifier_id = $1,
            updated_at = $2
        WHERE id = $3 AND version = $4
        RETURNING id`

	getUserQuery = `
        SELECT *
        FROM users
        WHERE id = $1 AND active = true`

	getTotalCount = `
        SELECT COUNT(id) 
        FROM users 
        WHERE active = true 
        AND (identity_no ILIKE '%' || $1 || '%' OR role ILIKE '%' || $1 || '%')`

	findUsers = `
        SELECT *
        FROM users 
        WHERE active = true 
        AND (identity_no ILIKE '%' || $1 || '%' OR role ILIKE '%' || $1 || '%')
        ORDER BY identity_no, role
        OFFSET $2 LIMIT $3`

	getTotal = `
        SELECT COUNT(id) 
        FROM users 
        WHERE active = true`

	getUsers = `
        SELECT *
        FROM users 
        WHERE active = true 
        ORDER BY COALESCE(NULLIF($1, ''), identity_no), role
        OFFSET $2 LIMIT $3`

	findUserByIdentity = `
        SELECT *
        FROM users
        WHERE identity_no = $1 AND active = true`

	findUserByUserAddress = `
        SELECT *
        FROM users
        WHERE user_address = $1 AND active = true`

	getUserIdentityAndNameByAddress = `
        SELECT identity_no, full_name
        FROM users
        WHERE user_address = $1 AND active = true`

	checkUserAddressLinked = `
        SELECT EXISTS(
        SELECT 1 
        FROM users 
        WHERE identity_no = $1 
          AND TRIM(user_address) <> '' 
          AND active = true
        )`

	linkWalletAddressQuery = `
        UPDATE users 
        SET user_address = $1, updated_at = now(), version = version + 1
        WHERE identity_no = $2
        AND active = true
        AND (user_address IS NULL OR TRIM(user_address) = '')
        RETURNING id`

	syncDriverLicenseWalletByIdentityQuery = `
        UPDATE driver_licenses
        SET wallet_address = $1,
            updated_at = now(),
            version = version + 1
        WHERE identity_no = $2
          AND active = true`

	unlinkWalletAddressQuery = `
        UPDATE users 
        SET user_address = NULL, updated_at = now(), version = version + 1
        WHERE identity_no = $1 AND active = true`

	clearDriverLicenseWalletByIdentityQuery = `
        UPDATE driver_licenses
        SET wallet_address = NULL,
            updated_at = now(),
            version = version + 1
        WHERE identity_no = $1
          AND active = true`
)
