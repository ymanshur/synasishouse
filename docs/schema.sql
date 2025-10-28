CREATE SCHEMA "inventory";

CREATE SCHEMA "order";

CREATE TYPE "stock_status" AS ENUM (
  'reserved',
  'setlled',
  'canceled'
);

CREATE TABLE "inventory"."products" (
  "id" uuid PRIMARY KEY,
  "code" varchar UNIQUE NOT NULL,
  "total" int NOT NULL DEFAULT 0,
  "hold" int NOT NULL DEFAULT 0,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order"."orders" (
  "id" uuid PRIMARY KEY,
  "order_no" varchar UNIQUE NOT NULL,
  "user_id" uuid NOT NULL,
  "status" stock_status NOT NULL,
  "expired_at" timestamptz,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order"."order_details" (
  "id" uuid PRIMARY KEY,
  "order_id" varchar NOT NULL,
  "product_code" varchar NOT NULL,
  "amount" int NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "inventory"."products" ("code");

CREATE INDEX ON "order"."orders" ("order_no");

CREATE INDEX ON "order"."order_details" ("order_id");

ALTER TABLE "order"."order_details" ADD FOREIGN KEY ("order_id") REFERENCES "order"."orders" ("id");
