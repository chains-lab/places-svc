-- +migration Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "place_category" AS ENUM (
    'store',
    'restaurant',
    'zoo'
);

CREATE TABLE "places_types"  (
    "id"         UUID           PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "name"       VARCHAR(256)   NOT NULL UNIQUE,
    "category"   place_category NOT NULL,
    "updated_at" TIMESTAMP      NOT NULL DEFAULT now(),
    "created_at" TIMESTAMP      NOT NULL DEFAULT now()
);

-- +migration Down
DROP TABLE IF EXISTS "places_types" CASCADE;

DROP TYPE IF EXISTS "place_category" CASCADE;