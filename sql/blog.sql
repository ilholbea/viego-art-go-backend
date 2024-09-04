CREATE TABLE IF NOT EXISTS blog
(
    id      UUID    NOT NULL
        CONSTRAINT pk_blog_id PRIMARY KEY,
    content VARCHAR,
    title   VARCHAR NOT NULL
        CONSTRAINT uq_title_email UNIQUE
);
alter table blog
    owner to postgres;

