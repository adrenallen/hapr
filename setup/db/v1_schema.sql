CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY unique primary key,
    username VARCHAR NOT NULL unique,
    password CHAR(60) NOT NULL,
    email varchar not null unique
);

CREATE TABLE user_changelog_seen(
    id INT GENERATED ALWAYS AS IDENTITY unique primary key,
    user_id INT REFERENCES users (id) NOT NULL,
    version_string varchar not null,
    unique(user_id)
);

CREATE TABLE factor_types(
    id int unique primary key,
    factor_type VARCHAR NOT NULL unique
);

CREATE TABLE factors(
    id INT GENERATED ALWAYS AS IDENTITY unique primary key,
    user_id INT REFERENCES users (id) NOT NULL,
    factor VARCHAR NOT NULL,
    unique(factor, user_id)
);

CREATE TABLE ratings(
    id INT GENERATED ALWAYS AS IDENTITY unique primary key,
    user_id INT REFERENCES users (id) NOT NULL,
    rating smallint NOT NULL,
    created_datetime timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE rating_factors(
    id INT GENERATED ALWAYS AS IDENTITY unique primary key,
    rating_id INT REFERENCES ratings (id) NOT NULL,
    factor_id INT REFERENCES factors (id) NOT NULL,
    factor_type_id INT REFERENCES factor_types (id) NOT NULL,
    rank smallint NOT NULL,
    unique(rating_id, factor_id, factor_type_id)
--    ,unique(rating_id, factor_type_id, rank)
);

CREATE TABLE session_guids(
  id INT GENERATED ALWAYS AS IDENTITY unique primary key,
  user_id INT REFERENCES users (id) NOT NULL,
  guid CHAR(36) NOT NULL,
  created_datetime timestamp with time zone NOT NULL DEFAULT now(),
  active BOOLEAN,
  unique(guid, active)
);

CREATE TABLE password_resets (
    id INT GENERATED ALWAYS AS IDENTITY unique primary key,
    user_id INT REFERENCES users (id) NOT NULL,
    created_datetime timestamp with time zone NOT NULL DEFAULT now(),
    token char(36) NOT NULL
);



alter table factors add column factor_encrypted BYTEA;
update factors set factor_encrypted=encrypt(factor::bytea, 'Aw5GCs6plC'::bytea, 'aes'::text);
alter table factors drop column factor;
alter table factors rename column factor_encrypted to factor;
alter table factors alter column factor set not null;
alter table factors add constraint UQ_factors_factor_user_id unique(factor, user_id);
alter table factors add  column archived bool default false;

alter table
	rating_factors drop
		constraint rating_factors_rating_id_factor_id_factor_type_id_key;

create
	table
		factor_aspects( id int generated always as identity unique primary key,
		factor_id int references factors(id) not null,
		factor_aspect BYTEA not null,
		archived bool not null default false,
		unique(factor_id,
		factor_aspect) );

alter table
	rating_factors add column factor_aspect_id int references factor_aspects(id);

alter table
	rating_factors add constraint UQ_rating_factors_rating_id_factor_id_factory_type_id_factor_aspect_id unique(rating_id,
	factor_id,
	factor_type_id,
	factor_aspect_id);