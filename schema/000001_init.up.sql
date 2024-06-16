CREATE TABLE authors
(
    id            serial PRIMARY KEY,
    first_name    varchar(50) NOT NULL DEFAULT '',
    last_name     varchar(50) NOT NULL DEFAULT '',
    bio           text        NOT NULL DEFAULT '',
    date_of_birth date
);

CREATE TABLE books
(
    id serial PRIMARY KEY,
        title varchar (255) not null,
        author integer references authors(id),
        year_of_publication integer not null default 0,
        isbn varchar (17) unique
);

