-- name: CreatePost :one
INSERT INTO posts (
        id,
        feed_id,
        title,
        url,
        description,
        published_at,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (title, url) DO
UPDATE
SET title = $3,
    description = $5,
    updated_at = NOW()
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
from posts
    JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;