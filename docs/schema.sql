CREATE SCHEMA "inventory";

CREATE SCHEMA "transaction";

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

CREATE TABLE "transaction"."stocks" (
  "id" uuid PRIMARY KEY,
  "order_no" varchar NOT NULL,
  "product_code" varchar NOT NULL,
  "amount" int NOT NULL,
  "status" stock_status NOT NULL,
  "expired_at" timestamptz,
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

CREATE INDEX ON "transaction"."stocks" ("order_no");

CREATE INDEX ON "transaction"."stocks" ("product_code");

CREATE INDEX ON "transaction"."stocks" ("expired_at");

CREATE INDEX ON "transaction"."stocks" ("status");

CREATE INDEX ON "order"."orders" ("order_no");

CREATE INDEX ON "order"."order_details" ("order_id");

ALTER TABLE "transaction"."stocks" ADD FOREIGN KEY ("product_code") REFERENCES "inventory"."products" ("code");

ALTER TABLE "order"."order_details" ADD FOREIGN KEY ("order_id") REFERENCES "order"."orders" ("id");
