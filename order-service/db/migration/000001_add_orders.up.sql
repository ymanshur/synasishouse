CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "orders" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "order_no" varchar UNIQUE NOT NULL,
  "user_id" uuid NOT NULL,
  "status" varchar(100) NOT NULL,
  "expired_at" timestamptz,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "orders" ("order_no");

CREATE TABLE "order_details" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "order_id" uuid NOT NULL,
  "product_code" varchar NOT NULL,
  "amount" int NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "order_details" ("order_id");

ALTER TABLE "order_details" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
