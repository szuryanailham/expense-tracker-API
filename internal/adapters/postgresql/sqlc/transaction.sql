-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id,
    category_id,
    amount_cents,
    note,
    transaction_date
)
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING
    id,
    user_id,
    category_id,
    amount_cents,
    note,
    transaction_date,
    created_at,
    updated_at;


-- name: ListTransactionsByUser :many
SELECT
    t.id,
    t.user_id,
    t.category_id,
    c.name AS category_name,
    c.type AS category_type,
    t.amount_cents,
    t.note,
    t.transaction_date,
    t.created_at
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
ORDER BY t.transaction_date DESC;

-- name: FindTransactionByID :one
SELECT
    t.id,
    t.user_id,
    t.category_id,
    c.name AS category_name,
    c.type AS category_type,
    t.amount_cents,
    t.note,
    t.transaction_date,
    t.created_at
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.id = $1
  AND t.user_id = $2;


-- name: UpdateTransaction :one
UPDATE transactions
SET
    category_id = $1,
    amount_cents = $2,
    note = $3,
    transaction_date = $4,
    updated_at = NOW()
WHERE id = $5
  AND user_id = $6
RETURNING
    id,
    user_id,
    category_id,
    amount_cents,
    note,
    transaction_date,
    created_at,
    updated_at;


-- name: DeleteTransaction :execrows
DELETE FROM transactions
WHERE id = $1
  AND user_id = $2;

-- name: GetTransactionSummary :one
SELECT
    COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount_cents ELSE 0 END), 0) AS total_income,
    COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount_cents ELSE 0 END), 0) AS total_expense
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1;

-- name: GetMonthlySummary :many
SELECT
    DATE_TRUNC('month', t.transaction_date)::DATE AS month,

    COALESCE(
        SUM(CASE WHEN c.type = 'income' THEN t.amount_cents ELSE 0 END),
        0
    ) AS total_income,

    COALESCE(
        SUM(CASE WHEN c.type = 'expense' THEN t.amount_cents ELSE 0 END),
        0
    ) AS total_expense
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
GROUP BY month
ORDER BY month;


