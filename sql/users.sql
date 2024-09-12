CREATE TABLE IF NOT EXISTS USERS
(
    id            UUID PRIMARY KEY,
    username      VARCHAR(50) UNIQUE  NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100)        NOT NULL,
    user_type    VARCHAR NOT NULL default 'non-admin',
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
alter table USERS
    owner to postgres;


