create table feature_flags (
    id int generated always as identity unique primary key,
    description varchar not null unique,
    enabled bool not null default false
);

create table early_access_codes(
    id int generated always as identity unique primary key,
    created_datetime timestamp with time zone NOT NULL DEFAULT now(),
    code char(6) not null,
    claimed_by_user_id INT REFERENCES users (id),
    claimed_datetime timestamp with time zone,
    emailed_to varchar
);

alter table ratings add column journal_entry text not null default '';

create table firebase_messaging_devices (
	id int generated always as identity unique primary key,
	device_id varchar(255) not null unique,
	user_id int references users(id)
);


create table user_reminder_settings(
    id int generated always as identity unique primary key,
	user_id int references users(id),
    reminder_time smallint not null,
    sunday boolean not null default false,
    monday boolean not null default false,
    tuesday boolean not null default false,
    wednesday boolean not null default false,
    thursday boolean not null default false,
    friday boolean not null default false,
    saturday boolean not null default false,
    email boolean not null default false,
    push boolean not null default false,
    sms boolean not null default false
);

alter table users add column mobile_number varchar(11);