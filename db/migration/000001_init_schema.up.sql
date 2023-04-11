CREATE TABLE "vendors" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "payment_info" jsonb NOT NULL,
  "social_links" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "vendor_permissions" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "vendor_user_permissions" (
  "id" bigserial PRIMARY KEY,
  "user_id" int NOT NULL,
  "permission_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "vendor_users" (
  "id" int PRIMARY KEY,
  "username" varchar,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "role" varchar NOT NULL,
  "vendor_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "operation_times" (
  "id" bigserial PRIMARY KEY,
  "brand_id" int NOT NULL,
  "day_of_week" int NOT NULL,
  "start_time" time NOT NULL,
  "end_time" time NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "zones" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "boundary" polygon NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "brand_zones" (
  "id" int PRIMARY KEY,
  "brand_id" int NOT NULL,
  "zone_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "brands" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "meta_tags" varchar NOT NULL,
  "slug" varchar NOT NULL,
  "type" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "email" varchar NOT NULL,
  "email_verified" bool NOT NULL,
  "logo" varchar NOT NULL,
  "banner" varchar NOT NULL,
  "rating" int NOT NULL DEFAULT 0,
  "vendor_id" bigint NOT NULL,
  "prefix" varchar NOT NULL,
  "status" varchar NOT NULL,
  "availability" bool NOT NULL,
  "location" point NOT NULL,
  "address" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "image" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "brand_categories" (
  "id" bigserial PRIMARY KEY,
  "brand_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "image" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "type" varchar NOT NULL,
  "name" varchar NOT NULL,
  "category_id" bigint NOT NULL,
  "slug" varchar NOT NULL,
  "image" varchar NOT NULL,
  "details" varchar NOT NULL,
  "price" float NOT NULL,
  "status" varchar NOT NULL,
  "brand_id" bigint NOT NULL,
  "availability" bool NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "inventory_history" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint NOT NULL,
  "quantity" int NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product_inventory" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint NOT NULL,
  "quantity" int NOT NULL,
  "purchase_price" float NOT NULL,
  "selling_price" float NOT NULL,
  "expire_date" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "product_variants" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint NOT NULL,
  "title" varchar NOT NULL,
  "min_select" int NOT NULL DEFAULT 1,
  "max_select" int NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "product_variant_items" (
  "id" bigserial PRIMARY KEY,
  "variant_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "vendors" ("name");

CREATE INDEX ON "vendor_user_permissions" ("user_id");

CREATE INDEX ON "vendor_user_permissions" ("permission_id");

CREATE INDEX ON "vendor_user_permissions" ("user_id", "permission_id");

CREATE INDEX ON "vendor_users" ("username");

CREATE INDEX ON "vendor_users" ("email");

CREATE UNIQUE INDEX ON "vendor_users" ("username", "email");

CREATE INDEX ON "operation_times" ("brand_id");

CREATE INDEX ON "operation_times" ("day_of_week");

CREATE INDEX ON "operation_times" ("brand_id", "day_of_week");

CREATE INDEX ON "brand_zones" ("brand_id");

CREATE INDEX ON "brand_zones" ("zone_id");

CREATE INDEX ON "brand_zones" ("brand_id", "zone_id");

CREATE INDEX ON "brands" ("slug");

CREATE INDEX ON "brands" ("email");

CREATE UNIQUE INDEX ON "brands" ("slug", "email");

CREATE INDEX ON "categories" ("name");

CREATE INDEX ON "brand_categories" ("brand_id");

CREATE INDEX ON "products" ("slug");

CREATE INDEX ON "products" ("brand_id");

CREATE INDEX ON "inventory_history" ("product_id");

CREATE INDEX ON "product_inventory" ("product_id");

CREATE INDEX ON "product_variants" ("product_id");

CREATE INDEX ON "product_variant_items" ("product_id");

CREATE INDEX ON "product_variant_items" ("variant_id");

CREATE INDEX ON "product_variant_items" ("product_id", "variant_id");

COMMENT ON COLUMN "inventory_history"."type" IS 'it would be like input, sold, restock, damage etc';

ALTER TABLE "vendor_user_permissions" ADD FOREIGN KEY ("user_id") REFERENCES "vendor_users" ("id");

ALTER TABLE "vendor_user_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "vendor_permissions" ("id");

ALTER TABLE "vendor_users" ADD FOREIGN KEY ("vendor_id") REFERENCES "vendors" ("id");

ALTER TABLE "operation_times" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "brand_zones" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "brand_zones" ADD FOREIGN KEY ("zone_id") REFERENCES "zones" ("id");

ALTER TABLE "brands" ADD FOREIGN KEY ("vendor_id") REFERENCES "vendors" ("id");

ALTER TABLE "brand_categories" ADD FOREIGN KEY ("brand_id") REFERENCES "categories" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "brand_categories" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "inventory_history" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_inventory" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_variants" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_variant_items" ADD FOREIGN KEY ("variant_id") REFERENCES "product_variants" ("id");

ALTER TABLE "product_variant_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
