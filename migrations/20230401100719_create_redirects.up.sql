CREATE TABLE redirects(
    id bigserial not null primary key,
    source varchar not null,
    destination varchar not null,
    is_active boolean not null
);