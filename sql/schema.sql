--table of courses
CREATE TABLE courses
(
    id          BIGSERIAL   PRIMARY KEY,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    status      TEXT        NOT NULL DEFAULT 'Not Started',
    created     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- table of users
CREATE TABLE users 
(
    id       BIGSERIAL   PRIMARY KEY,
    username TEXT        NOT NULL UNIQUE,
    password TEXT        NOT NULL,
    is_admin BOOLEAN     NOT NULL DEFAULT FALSE,
    active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- table of users_courses
CREATE TABLE users_courses
(
    user_id     BIGINT      NOT NULL REFERENCES users,
    course_id   BIGINT      NOT NULL REFERENCES courses,
    created     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--table of users_tokens
CREATE TABLE users_tokens
(
    token       TEXT        NOT NULL    UNIQUE,
    user_id     BIGINT      NOT NULL    REFERENCES users,
    expires     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);
