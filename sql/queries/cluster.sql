-- name: CreateCluster :one
INSERT INTO
    clusters (name, image_url)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: DeleteCluster :exec
DELETE FROM clusters
WHERE
    id = $1;

-- name: GetCluster :one
SELECT
    *
FROM
    clusters
WHERE
    id = $1;

-- name: ListClusters :many
SELECT
    *
FROM
    clusters
ORDER BY
    id ASC;

-- name: PatchCluster :exec
UPDATE clusters
SET
    image_url = $2,
    name = $3
WHERE
    id = $1;
