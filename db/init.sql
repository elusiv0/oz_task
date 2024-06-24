CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    _text VARCHAR,
    title VARCHAR,
    closed BOOLEAN,
    created_at timestamp not null default current_timestamp
);
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    _text VARCHAR(2000), 
    article_id int REFERENCES posts (id),
    parent_id int REFERENCES comments (id),
    created_at timestamp not null default current_timestamp
);