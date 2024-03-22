-- name: AddMovie :one
INSERT INTO movies (
    title,
    description,
    rating,
    image
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListMovies :many
SELECT * FROM movies
LIMIT $1;

-- name: ListMoviesByTitle :many
SELECT * FROM movies
WHERE title = $1
LIMIT $2;

-- name: DetailMovie :one
SELECT * FROM movies
WHERE id = $1
LIMIT 1;

-- name: UpdateMovie :one
UPDATE movies
SET
    title = $1,
    description = $2,
    rating = $3,
    image = $4,
    updated_at = 'now()'
WHERE id = $5
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;