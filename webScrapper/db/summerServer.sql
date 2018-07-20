DROP DATABASE IF EXISTS "summerserver";
CREATE DATABASE "summerserver";

CREATE TABLE events (
	id serial NOT NULL,
	name VARCHAR(255) NOT NULL,
	link VARCHAR(255) NOT NULL,
	month VARCHAR(255),
	days VARCHAR(255),
	year VARCHAR(255),
	individual_days VARCHAR(255),
	festival_length VARCHAR(255),
	CONSTRAINT Events_pk PRIMARY KEY (id)
);



