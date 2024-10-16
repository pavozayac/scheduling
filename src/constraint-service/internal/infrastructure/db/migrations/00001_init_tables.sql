-- +goose Up

CREATE TABLE schedules (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL
);

CREATE TABLE workers (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    schedule_id integer REFERENCES schedules(id) NOT NULL
);

CREATE TABLE tasks (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    story varchar(1000) NOT NULL,
    schedule_id integer REFERENCES schedules(id) NOT NULL
);

CREATE TABLE locations (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    story text(1000) NOT NULL,
    schedule_id integer REFERENCES schedules(id) NOT NULL
);

CREATE TYPE constraint_type as ENUM ('must', 'cant');

CREATE TABLE constraints (
    schedule_id integer REFERENCES schedules(id),
    location_id integer REFERENCES locations(id),
    task_id integer REFERENCES tasks(id),
    worker_id integer REFERENCES workers(id),
    start_slot integer,
    end_slot integer,
    kind constraint_type NOT NULL,
    UNIQUE (location_id, task_id, worker_id, start_slot, end_slot, kind)
);

-- +goose Down

DROP TABLE schedules;
DROP TABLE workers;
DROP TABLE tasks;
DROP TABLE locations;
DROP TYPE constraint_type;
DROP TABLE constraints;
