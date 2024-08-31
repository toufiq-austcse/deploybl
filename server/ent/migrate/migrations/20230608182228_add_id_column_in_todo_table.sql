-- Create "todos" table
CREATE TABLE "todos" ("id" uuid NOT NULL, "title" character varying NOT NULL DEFAULT '', "status" character varying NOT NULL DEFAULT '', PRIMARY KEY ("id"));
