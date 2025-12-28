-- name: CreateCategory :one
INSERT INTO categories (
    user_id,
    name,
    type
)
VALUES (
    $1, $2, $3
)
RETURNING id, user_id, name, type, created_at;

-- name: ListCategoriesByUser :many
SELECT
    id,
    user_id,
    name,
    type,
    created_at
FROM categories
WHERE user_id IS NULL
   OR user_id = $1
ORDER BY name ASC;

-- name: FindCategoryByID :one
SELECT
    id,
    user_id,
    name,
    type,
    created_at
FROM categories
WHERE id = $1
  AND (user_id IS NULL OR user_id = $2);


-- name: UpdateCategory :one
UPDATE categories
SET
    name = $1,
    type = $2
WHERE id = $3
  AND user_id = $4
RETURNING id, user_id, name, type, created_at;

-- name: DeleteCategory :execrows
DELETE FROM categories
WHERE id = $1
  AND user_id = $2;
