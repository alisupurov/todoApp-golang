CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
    id            SERIAL        PRIMARY KEY,
    version       INTEGER       NOT NULL DEFAULT 1,
    full_name     VARCHAR(100)  NOT NULL CHECK (char_length(full_name) >= 3),
    phone_number  VARCHAR(15)   CHECK (
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length(phone_number) >= 10
    )
);

CREATE TABLE todoapp.accounts (
    id            SERIAL        PRIMARY KEY,
    email         VARCHAR(255)  NOT NULL UNIQUE,
    password_hash VARCHAR(255)  NOT NULL
);

CREATE TABLE todoapp.tasks (
    id            SERIAL                  PRIMARY KEY,
    version       INTEGER       NOT NULL  DEFAULT 1,
    title         VARCHAR(100)  NOT NULL  CHECK (char_length(title) >= 1),
    description   VARCHAR(1000)           CHECK (char_length(description) >= 1),
    completed     BOOLEAN       NOT NULL,
    created_at    TIMESTAMPTZ   NOT NULL,
    completed_at  TIMESTAMPTZ,

    CHECK (
        (completed=FALSE AND completed_at IS NULL)
        OR
        (completed=TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
    ),

    author_user_id INTEGER     NOT NULL   REFERENCES todoapp.users(id)
);

CREATE INDEX idx_tasks_author_user_id ON todoapp.tasks(author_user_id);
CREATE INDEX idx_tasks_created_at     ON todoapp.tasks(created_at);
