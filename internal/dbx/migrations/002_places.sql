-- +migration Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE "place_statuses" AS ENUM (
    'active',
    'inactive',
    'blocked'
);

CREATE TYPE "ownership_type" AS ENUM (
    'private',
    'public',
    'municipal',
);

CREATE TABLE "places" (
    "id"             UUID             PRIMARY KEY DEFAULT uuid_generate_v4(),
    "city_id"        UUID             NOT NULL,
    "distributor_id" UUID,
    "type_id"        UUID             NOT NULL REFERENCES "places_types" ("id") ON DELETE CASCADE,
    "ownership"      ownership_type   NOT NULL,
    "status"         place_statuses   NOT NULL,

    "coords"         geography(POINT, 4326) NOT NULL,
    "name"           VARCHAR   NOT NULL,
    "description"    VARCHAR   NOT NULL,
    "address"        VARCHAR   NOT NULL,
    "website"        VARCHAR   NOT NULL,
    "phone"          VARCHAR   NOT NULL,
    "updated_at"     TIMESTAMP NOT NULL DEFAULT now(),
    "created_at"     TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT check_distributor_only_for_private
    CHECK (
        (ownership = 'private'    AND distributor_id IS NOT NULL)
     OR (ownership <> 'private'   AND distributor_id IS NULL)
    )

    CONSTRAINT check_website_not_empty CHECK (website <> ''),
    CONSTRAINT check_phone_format_e164 CHECK (phone ~ '^\+?[0-9]{7,15}$')
);

-- +migration Down
DROP TABLE IF EXISTS "places" CASCADE;

DROP TYPE IF EXISTS "ownership_type" CASCADE;
DROP TYPE IF EXISTS "place_statuses" CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp";
