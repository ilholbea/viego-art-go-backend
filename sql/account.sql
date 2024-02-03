CREATE TABLE IF NOT EXISTS account
(
    id UUID NOT NULL CONSTRAINT pk_account_id PRIMARY KEY,
    full_name VARCHAR,
    email VARCHAR NOT NULL CONSTRAINT uq_account_email UNIQUE,
    hashed_password VARCHAR,
    account_type VARCHAR NOT NULL default 'user'
);

alter table account
    owner to postgres;