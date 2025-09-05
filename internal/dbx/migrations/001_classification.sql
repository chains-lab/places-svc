-- +migrate Up
CREATE TYPE categories_statuses AS ENUM (
    'active',
    'deprecated'
);

CREATE TABLE place_categories (
    code       VARCHAR(16)         PRIMARY KEY,
    status     categories_statuses NOT NULL DEFAULT 'active',
    icon       VARCHAR(255)        NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CHECK (code ~ '^[a-z_]{1,16}$')
);

CREATE TYPE kind_statuses AS ENUM (
    'active',
    'deprecated'
);

CREATE TABLE place_kinds (
    code          VARCHAR(16)   PRIMARY KEY,
    category_code VARCHAR(16)   NOT NULL REFERENCES place_categories(code) ON DELETE RESTRICT ON UPDATE CASCADE,
    status        kind_statuses NOT NULL DEFAULT 'active',
    icon          VARCHAR(255)  NOT NULL,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CHECK (code ~ '^[a-z_]{1,16}$')
);

CREATE TABLE place_category_i18n (
    category_code VARCHAR(16) NOT NULL REFERENCES place_categories(code) ON DELETE CASCADE,
    locale        varchar(2) NOT NULL,
    name          VARCHAR(255) NOT NULL,

    UNIQUE (name, locale),
    PRIMARY KEY (category_code, locale),
    CHECK (locale ~ '^[a-z]{2}(-[A-Z]{2})?$')
);

CREATE TABLE place_kind_i18n (
    kind_code     VARCHAR(16) NOT NULL REFERENCES place_kinds(code) ON DELETE CASCADE,
    locale        varchar(2) NOT NULL,
    name          VARCHAR(255) NOT NULL,

    UNIQUE (name, locale),
    PRIMARY KEY (kind_code, locale),
    CHECK (locale ~ '^[a-z]{2}(-[A-Z]{2})?$')
);

-- +migrate Down
DROP TABLE IF EXISTS place_category_i18n CASCADE;
DROP TABLE IF EXISTS place_kind_i18n CASCADE;
DROP TABLE IF EXISTS place_categories CASCADE;
DROP TABLE IF EXISTS place_kinds CASCADE;

DROP TYPE IF EXISTS kind_statuses;
DROP TYPE IF EXISTS categories_statuses;
