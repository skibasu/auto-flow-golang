package repository

const GET_AUTH_BY_EMAIL = `
	SELECT 
		u.id,
		u.password,
		COALESCE(json_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '[]') AS roles
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	WHERE u.email = $1
	GROUP BY u.id, u.password
	`
const GET_ME = `
	SELECT 
		u.id,
		u.email,
		COALESCE(u.first_name, '') AS first_name,
		COALESCE(u.last_name, '') AS last_name,
		COALESCE(u.phone_number, '') AS phone_number,
		COALESCE(u.avatar, '') AS avatar,
		COALESCE(json_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '[]') AS roles
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	WHERE u.id = $1
	GROUP BY 
	u.id,
	u.email,
	u.first_name,
	u.last_name,
	u.phone_number,
	u.avatar;
	`

const CREATE_USER = `
WITH new_user AS (
	INSERT INTO users (
		email,
		first_name,
		last_name,
		phone_number,
		password
	)
	VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'admin'))
	RETURNING id, email, first_name, last_name, phone_number, avatar
),
insert_roles AS (
	INSERT INTO user_roles (user_id, role)
	SELECT id, unnest($6::user_role[])
	FROM new_user
)
SELECT 
	id,
	email,
	COALESCE(first_name, '') AS first_name,
	COALESCE(last_name, '') AS last_name,
	COALESCE(phone_number, '') AS phone_number,
	COALESCE(avatar, '') AS avatar
FROM new_user;
`

const GET_USERS_PART_1 = `SELECT 
		u.id,
		u.email,
		COALESCE(u.first_name, '') AS first_name,
		COALESCE(u.last_name, '') AS last_name,
		COALESCE(u.phone_number, '') AS phone_number,
		COALESCE(u.avatar, '') AS avatar,
		COALESCE(json_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '[]') AS roles
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	`

const GET_USERS_PART_2 = `
	GROUP BY 
		u.id,
		u.email,
		u.first_name,
		u.last_name,
		u.phone_number,
		u.avatar
	`
