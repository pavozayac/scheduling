-- +goose Up

CREATE TABLE schedules (
    id uuid PRIMARY KEY,
    title varchar(200) NOT NULL
);

CREATE TABLE workers (
    id uuid PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    schedule_id uuid REFERENCES schedules(id) NOT NULL
);

CREATE TABLE tasks (
    id uuid PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    story varchar(1000) NOT NULL,
    schedule_id uuid REFERENCES schedules(id) NOT NULL
);

CREATE TABLE locations (
    id uuid PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    story varchar(1000) NOT NULL,
    schedule_id uuid REFERENCES schedules(id) NOT NULL
);

CREATE TYPE constraint_type as ENUM ('must', 'cannot');

CREATE TABLE constraints (
    schedule_id uuid REFERENCES schedules(id),
    location_id uuid REFERENCES locations(id),
    task_id uuid REFERENCES tasks(id),
    worker_id uuid REFERENCES workers(id),
    start_slot integer,
    end_slot integer,
    kind constraint_type NOT NULL,
    UNIQUE (location_id, task_id, worker_id, start_slot, end_slot, kind)
);

-- +goose Down

DROP TABLE constraints;
DROP TABLE locations;
DROP TABLE tasks;
DROP TABLE workers;
DROP TABLE schedules;
DROP TYPE constraint_type;
