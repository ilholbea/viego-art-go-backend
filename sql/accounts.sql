CREATE TABLE IF NOT EXISTS accounts
(
    id              UUID    NOT NULL
        CONSTRAINT pk_account_id PRIMARY KEY,
    name            VARCHAR,
    email           VARCHAR NOT NULL
        CONSTRAINT uq_account_email UNIQUE,
    password VARCHAR,
    account_type    VARCHAR NOT NULL default 'user'
);
alter table accounts
    owner to postgres;

