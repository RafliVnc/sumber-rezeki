-- create user_role
CREATE TYPE user_role AS ENUM ('SUPER_ADMIN', 'OWNER', 'WAREHOUSE_HEAD', 'TREASURER');

-- create table user
CREATE TABLE users (
    id VARCHAR(200) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    role user_role NOT NULL,
    password VARCHAR(200) NOT NULL
);
