-- +migrate Up
CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TYPE place_class_statuses AS ENUM (
    'active',
    'inactive'
);

-- Основная таблица-дерево
CREATE TABLE IF NOT EXISTS place_classes (
    code         VARCHAR(32)          PRIMARY KEY,
    parent       VARCHAR(32)          NULL REFERENCES place_classes(code) ON DELETE RESTRICT ON UPDATE CASCADE,
    status       place_class_statuses NOT NULL DEFAULT 'active',
    icon         VARCHAR(255)         NOT NULL,

    path         LTREE         NOT NULL,      -- материализованный путь, напр. 'cars.suv.compact'
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CHECK (code ~ '^[a-z_]{1,16}$'),
    CHECK (parent IS NULL OR parent <> code)
);

CREATE TABLE IF NOT EXISTS place_class_i18n (
    class  VARCHAR(32)  NOT NULL REFERENCES place_classes(code) ON DELETE CASCADE,
    locale VARCHAR(2)   NOT NULL,
    name   VARCHAR(32) NOT NULL,

    CHECK (locale ~ '^[a-z]{2}$'),
    PRIMARY KEY (class, locale),
    UNIQUE (name, locale)
);

-- 1) BEFORE INSERT/UPDATE: вычислить/пересчитать path
--   - на INSERT: path = (parent.path || code) или просто code
--   - на UPDATE code/parent: пересчитать path узла и всего поддерева
--   - анти-цикл: запрещаем ставить parent из своего поддерева

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION place_class_before_ins_upd_path() RETURNS TRIGGER
AS $pc$
DECLARE
    parent_path  ltree;
    old_path     ltree;
    new_path     ltree;
    tail_start   int;
BEGIN
    IF TG_OP = 'INSERT' THEN
    -- получить путь родителя
        IF NEW.parent IS NOT NULL THEN
            SELECT path INTO parent_path FROM place_classes WHERE code = NEW.parent;
            IF parent_path IS NULL THEN
                RAISE EXCEPTION 'parent % not found for %', NEW.parent, NEW.code;
            END IF;
            NEW.path := parent_path || NEW.code::ltree;
        ELSE
            NEW.path := NEW.code::ltree;
        END IF;

        NEW.updated_at := now() AT TIME ZONE 'UTC';
        NEW.created_at := COALESCE(NEW.created_at, now() AT TIME ZONE 'UTC');
        RETURN NEW;
    END IF;

    -- UPDATE
    IF TG_OP = 'UPDATE' THEN
    -- если ни code, ни parent не менялись — просто обновим updated_at и выйдем
        IF NEW.code = OLD.code AND COALESCE(NEW.parent,'') = COALESCE(OLD.parent,'') THEN
            NEW.updated_at := now() AT TIME ZONE 'UTC';
            RETURN NEW;
        END IF;

        old_path := OLD.path;

        -- анти-цикл: нельзя поставить родителя из своего поддерева
        IF NEW.parent IS NOT NULL THEN
            PERFORM 1
            FROM place_classes p
            WHERE p.code = NEW.parent
                AND p.path <@ old_path;  -- родитель лежит внутри нашего поддерева?
            IF FOUND THEN
                RAISE EXCEPTION 'cycle detected: % cannot be parent of %', NEW.parent, NEW.code;
            END IF;

            SELECT path INTO parent_path FROM place_classes WHERE code = NEW.parent;
            IF parent_path IS NULL THEN
                RAISE EXCEPTION 'parent % not found for %', NEW.parent, NEW.code;
            END IF;
            new_path := parent_path || NEW.code::ltree;
        ELSE
            new_path := NEW.code::ltree;
        END IF;

        NEW.path := new_path;
        NEW.updated_at := now() AT TIME ZONE 'UTC';

        -- если путь изменился (сменился код или родитель) — пересчитать поддерево
        IF new_path <> old_path THEN
            -- пересчёт путей потомков:
            -- tail = subpath(child.path, nlevel(old_path))
            -- child.path = new_path || tail
            WITH RECURSIVE subtree AS (
                SELECT code, path
                FROM place_classes
                WHERE path <@ old_path AND code <> OLD.code  -- все потомки (без самого узла)
            )
            UPDATE place_classes c
            SET path = NEW.path || subpath(c.path, nlevel(old_path)),
                updated_at = now() AT TIME ZONE 'UTC'
            FROM subtree s
            WHERE c.code = s.code;
        END IF;

        RETURN NEW;
    END IF;

    RETURN NEW;
