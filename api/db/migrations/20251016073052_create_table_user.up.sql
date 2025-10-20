-- create user_role
CREATE TYPE user_role AS ENUM ('SUPER_ADMIN', 'OWNER', 'WAREHOUSE_HEAD', 'TREASURER');

-- create table user
CREATE TABLE users (
    id VARCHAR(200) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    role user_role NOT NULL,
    password VARCHAR(200) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);


INSERT INTO users (id, name, username, phone, role, password) 
VALUES ('7f5dc73e-e097-4e8e-ba7c-5ed828fabc74', 'Super Admin', 'superadmin', '0888888888', 'SUPER_ADMIN', '$2a$10$ibmW.pidc9yRifFckQJFZ.1Hs4BEkf8.B.b5.xSJio3nsR.3y7rQO');