-- name: InsertPatrol :one
INSERT INTO patrol_profile (
  user_id,
  officer_id,
  officer_name,
  status,
  street,
  city,
  state,
  latitude,
  longitude,
  created_at,
  updated_at
) VALUES (
  sqlc.arg('user_id'),
  sqlc.arg('officer_id'),
  sqlc.arg('officer_name'),
  sqlc.arg('status'),
  sqlc.narg('street'),
  sqlc.narg('city'),
  sqlc.narg('state'),
  sqlc.narg('latitude'),
  sqlc.narg('longitude'),
  COALESCE(sqlc.narg('created_at'), now()),
  COALESCE(sqlc.narg('updated_at'), now())
)
RETURNING user_id;

-- name: GetAllPatrols :many
SELECT * FROM patrol_profile;

-- name: GetPatrolByUserID :one
SELECT * FROM patrol_profile
WHERE user_id = sqlc.arg('user_id');

-- name: IgnoreNullUpdatePatrolByUserID :execrows
UPDATE patrol_profile
SET
  officer_id  = COALESCE(sqlc.narg('officer_id'), officer_id),
  officer_name = COALESCE(sqlc.narg('officer_name'), officer_name),
  status       = COALESCE(sqlc.narg('status'),       status),
  street       = COALESCE(sqlc.narg('street'),       street),
  city         = COALESCE(sqlc.narg('city'),         city),
  state        = COALESCE(sqlc.narg('state'),        state),
  latitude     = COALESCE(sqlc.narg('latitude'),     latitude),
  longitude    = COALESCE(sqlc.narg('longitude'),    longitude),
  updated_at   = NOW()          -- always refresh
WHERE user_id   = sqlc.arg('user_id');


-- name: UpdateAllPatrolByUserID :execrows
UPDATE patrol_profile
SET
  officer_id   = sqlc.arg('officer_id'),
  officer_name = sqlc.arg('officer_name'),
  status       = sqlc.arg('status'),
  street       = sqlc.narg('street'),
  city         = sqlc.narg('city'),
  state        = sqlc.narg('state'),
  latitude     = sqlc.narg('latitude'),
  longitude    = sqlc.narg('longitude'),
  updated_at   = NOW()          -- always refresh
WHERE user_id   = sqlc.arg('user_id');

-- name: DeletePatrolByUserID :execrows
DELETE FROM patrol_profile
WHERE user_id = sqlc.arg('user_id');

-- name: UpdatePatrolStatusByUserID :execrows
UPDATE patrol_profile
SET status = sqlc.arg('status')
WHERE user_id = sqlc.arg('user_id');

-- name: UpdatePatrolLocationByUserID :execrows
UPDATE patrol_profile
SET street = sqlc.arg('street'),
    city = sqlc.arg('city'),
    state = sqlc.arg('state'),
    latitude = sqlc.arg('latitude'),
    longitude = sqlc.arg('longitude')
WHERE user_id = sqlc.arg('user_id');