CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.sessions (
    session_uid VARCHAR(128)    NOT NULL    PRIMARY KEY    CHECK (session_uid <> ''),
    customer_id INTEGER         NOT NULL                   CHECK (customer_id > 0),
    created_at  INTEGER         NOT NULL DEFAULT extract(epoch FROM current_timestamp(0))
);

CREATE INDEX IF NOT EXISTS customer_id_idx ON auth.sessions (customer_id);

CREATE TABLE IF NOT EXISTS auth.customers (
    id          SERIAL          NOT NULL PRIMARY KEY,
    phone       VARCHAR(18)     NOT NULL,
    email       VARCHAR(320)    NOT NULL DEFAULT '',
    firstname   VARCHAR(256)    NOT NULL,
    lastname    VARCHAR(256)    NOT NULL,
    patronymic  VARCHAR(256)    DEFAULT '',
    passphrase  VARCHAR(128)    NOT NULL,
    UNIQUE (id),
    UNIQUE (phone),
    UNIQUE (email)
);
