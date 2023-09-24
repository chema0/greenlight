CREATE DATABASE greenlight;

CREATE ROLE greenlight WITH LOGIN PASSWORD 'greenlight';

ALTER DATABASE greenlight OWNER TO greenlight;

\connect greenlight;

CREATE EXTENSION IF NOT EXISTS citext;