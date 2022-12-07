-- name: FindMemberByID :one
SELECT * FROM members
WHERE id = ? LIMIT 1;

-- name: FindMemberBySlackUID :one
SELECT * FROM members
WHERE slack_uid = ? LIMIT 1;

-- name: ListMembers :many
SELECT * FROM members
ORDER BY meals_cooked ASC, meals_eaten DESC;

-- name: CreateMember :one
INSERT INTO members (
    slack_uid, full_name, leader
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: UpdateMemberLeaderStatus :exec
UPDATE members
set leader = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateMemberMealsCooked :exec
UPDATE members
set meals_cooked = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateMemberMealsEaten :exec
UPDATE members
set meals_eaten = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteMember :exec
DELETE FROM members
WHERE id = ?;

-- name: FindMealByID :one
SELECT * FROM meals
WHERE id = ? LIMIT 1;

-- name: FindMealByDate :one
SELECT * FROM meals
WHERE year = ? AND month = ? AND day = ? LIMIT 1;

-- name: FindMealBySlackMessageID :one
SELECT * FROM meals
WHERE slack_message_id = ? LIMIT 1;

-- name: CreateMeal :one
INSERT INTO meals (
    cook_slack_uid, year, month, day
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateMealSlackMessageID :exec
UPDATE meals
set slack_message_id = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateMealDescription :exec
UPDATE meals
set description = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateMealSlackUID :exec
UPDATE meals
set cook_slack_uid = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteMeal :exec
DELETE FROM meals
WHERE id = ?;

-- name: CountMealsByDate :one
SELECT count(*) FROM meals WHERE year = ? AND month = ? AND day = ?;
