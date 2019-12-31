DROP DATABASE IF EXISTS "summerserver";
CREATE DATABASE "summerserver";

CREATE TABLE events (
	id serial NOT NULL,
	name VARCHAR(255) NOT NULL,
	link VARCHAR(255) NOT NULL,
	month VARCHAR(255),
	year VARCHAR(255),
	location VARCHAR(255),
	CONSTRAINT Events_pk PRIMARY KEY (id)
);


CREATE TABLE days (
	id serial NOT NULL,
	day integer NOT NULL,
	event_id integer REFERENCES events
);

CREATE INDEX events_id ON events (id);
CREATE INDEX days_id ON days (id);
CREATE INDEX event_id ON days (event_id);


