CREATE TYPE timeline_status_type AS ENUM ('ready', 'expired', 'completed');

ALTER TABLE timeline
ADD status timeline_status_type NOT NULL DEFAULT 'ready';
