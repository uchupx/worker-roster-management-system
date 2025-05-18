package sqlite

var (
	CreateShiftQuery = `
		INSERT INTO 
			shifts (user_id, start_time, end_time, shift_date, status)
			VALUES 
				(?, ?, ?, ?, ?);`

	UpdateShiftQuery = `UPDATE shifts SET start_time = ?, end_time = ?, shift_date = ?, status = ? WHERE id = ?;`

	FindShiftByIDQuery = `
		SELECT 
			s.id, s.user_id, s.start_time, s.end_time, s.shift_date, s.status 
		FROM shifts s 
		WHERE s.id = ?;`

	FindShiftQuery = `
		SELECT 
				s.id, s.user_id, s.start_time, s.end_time, s.shift_date, s.status, u.name
			FROM shifts s
				JOIN users u ON s.user_id = u.id
			WHERE 
				(CASE WHEN ? IS NOT NULL THEN s.id = ? ELSE 1 END) AND 
				(CASE WHEN ? IS NOT NULL THEN s.status = ? ELSE 1 END) AND 
				(CASE WHEN ? IS NOT NULL THEN u.id = ? ELSE 1 END);`	

	FindShiftByDateQuery   = `
		SELECT 
				s.id, s.user_id, s.start_time, s.end_time, s.shift_date, s.status, u.name
			FROM shifts s
				JOIN users u ON s.user_id = u.id
		WHERE 
			shift_date BETWEEN ? AND ?;`

	FindShiftByUserIDDateQuery = `
			SELECT 
				s.id, s.user_id, s.start_time, s.end_time, s.shift_date, s.status, u.name
			FROM shifts s
				JOIN users u ON s.user_id = u.id
			WHERE 
				(CASE WHEN ? IS NOT NULL THEN s.user_id = ? ELSE 1 END)
				AND DATE(s.shift_date) BETWEEN DATE(?) AND DATE(?);
				s.status IN (0, 1)`
)
