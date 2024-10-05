CREATE TABLE schedules (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
);

CREATE TABLE workers (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    schedule_id integer REFERENCES schedules(id),
);

CREATE TABLE tasks (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    story varchar(1000) NOT NULL,
    schedule_id integer REFERENCES schedules(id),
);

CREATE TABLE locations (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title varchar(200) NOT NULL,
    story text(1000) NOT NULL,
    schedule_id integer REFERENCES schedules(id),
);

CREATE TYPE constraint_type as ENUM ('must', 'cant');

CREATE TABLE worker_time_constraints (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    worker_id integer REFERENCES workers(id);
    start_slot integer NOT NULL;
    end_slot integer NOT NULL;
    kind constraint_type NOT NULL;
);

CREATE TABLE location_time_constraints (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    location_id integer REFERENCES locations(id);
    start_slot integer NOT NULL;
    end_slot integer NOT NULL;
    kind constraint_type NOT NULL;
);

CREATE TABLE task_time_constraints (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    task_id integer REFERENCES tasks(id);
    start_slot integer NOT NULL;
    end_slot integer NOT NULL;
    kind constraint_type NOT NULL;
);

CREATE TABLE worker_task_constraints (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    worker_id integer REFERENCES workers(id);
    task_id integer REFERENCES tasks(id);
    kind constraint_type NOT NULL;
);

CREATE TABLE location_task_constraints (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    location_id integer REFERENCES locations(id);
    task_id integer REFERENCES tasks(id);
    kind constraint_type NOT NULL;
);