
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE controller_events
(
  controller_event_id serial NOT NULL,
  event_timestamp timestamp with time zone,
  name character varying,
  role character varying,
  callsign character varying,
  primary_frequency integer,
  latitude float,
  longitude float,
  CONSTRAINT controller_events_pkey PRIMARY KEY (controller_event_id)
);

CREATE INDEX controller_event_id_index
  ON controller_events (controller_event_id );

CREATE TABLE pilot_events
(
  pilot_event_id serial NOT NULL,
  event_timestamp timestamp with time zone,
  cid integer,
  name character varying,
  equipment character varying,
  callsign character varying,
  frequency integer,
  radio character varying,
  desired_role character varying,
  latitude float,
  longitude float,
  altitude float,
  ground_speed float,
  true_heading float,
  flight_plan_origin character varying,
  flight_plan_destination character varying,
  flight_plan_route character varying,
  flight_plan_remarks character varying,
  CONSTRAINT pilot_events_pkey PRIMARY KEY (pilot_event_id)
);

CREATE INDEX pilot_event_id_index
  ON pilot_events (pilot_event_id );

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE controller_events;

DROP TABLE pilot_events;