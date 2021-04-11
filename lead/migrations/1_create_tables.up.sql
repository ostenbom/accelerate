CREATE TABLE task (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE
);
