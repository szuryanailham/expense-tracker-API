-- name: CreateUser :one
INSERT INTO users (
    first_name,
    last_name,
    email,
    password
)
VALUES (
    $1, $2, $3,$4
)
RETURNING id, first_name,last_name, email, created_at;


-- name: FindUserByEmail :one
SELECT
    id,
    first_name,
    last_name,
    email,
    password,
    created_at,
    updated_at
FROM users
WHERE email = $1;


-- name: FindUserByID :one
SELECT
    id,
    first_name,
    last_name,
    email,
    created_at,
    updated_at
FROM users
WHERE id = $1;
