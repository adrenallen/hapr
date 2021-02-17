
create table interested_people_emails (
	id int generated always as identity unique primary key,
	email varchar not null unique,
	created_datetime timestamp with time zone NOT NULL DEFAULT now()
);