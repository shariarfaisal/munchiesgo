CREATE TABLE "vendors" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar NOT NULL,
  "payment_info" jsonb NOT NULL,
  "social_links" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "vendor_users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "role" varchar NOT NULL,
  "vendor_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "operation_times" (
  "id" bigserial PRIMARY KEY,
  "brand_id" bigint NOT NULL,
  "day_of_week" int NOT NULL,
  "start_time" time NOT NULL,
  "end_time" time NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "zones" (
  "id" int PRIMARY KEY,
  "name" varchar NOT NULL,
  "boundary" Polygon NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "brand_zones" (
  "id" int PRIMARY KEY,
  "brand_id" bigint NOT NULL,
  "zone_id" bigint NOT NULL,
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
  "location" Point NOT NULL,
  "address" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "image" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "brand_categories" (
  "id" bigserial PRIMARY KEY,
  "brand_id" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "name" varchar NOT NULL,
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
  "use_inventory" bool NOT NULL DEFAULT false,
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

CREATE TABLE "customers" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "image" varchar NOT NULL,
  "email_verified" bool NOT NULL DEFAULT false,
  "nid" varchar NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "customer_addresses" (
  "id" bigserial PRIMARY KEY,
  "customer_id" bigint NOT NULL,
  "label" varchar NOT NULL,
  "address" varchar NOT NULL,
  "geo_point" Point NOT NULL,
  "apartment" varchar NOT NULL,
  "area" varchar NOT NULL,
  "floor" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "delivery_address" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "address" varchar NOT NULL,
  "geo_point" Point NOT NULL,
  "apartment" varchar NOT NULL,
  "area" varchar NOT NULL,
  "floor" varchar NOT NULL,
  "phone" varchar NOT NULL
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "customer_id" bigint NOT NULL,
  "status" varchar NOT NULL,
  "payment_method" varchar NOT NULL,
  "payment_status" varchar NOT NULL,
  "rider_note" varchar NOT NULL,
  "dispatch_time" timestamptz NOT NULL DEFAULT (now()),
  "total" float NOT NULL,
  "total_discount" float NOT NULL DEFAULT 0,
  "service_charge" float NOT NULL DEFAULT 0,
  "payable" float NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "brand_id" bigint NOT NULL,
  "price" float NOT NULL,
  "quantity" int NOT NULL,
  "discount" float NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "brand_orders" (
  "id" bigserial PRIMARY KEY,
  "brand_id" bigint NOT NULL,
  "order_id" bigint NOT NULL,
  "status" varchar NOT NULL,
  "total" float NOT NULL,
  "discount" float NOT NULL DEFAULT 0,
  "note" varchar NOT NULL DEFAULT ''
);

CREATE TABLE "rider_assign" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "rider_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "riders" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "hashed_password" varchar NOT NULL
);

CREATE TABLE "order_logs" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "user_type" varchar NOT NULL,
  "user_name" varchar NOT NULL,
  "user_id" bigint NOT NULL,
  "action_type" varchar NOT NULL,
  "prev_value" jsonb NOT NULL,
  "current_value" jsonb NOT NULL,
  "message" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "vendors" ("name");

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

CREATE INDEX ON "brand_categories" ("category_id");

CREATE INDEX ON "brand_categories" ("brand_id", "category_id");

CREATE INDEX ON "products" ("slug");

CREATE INDEX ON "products" ("brand_id");

CREATE INDEX ON "products" ("category_id");

CREATE INDEX ON "products" ("brand_id", "category_id");

CREATE INDEX ON "inventory_history" ("product_id");

CREATE INDEX ON "product_inventory" ("product_id");

CREATE INDEX ON "product_variants" ("product_id");

CREATE INDEX ON "product_variant_items" ("product_id");

CREATE INDEX ON "product_variant_items" ("variant_id");

CREATE INDEX ON "product_variant_items" ("product_id", "variant_id");

CREATE INDEX ON "customers" ("phone");

CREATE INDEX ON "customer_addresses" ("customer_id");

CREATE INDEX ON "delivery_address" ("order_id");

CREATE INDEX ON "delivery_address" ("customer_id");

CREATE INDEX ON "orders" ("customer_id");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "order_items" ("product_id");

CREATE INDEX ON "order_items" ("brand_id");

CREATE INDEX ON "brand_orders" ("brand_id");

CREATE INDEX ON "brand_orders" ("order_id");

COMMENT ON COLUMN "inventory_history"."type" IS 'it would be like input, sold, restock, damage etc';

ALTER TABLE "vendor_users" ADD FOREIGN KEY ("vendor_id") REFERENCES "vendors" ("id");

ALTER TABLE "operation_times" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "brand_zones" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "brand_zones" ADD FOREIGN KEY ("zone_id") REFERENCES "zones" ("id");

ALTER TABLE "brands" ADD FOREIGN KEY ("vendor_id") REFERENCES "vendors" ("id");

ALTER TABLE "brand_categories" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "brand_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "brand_categories" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "inventory_history" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_inventory" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_variants" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_variant_items" ADD FOREIGN KEY ("variant_id") REFERENCES "product_variants" ("id");

ALTER TABLE "product_variant_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "customer_addresses" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "delivery_address" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "delivery_address" ADD FOREIGN KEY ("customer_id") REFERENCES "orders" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "brand_orders" ADD FOREIGN KEY ("brand_id") REFERENCES "brands" ("id");

ALTER TABLE "brand_orders" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "rider_assign" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "rider_assign" ADD FOREIGN KEY ("rider_id") REFERENCES "riders" ("id");

ALTER TABLE "order_logs" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
