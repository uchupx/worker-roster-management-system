package sqlite

var (
	insertUserQuery     = ` INSERT INTO users (name, email, password) VALUES (?, ?, ?);`
	insertUserRoleQuery = ` INSERT INTO user_roles (user_id, role_id) VALUES (?, ?);`

	findUserByEmailQuery = ` SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?;`
	findUserByIdQuery    = ` SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = ?;`
	findUserQuery =  `
	SELECT DISTINCT u.id, u.name, u.email, u.password, u.created_at, u.updated_at 
	FROM users as u
	LEFT JOIN user_roles as ur ON u.id = ur.user_id
	WHERE 
		(CASE WHEN ? IS NOT NULL THEN u.email = ? ELSE 1 END) AND 
		(CASE WHEN ? IS NOT NULL THEN ur.role_id = ? ELSE 1 END) AND 
		(CASE WHEN ? IS NOT NULL THEN u.id = ? ELSE 1 END);`	
)
