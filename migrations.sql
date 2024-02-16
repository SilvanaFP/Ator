BEGIN;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    password TEXT
);

CREATE TABLE rewards (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    required_points INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    reward_id BIGINT REFERENCES rewards(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE TABLE tasks_done (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    task_id BIGINT REFERENCES tasks(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    verified BOOLEAN
);

COMMIT;
