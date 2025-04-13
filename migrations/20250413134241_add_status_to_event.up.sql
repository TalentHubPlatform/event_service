CREATE TYPE event_status_type AS ENUM ('planned', 'in_process', 'completed');

ALTER TABLE event
ADD status event_status_type NOT NULL DEFAULT 'planned';
