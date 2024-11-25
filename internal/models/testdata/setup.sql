CREATE SCHEMA test_schema;

SET search_path TO test_schema;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(15) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL
);

CREATE TABLE thoughts (
    id SERIAL PRIMARY KEY,
    content VARCHAR(200) NOT NULL,
    created TIMESTAMP NOT NULL,
    agreecount INTEGER NOT NULL DEFAULT 0,
    disagreecount INTEGER NOT NULL DEFAULT 0,
    user_id INTEGER NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);

CREATE INDEX idx_created ON thoughts (created);

INSERT INTO users(username, email, name, hashed_password, created) VALUES (
    'amorgan',
    'morgan@rdr.com',
    'Arthur Morgan',
    '$2a$12$/q89avx.dI4Mg5xAWlbRjOsCte1GtdCCrcttimfI9OS1iZlPI/PZq',
    NOW()
);

INSERT INTO thoughts(content, created, user_id) VALUES (
    'This is some test content that lies on the Test Database',
    NOW(),
    1
);