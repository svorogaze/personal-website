CREATE TABLE author (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(75) NOT NULL,
    avatar_link text NOT NULL,
    login varchar(32) NOT NULL UNIQUE,
    password varchar(255) NOT NULL
);

CREATE TABLE blog (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(255) NOT NULL,
    description text NOT NULL,
    contents text NOT NULL,
    creation_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    authorid integer NOT NULL,
    picturelink text NOT NULL,
    FOREIGN KEY (authorid) REFERENCES author(id)
);