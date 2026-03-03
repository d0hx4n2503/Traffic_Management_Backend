package repository

const (
	createNotificationQuery = `
	INSERT INTO notifications (
		id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	)VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	)
	RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	`

	updateNotificationQuery = `
		UPDATE notifications
		SET 
			code = COALESCE(NULLIF($1, ''), code),
			title = COALESCE(NULLIF($2, ''), title),
			content = COALESCE(NULLIF($3, ''), content),
			type = COALESCE(NULLIF($4, ''), type),
			target = COALESCE(NULLIF($5, ''), target),
			target_user = COALESCE(NULLIF($6, ''), target_user),
			status = COALESCE(NULLIF($7, ''), status),
			modifier_id = COALESCE($8, modifier_id),
			version = version + 1,
			updated_at = $9
		WHERE id = $10 AND active = true
		RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	`

	deleteNotificationQuery = `
	UPDATE notifications
	SET
		active = false,
		version = version + 1,
		modifier_id = $1,
		updated_at = $2
	WHERE id = $3 AND active = true
	RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	`

	getNotificationByIdQuery = `
	SELECT id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE id = $1 AND active = true
	`

	getTotalCount = `
	SELECT COUNT(id)
	FROM notifications
	WHERE active = true
	`

	searchByTitleQuery = `
	SELECT id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE active = true AND title ILIKE '%' || $1 || '%'
	ORDER BY created_at DESC
	OFFSET $2 LIMIT $3`

	findByTitleCount = `
	SELECT COUNT(*)
	FROM notifications
	WHERE title ILIKE '%' || $1 || '%' AND active = true
	`

	findByTitle = `
	SELECT *
	FROM notifications
	WHERE title = $1 AND active = true
	`

	getNotification = `
	SELECT id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
	FROM notifications
	WHERE active = true
	ORDER BY updated_at DESC, created_at DESC OFFSET $1 LIMIT $2
	`

	getNotificationsForUser = `
        SELECT id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
        FROM notifications
        WHERE active = true
          AND created_at > $1  -- sau thời điểm user tạo tài khoản
          AND (
            target = 'all' 
            OR (target = 'personal' AND target_user = $2)
          )
        ORDER BY created_at DESC
        OFFSET $3 LIMIT $4
    `

	getTotalNotificationsForUserCount = `
        SELECT COUNT(*)
        FROM notifications
        WHERE active = true
          AND created_at > $1
          AND (
            target = 'all' 
            OR (target = 'personal' AND target_user = $2)
          )
    `

	markNotificationAsReadQuery = `
        UPDATE notifications
        SET status = 'read',
            updated_at = NOW()
        WHERE id = $1 
          AND target = 'personal' 
          AND target_user = $2
          AND active = true
        RETURNING id, code, title, content, type, target, target_user, status, creator_id, created_at, updated_at, active
    `
)
