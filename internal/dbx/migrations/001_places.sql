-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE TABLE place_categories (
    id          VARCHAR(50) PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CHECK (id ~ '^[a-z_]{1,50}$')
);

CREATE TABLE place_kinds (
    id           VARCHAR(50) PRIMARY KEY,
    category_id  VARCHAR(50) NOT NULL REFERENCES place_categories(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    name         VARCHAR(50) NOT NULL UNIQUE,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CHECK (id ~ '^[a-z_]{1,50}$')
);


CREATE TYPE "place_statuses" AS ENUM (
    'active',
    'inactive',
    'blocked'
);

CREATE TYPE "place_ownership" AS ENUM (
    'private',
    'common'
);

CREATE TABLE "places" (
    "id"             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "city_id"        UUID NOT NULL,
    "distributor_id" UUID,
    "type_id"        VARCHAR(50) NOT NULL REFERENCES place_kinds(id) ON DELETE RESTRICT ON UPDATE CASCADE,

    "status"         place_statuses         NOT NULL,
    "verified"       BOOLEAN                NOT NULL DEFAULT FALSE,
    "ownership"      place_ownership        NOT NULL,
    "point"          geography(POINT, 4326) NOT NULL,

    "website"        VARCHAR(255),
    "phone"          VARCHAR(255),

    "created_at"     TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    "updated_at"     TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TYPE "language_code" AS ENUM (
    'en','uk','ru'
);

CREATE TABLE "place_details" (
    "place_id"    UUID          NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    "language"    language_code NOT NULL,
    "name"        VARCHAR       NOT NULL,
    "address"     VARCHAR       NOT NULL,
    "description" VARCHAR,

    "created_at"   TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    "updated_at"   TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    PRIMARY KEY (place_id, language)
);

-- +migrate Down
DROP TABLE IF EXISTS "place_details" CASCADE;
DROP TABLE IF EXISTS "places" CASCADE;
DROP TABLE IF EXISTS "place_kinds" CASCADE;
DROP TABLE IF EXISTS "place_categories" CASCADE;

DROP TYPE IF EXISTS "place_statuses";
DROP TYPE IF EXISTS "place_ownership";
DROP TYPE IF EXISTS "language_code";

DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "postgis";
