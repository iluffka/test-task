SELECT 'CREATE DATABASE counter'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'counter')\gexec

CREATE SEQUENCE count_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

CREATE TABLE request_count (
  id INT NOT NULL default nextval('count_id_seq'),
  request_time INTEGER NOT NULL
);