-- Tabel SALES
CREATE TABLE sales (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP 
);

-- Tabel ROUTE
CREATE TABLE "route" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP 
);

-- Tabel MANY-TO-MANY
CREATE TABLE sales_routes (
    id SERIAL PRIMARY KEY,
    sales_id INT NOT NULL,
    route_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP ,
    CONSTRAINT fk_sales FOREIGN KEY (sales_id)
        REFERENCES sales(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_route FOREIGN KEY (route_id)
        REFERENCES route(id)
        ON DELETE CASCADE,
    CONSTRAINT uq_sales_route UNIQUE (sales_id, route_id)
);
