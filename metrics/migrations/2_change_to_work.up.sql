DROP TABLE task;

CREATE TABLE work (
    id BIGSERIAL PRIMARY KEY,
    branch TEXT NOT NULL,
    pull_request INTEGER,
    merge_commit CHARACTER VARYING(40),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    merged_time TIMESTAMP WITH TIME ZONE,
    deployed_time TIMESTAMP WITH TIME ZONE
);
