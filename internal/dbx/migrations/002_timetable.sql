-- +migration Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE place_timetable (
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    place_id  UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    start_min INT  NOT NULL,
    end_min   INT  NOT NULL,

    CHECK (start_min >= 0 AND end_min <= 10080 AND end_min > start_min),

    EXCLUDE USING gist (
        place_id WITH =,
        int4range(start_min, end_min, '[)') WITH &&
    )
);

CREATE INDEX place_timetable_place_idx ON place_timetable (place_id);

-- +migration Down
DROP INDEX IF EXISTS place_timetable_place_idx; -- если создашь выше
DROP TABLE IF EXISTS place_timetable CASCADE;

DROP EXTENSION IF EXISTS btree_gist;
DROP EXTENSION IF EXISTS "uuid-ossp";
