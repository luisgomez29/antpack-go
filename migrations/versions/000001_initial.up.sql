/*
 * Author: Luis Guillermo Gómez Galeano
 *
 * Create initial database schema.
 *
 * The created_at and updated_at fields have the default value CURRENT_TIMESTAMP.
 */

-- Create users table
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(40)  NOT NULL,
    last_name  VARCHAR(40)  NOT NULL,
    email      VARCHAR(60)  NOT NULL UNIQUE,
    password   VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