END;
$pc$ LANGUAGE plpgsql;

CREATE TRIGGER trg_place_class_before_ins_upd_path
BEFORE INSERT OR UPDATE OF code, parent
ON place_classes
FOR EACH ROW
EXECUTE FUNCTION place_class_before_ins_upd_path();
-- +migrate StatementEnd

-- 2) AFTER UPDATE OF status: каскадная деактивация
--   Если узел стал inactive, все потомки тоже становятся inactive.
--   (Т.к. дерево маленькое, допускаем повторные срабатывания у потомков — это ок.)

-- +migrate StatementBegin
CREATE FUNCTION place_class_after_update_status() RETURNS TRIGGER
AS $pc$
BEGIN
    IF TG_OP = 'UPDATE' AND NEW.status = 'inactive' AND OLD.status <> 'inactive' THEN
        UPDATE place_classes
        SET status = 'inactive',
            updated_at = now() AT TIME ZONE 'UTC'
        WHERE path <@ NEW.path         -- селектим узел и всех потомков
            AND code <> NEW.code         -- кроме уже обновлённого
            AND status <> 'inactive';  -- обновляем только активных
    END IF;

    RETURN NEW;
END;
$pc$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER trg_place_class_after_update_status
AFTER UPDATE OF status
ON place_classes
FOR EACH ROW
EXECUTE FUNCTION place_class_after_update_status();


-- 3) (опционально, но полезно) запретить активацию под предком-inactive
--   Нельзя поставить active, если любой предок inactive.

-- +migrate StatementBegin
CREATE FUNCTION place_class_check_activate_under_inactive() RETURNS TRIGGER
AS $pc$
DECLARE
    has_depr boolean;
BEGIN
    IF TG_OP = 'UPDATE' AND NEW.status = 'active' AND OLD.status <> 'active' THEN
    -- предки: те, кто является префиксом нашего пути, исключая сам узел
        SELECT EXISTS (
            SELECT 1
            FROM place_classes anc
            WHERE NEW.path <@ anc.path   -- anc — предок NEW
                AND anc.code <> NEW.code
                AND anc.status = 'inactive'
        ) INTO has_depr;

        IF has_depr THEN
            RAISE EXCEPTION 'cannot activate node % under inactive ancestor', NEW.code;
        END IF;
    END IF;
    RETURN NEW;
END;
$pc$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER trg_place_class_check_activate_under_inactive
BEFORE UPDATE OF status
ON place_classes
FOR EACH ROW
EXECUTE FUNCTION place_class_check_activate_under_inactive();

CREATE INDEX IF NOT EXISTS place_class_i18n_locale_idx ON place_class_i18n (locale, name);
CREATE INDEX IF NOT EXISTS place_class_parent_idx ON place_classes (parent);
CREATE INDEX IF NOT EXISTS place_class_path_gist   ON place_classes USING GIST (path);
CREATE INDEX IF NOT EXISTS place_class_status_idx  ON place_classes (status);

-- +migrate Down

DROP TRIGGER IF EXISTS trg_place_class_before_ins_upd_path ON place_classes;
DROP TRIGGER IF EXISTS trg_place_class_after_update_status ON place_classes;
DROP TRIGGER IF EXISTS trg_place_class_check_activate_under_inactive ON place_classes;

DROP FUNCTION IF EXISTS place_class_before_ins_upd_path();
DROP FUNCTION IF EXISTS place_class_after_update_status();
DROP FUNCTION IF EXISTS place_class_check_activate_under_inactive();

DROP TABLE IF EXISTS place_class_i18n;

DROP TABLE IF EXISTS place_classes;

DROP TYPE IF EXISTS place_class_statuses;

DROP EXTENSION IF EXISTS ltree;
