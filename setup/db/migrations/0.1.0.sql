begin transaction;

CREATE TABLE report_types(
    id int unique primary key,
    report_type VARCHAR NOT NULL unique
);

create table report_requests(
	id int generated always as identity unique primary key,
	user_id int references users(id),
	report_type_id int references report_types(id),
	requested_datetime timestamp with time zone NOT NULL DEFAULT now(),
	additional_parameters jsonb,
	completed_datetime timestamp with time zone
);

create table report_results(
	id int generated always as identity unique primary key,
	user_id int references users(id),
	report_type_id int references report_types(id),
	report_request_id int references report_requests(id),
	completed_datetime timestamp with time zone NOT NULL DEFAULT now(),
	result jsonb not NULL
);

insert into report_types values (1, 'Factor Impact');


commit;