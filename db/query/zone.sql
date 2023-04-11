-- CREATE TABLE "zones" (
--   "id" int PRIMARY KEY,
--   "name" varchar,
--   "boundary" "geography(Polygon, 4326)",
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

-- CREATE TABLE "brand_zones" (
--   "id" int PRIMARY KEY,
--   "brand_id" int NOT NULL,
--   "zone_id" int NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

-- name: CreateZone :one
INSERT INTO zones (
  name,
  boundary
) VALUES ($1, $2)
RETURNING *;

-- name: UpdateZone :one
UPDATE zones SET
  name = $2,
  boundary = $3
WHERE id = $1
RETURNING *;

-- name: DeleteZone :exec
DELETE FROM zones WHERE id = $1;

-- name: GetZone :one
SELECT * FROM zones WHERE id = $1;

-- name: ListZones :many
SELECT * FROM zones ORDER BY id LIMIT $1 OFFSET $2;

-- name: ListZonesByBrandID :many
SELECT * FROM zones WHERE id IN (SELECT zone_id FROM brand_zones WHERE brand_id = $1) ORDER BY id LIMIT $2 OFFSET $3;

-- name: CreateBrandZone :one
INSERT INTO brand_zones (
  brand_id,
  zone_id
) VALUES ($1, $2)
RETURNING *;

-- name: DeleteBrandZone :exec
DELETE FROM brand_zones WHERE brand_id = $1 AND zone_id = $2;

-- name: ListBrandZones :many
SELECT * FROM brand_zones WHERE brand_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountZones :one
SELECT COUNT(*) FROM zones;

-- name: CountZonesByBrandID :one
SELECT COUNT(*) FROM zones WHERE id IN (SELECT zone_id FROM brand_zones WHERE brand_id = $1);

-- name: SearchZones :many
SELECT * FROM zones
WHERE name ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: SearchZonesByBrandID :many
SELECT * FROM zones
WHERE id IN (SELECT zone_id FROM brand_zones WHERE brand_id = $2) AND name ILIKE $1 ORDER BY id LIMIT $3 OFFSET $4;

-- name: SearchZonesByGeoPoint :many
SELECT * FROM zones
WHERE ST_Intersects(boundary, ST_GeomFromText('POINT(' || $1 || ' ' || $2 || ')', 4326))
ORDER BY id LIMIT $3 OFFSET $4;

-- name: SearchZonesByGeoPointAndBrandID :many
SELECT * FROM zones
WHERE id IN (SELECT zone_id FROM brand_zones WHERE brand_id = $3) AND ST_Intersects(boundary, ST_GeomFromText('POINT(' || $1 || ' ' || $2 || ')', 4326))
ORDER BY id LIMIT $4 OFFSET $5;

-- name: SearchZonesByGeoPolygon :many
SELECT * FROM zones
WHERE ST_Intersects(boundary, ST_GeomFromText('POLYGON((' || $1 || '))', 4326))
ORDER BY id LIMIT $2 OFFSET $3;

-- name: SearchZonesByGeoPolygonAndBrandID :many
SELECT * FROM zones
WHERE id IN (SELECT zone_id FROM brand_zones WHERE brand_id = $2) AND ST_Intersects(boundary, ST_GeomFromText('POLYGON((' || $1 || '))', 4326))
ORDER BY id LIMIT $3 OFFSET $4;