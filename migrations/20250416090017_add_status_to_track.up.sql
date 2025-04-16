CREATE TYPE track_status_type AS ENUM ('planned', 'in_process', 'completed');

ALTER TABLE track
ADD status track_status_type NOT NULL DEFAULT 'planned';
