--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.5
-- Dumped by pg_dump version 9.6.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;


create table  if not exists usr
(
    id            bigint generated by default as identity
        primary key,
    password         text not null,
    email         text not null
        unique,
    shortening_url_limit    bigint default 1 not null,
    account_type text not null default 'free'::text, -- could be free or premium
    is_active     boolean default true not null,
    zlins_dttm    timestamp with time zone,
    zlupd_dttm    timestamp with time zone
);

create table if not exists url
(
    id bigint generated by default as identity
        constraint url_pkey
            primary key,
    long_version text not null unique,
    shortening_version text not null unique,
    usr_id      bigint
        constraint fk_url_usr
            references usr
            on update cascade,

    zlins_dttm    timestamp with time zone,
    zlupd_dttm    timestamp with time zone
);



create table if not exists logjwt
(
    id bigint generated by default as identity
        constraint logjwt_pkey
            primary key,
    dttm       timestamp with time zone not null,
    usr_id      bigint
        constraint fk_logjwt_usr
            references usr
            on update cascade,
    jwt        text                     not null,
    expires_on timestamp with time zone,
    is_invalid boolean default false    not null

);

create index ndx_logjwt_usr
    on logjwt (usr_id);






--
-- PostgreSQL database dump complete
--