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
    type = $3,
    note = $4,
    transaction_date = $5
WHERE id = $6
  AND user_id = $7
RETURNING
    id,
    user_id,
    category_id,
    amount_cents,
    type,
    note,
    transaction_date,
    created_at;


-- name: DeleteTransaction :execrows
DELETE FROM transactions
WHERE id = $1
  AND user_id = $2;

-- name: GetTransactionSummary :one
SELECT
    COALESCE(SUM(CASE WHEN type = 'income' THEN amount_cents ELSE 0 END), 0) AS total_income,
    COALESCE(SUM(CASE WHEN type = 'expense' THEN amount_cents ELSE 0 END), 0) AS total_expense
FROM transactions
WHERE user_id = $1;

-- name: GetMonthlySummary :many
SELECT
    DATE_TRUNC('month', transaction_date) AS month,
    SUM(CASE WHEN type = 'income' THEN amount_cents ELSE 0 END) AS total_income,
    SUM(CASE WHEN type = 'expense' THEN amount_cents ELSE 0 END) AS total_expense
FROM transactions
WHERE user_id = $1
GROUP BY month
ORDER BY month DESC;
