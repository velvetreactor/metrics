CREATE TABLE activities (
    id SERIAL PRIMARY KEY,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
)