DROP INDEX IF EXISTS "vendors_name_idx";
DROP INDEX IF EXISTS "vendor_users_username_idx";
DROP INDEX IF EXISTS "vendor_users_email_idx";
DROP INDEX IF EXISTS "vendor_users_username_email_unique_idx";
DROP INDEX IF EXISTS "operation_times_brand_id_idx";
DROP INDEX IF EXISTS "operation_times_day_of_week_idx";
DROP INDEX IF EXISTS "operation_times_brand_id_day_of_week_idx";
DROP INDEX IF EXISTS "brand_zones_brand_id_idx";
DROP INDEX IF EXISTS "brand_zones_zone_id_idx";
DROP INDEX IF EXISTS "brand_zones_brand_id_zone_id_idx";
DROP INDEX IF EXISTS "brands_slug_idx";
DROP INDEX IF EXISTS "brands_email_idx";
DROP INDEX IF EXISTS "brands_slug_email_unique_idx";
DROP INDEX IF EXISTS "categories_name_idx";
DROP INDEX IF EXISTS "brand_categories_brand_id_idx";
DROP INDEX IF EXISTS "brand_categories_category_id_idx";
DROP INDEX IF EXISTS "brand_categories_brand_id_category_id_idx";
DROP INDEX IF EXISTS "products_slug_idx";
DROP INDEX IF EXISTS "products_brand_id_idx";
DROP INDEX IF EXISTS "products_category_id_idx";
DROP INDEX IF EXISTS "products_brand_id_category_id_idx";
DROP INDEX IF EXISTS "inventory_history_product_id_idx";
DROP INDEX IF EXISTS "product_inventory_product_id_idx";
DROP INDEX IF EXISTS "product_variants_product_id_idx";
DROP INDEX IF EXISTS "product_variant_items_product_id_idx";
DROP INDEX IF EXISTS "product_variant_items_variant_id_idx";
DROP INDEX IF EXISTS "product_variant_items_product_id_variant_id_idx";

ALTER TABLE "vendor_users" DROP CONSTRAINT IF EXISTS "vendor_users_vendor_id_fkey";
ALTER TABLE "operation_times" DROP CONSTRAINT IF EXISTS "operation_times_brand_id_fkey";
ALTER TABLE "brand_zones" DROP CONSTRAINT IF EXISTS "brand_zones_brand_id_fkey";
ALTER TABLE "brand_zones" DROP CONSTRAINT IF EXISTS "brand_zones_zone_id_fkey";
ALTER TABLE "brands" DROP CONSTRAINT IF EXISTS "brands_vendor_id_fkey";
ALTER TABLE "brand_categories" DROP CONSTRAINT IF EXISTS "brand_categories_brand_id_fkey";
ALTER TABLE "brand_categories" DROP CONSTRAINT IF EXISTS "brand_categories_category_id_fkey";
ALTER TABLE "products" DROP CONSTRAINT IF EXISTS "products_category_id_fkey";
ALTER TABLE "products" DROP CONSTRAINT IF EXISTS "products_brand_id_fkey";
ALTER TABLE "inventory_history" DROP CONSTRAINT IF EXISTS "inventory_history_product_id_fkey";
ALTER TABLE "product_inventory" DROP CONSTRAINT IF EXISTS "product_inventory_product_id_fkey";
ALTER TABLE "product_variants" DROP CONSTRAINT IF EXISTS "product_variants_product_id_fkey";
ALTER TABLE "product_variant_items" DROP CONSTRAINT IF EXISTS "product_variant_items_variant_id_fkey";
ALTER TABLE "product_variant_items" DROP CONSTRAINT IF EXISTS "product_variant_items_product_id_fkey";


DROP INDEX IF EXISTS "vendors_name_idx";

DROP INDEX IF EXISTS "vendor_users_username_idx";

DROP INDEX IF EXISTS "vendor_users_email_idx";

DROP INDEX IF EXISTS "vendor_users_username_email_idx";

DROP INDEX IF EXISTS "operation_times_brand_id_idx";

