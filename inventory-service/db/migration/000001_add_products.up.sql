CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "products" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "code" varchar NOT NULL UNIQUE,
  "total" int NOT NULL DEFAULT 0,
  "hold" int NOT NULL DEFAULT 0,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "products" ("code");

