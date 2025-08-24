-- +migration Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

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

CREATE TABLE place_types  (
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
    'distributor',
    'municipality'
--  'unclaimed' TODO: future use
);

CREATE TABLE "places" (
    "id"             UUID                   PRIMARY KEY DEFAULT uuid_generate_v4(),
    "distributor_id" UUID,
    "type"           VARCHAR(256)           NOT NULL REFERENCES place_types(name) ON DELETE RESTRICT,
    "status"         place_statuses         NOT NULL,
    "ownership"      place_ownership        NOT NULL,
    "coords"         geography(POINT, 4326) NOT NULL,

    "website"        VARCHAR,
    "phone"          VARCHAR,
    "updated_at"     TIMESTAMP              NOT NULL DEFAULT now(),
    "created_at"     TIMESTAMP              NOT NULL DEFAULT now()
);

CREATE TABLE "place_details" (
    "place_id"    UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    "language"    VARCHAR(16) NOT NULL,
    "name"        VARCHAR    NOT NULL,
    "address"     VARCHAR    NOT NULL,
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
