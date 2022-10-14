CREATE TABLE "products" (
  "barcode" TEXT PRIMARY KEY,
  "name" TEXT,
  "quantity" TEXT,
  "image_url" TEXT,
  "energy_100g" FLOAT,
  "energy_serving" FLOAT,
  "nutrient_levels_id" INT,
  "nova_group" INT,
  "nutriscore_grade" TEXT
);

CREATE TABLE "brands" (
  "tag" TEXT PRIMARY KEY
);

CREATE TABLE "product_brands" (
  "barcode" TEXT,
  "tag" TEXT,
  PRIMARY KEY ("barcode", "tag")
);

CREATE TABLE "nutrient_levels" (
  "id" SERIAL PRIMARY KEY,
  "fat" TEXT NOT NULL,
  "saturated_fat" TEXT NOT NULL,
  "sugar" TEXT NOT NULL,
  "salt" TEXT NOT NULL
);

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT NOT NULL,
  "email" TEXT UNIQUE NOT NULL,
  "password" TEXT NOT NULL
);

CREATE TABLE "ssds" (
  "id" SERIAL NOT NULL,
  "user_id" INT NOT NULL,
  "mac_address" TEXT UNIQUE NOT NULL,
  UNIQUE ("id", "user_id"),
  PRIMARY KEY ("id")
);

CREATE TABLE "product_ssds" (
  "ssd_id" INT NOT NULL,
  "barcode" TEXT NOT NULL,
  "quantity" INT NOT NULL,
  PRIMARY KEY ("ssd_id", "barcode")
);

CREATE TABLE "product_orders" (
  "id" SERIAL NOT NULL,
  "ssd_id" INT,
  "product_id" TEXT,
  "timestamp" timestamp NOT NULL,
  "quantity" INT NOT NULL,
  "status" TEXT NOT NULL,
  PRIMARY KEY ("id")
);

ALTER TABLE "products" ADD FOREIGN KEY ("nutrient_levels_id") REFERENCES "nutrient_levels" ("id");

ALTER TABLE "product_brands" ADD FOREIGN KEY ("barcode") REFERENCES "products" ("barcode");

ALTER TABLE "product_brands" ADD FOREIGN KEY ("tag") REFERENCES "brands" ("tag");

ALTER TABLE "ssds" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "product_ssds" ADD FOREIGN KEY ("ssd_id") REFERENCES "ssds" ("id");

ALTER TABLE "product_ssds" ADD FOREIGN KEY ("barcode") REFERENCES "products" ("barcode");

ALTER TABLE "product_orders" ADD FOREIGN KEY ("ssd_id") REFERENCES "ssds" ("id");

ALTER TABLE "product_orders" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("barcode");