DROP INDEX IF EXISTS "operation_times_day_of_week_idx";

DROP INDEX IF EXISTS "operation_times_brand_id_day_of_week_idx";

DROP INDEX IF EXISTS "brand_zones_brand_id_idx";

DROP INDEX IF EXISTS "brand_zones_zone_id_idx";

DROP INDEX IF EXISTS "brand_zones_brand_id_zone_id_idx";

DROP INDEX IF EXISTS "brands_slug_idx";

DROP INDEX IF EXISTS "brands_email_idx";

DROP INDEX IF EXISTS "brands_slug_email_idx";

DROP INDEX IF EXISTS "categories_name_idx";

DROP INDEX IF EXISTS "brand_categories_brand_id_idx";

DROP INDEX IF EXISTS "brand_categories_category_id_idx";

DROP INDEX IF EXISTS "brand_categories_brand_id_category_id_idx";

DROP INDEX IF EXISTS "products_slug_idx";

DROP INDEX IF EXISTS "products_brand_id_idx";

DROP INDEX IF EXISTS "products_category_id_idx";

DROP INDEX IF EXISTS "products_brand_id_category_id_idx";

DROP INDEX IF EXISTS "inventory_history_product_id_idx";

DROP INDEX IF EXISTS "product_inventory_product_id_idx";

DROP INDEX IF EXISTS "product_variants_product_id_idx";

DROP INDEX IF EXISTS "product_variant_items_product_id_idx";

DROP INDEX IF EXISTS "product_variant_items_variant_id_idx";

DROP INDEX IF EXISTS "product_variant_items_product_id_variant_id_idx";

ALTER TABLE "vendor_users" DROP CONSTRAINT IF EXISTS "vendor_users_vendor_id_fkey";

ALTER TABLE "operation_times" DROP CONSTRAINT IF EXISTS "operation_times_brand_id_fkey";

ALTER TABLE "brand_zones" DROP CONSTRAINT IF EXISTS "brand_zones_brand_id_fkey";

ALTER TABLE "brand_zones" DROP CONSTRAINT IF EXISTS "brand_zones_zone_id_fkey";

ALTER TABLE "brands" DROP CONSTRAINT IF EXISTS "brands_vendor_id_fkey";

ALTER TABLE "brand_categories" DROP CONSTRAINT IF EXISTS "brand_categories_brand_id_fkey";

ALTER TABLE "brand_categories" DROP CONSTRAINT IF EXISTS "brand_categories_category_id_fkey";

ALTER TABLE "products" DROP CONSTRAINT IF EXISTS "products_category_id_fkey";

ALTER TABLE "products" DROP CONSTRAINT IF EXISTS "products_brand_id_fkey";

ALTER TABLE "inventory_history" DROP CONSTRAINT IF EXISTS "inventory_history_product_id_fkey";

ALTER TABLE "product_inventory" DROP CONSTRAINT IF EXISTS "product_inventory_product_id_fkey";

ALTER TABLE "product_variants" DROP CONSTRAINT IF EXISTS "product_variants_product_id_fkey";

ALTER TABLE "product_variant_items" DROP CONSTRAINT IF EXISTS "product_variant_items_variant_id_fkey";

ALTER TABLE "product_variant_items" DROP CONSTRAINT IF EXISTS "product_variant_items_product_id_fkey";


DROP TABLE IF EXISTS "product_variant_items";

DROP TABLE IF EXISTS "product_variants";

DROP TABLE IF EXISTS "product_inventory";

DROP TABLE IF EXISTS "inventory_history";

DROP TABLE IF EXISTS "products";

DROP TABLE IF EXISTS "brand_categories";

DROP TABLE IF EXISTS "categories";

DROP TABLE IF EXISTS "brands";

DROP TABLE IF EXISTS "brand_zones";

DROP TABLE IF EXISTS "zones";

DROP TABLE IF EXISTS "operation_times";

DROP TABLE IF EXISTS "vendor_users";

DROP TABLE IF EXISTS "vendor_permissions";

DROP TABLE IF EXISTS "vendors";
