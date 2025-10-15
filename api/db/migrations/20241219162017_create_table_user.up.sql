CREATE TABLE users (
    "id" serial NOT NULL PRIMARY KEY,
    "name" varchar(200) NOT NULL,
    "email" varchar(200) NOT NULL UNIQUE,
    "password" varchar(200) NOT NULL);