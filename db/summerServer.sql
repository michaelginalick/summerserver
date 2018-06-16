DROP DATABASE IF EXISTS "summerserver";
CREATE DATABASE "summerserver";

USE "summerserver";



CREATE TABLE "Events" (
	"id" serial NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"link" VARCHAR(255) NOT NULL,
	"month" VARCHAR(255),
	"days" integer,
	"individual_days" integer,
	"festival_length" integer,
	CONSTRAINT Events_pk PRIMARY KEY ("id")
)



