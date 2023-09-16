CREATE DATABASE greenlight;

CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';

ALTER DATABASE greenlight OWNER TO greenlight;

\connect greenlight;

CREATE EXTENSION IF NOT EXISTS citext;