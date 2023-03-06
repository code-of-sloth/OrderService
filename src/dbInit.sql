CREATE DATABASE IF NOT EXISTS order_db;

USE order_db;

CREATE TABLE IF NOT EXISTS orders (
    id int(11) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    currency_unit VARCHAR(10) NOT NULL DEFAULT 'INR',
    created timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS items (
    id int(11) NOT NULL AUTO_INCREMENT,
    description CHAR(100) NOT NULL DEFAULT '',
    price FLOAT DEFAULT 0,
    qty int(2) DEFAULT 0,
    order_ids TEXT,
    created timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);
