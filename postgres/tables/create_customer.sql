CREATE TABLE customer (
    customer_id SERIAL PRIMARY KEY,
    name varchar(20) NOT NULL,
    sex varchar(1) NOT NULL,
    age numeric(3,0) NOT NULL
);
