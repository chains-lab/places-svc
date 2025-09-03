-- +migration Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE TYPE "place_category" AS ENUM (
    'food and drinks',
    'shops',
    'services',
    'hotels and accommodation',
    'active leisure',
    'religion',
    'offices and factories',
    'residential buildings',
    'education'
);

CREATE TABLE "place_types"  (
    "name"       VARCHAR(256)   PRIMARY KEY,
    "category"   place_category NOT NULL,
    "updated_at" TIMESTAMP      NOT NULL DEFAULT now(),
    "created_at" TIMESTAMP      NOT NULL DEFAULT now()
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
    "status"         place_statuses         NOT NULL,
    "type"           VARCHAR(255)           NOT NULL REFERENCES place_types(name) ON DELETE RESTRICT,
    "verified"       BOOLEAN                NOT NULL DEFAULT FALSE,
    "coords"         geography(POINT, 4326) NOT NULL,
    "ownership"      place_ownership        NOT NULL,

    "website"        VARCHAR(255),
    "phone"          VARCHAR(255),

    "updated_at"     TIMESTAMP              NOT NULL DEFAULT now(),
    "created_at"     TIMESTAMP              NOT NULL DEFAULT now()
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

    PRIMARY KEY (place_id, language)
);

-- +migration Down
DROP TABLE IF EXISTS "places" CASCADE;
DROP TABLE IF EXISTS "place_types" CASCADE;
DROP TABLE IF EXISTS "place_details" CASCADE;

DROP TYPE IF EXISTS "place_category";
DROP TYPE IF EXISTS "place_statuses";
DROP TYPE IF EXISTS "place_ownership";

DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "postgis";
