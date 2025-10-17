-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE TYPE "place_statuses" AS ENUM (
    'active',
    'inactive',
    'blocked'
);

CREATE TABLE "places" (
    "id"         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "city_id"    UUID NOT NULL,
    "company_id" UUID,
    "class"      VARCHAR(32) NOT NULL REFERENCES place_classes(code) ON DELETE RESTRICT ON UPDATE CASCADE,

    "status"     place_statuses         NOT NULL,
    "verified"   BOOLEAN                NOT NULL DEFAULT FALSE,
    "point"      geography(POINT, 4326) NOT NULL,
    "address"    VARCHAR(255)           NOT NULL,

    "website"    VARCHAR(255),
    "phone"      VARCHAR(32),

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TABLE place_i18n (
    "place_id"    UUID         NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    "locale"      VARCHAR(2)   NOT NULL,
    "name"        VARCHAR(128)  NOT NULL,
    "description" VARCHAR(1024) NOT NULL,

    CHECK (locale ~ '^[a-z]{2}$'),
    PRIMARY KEY (place_id, locale)
);

-- +migrate Down
DROP TABLE IF EXISTS "place_i18n" CASCADE;
DROP TABLE IF EXISTS "places" CASCADE;

DROP TYPE IF EXISTS "place_statuses";
DROP TYPE IF EXISTS "place_ownership";

DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "postgis";
