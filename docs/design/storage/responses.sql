create table responses
(
	idempotency_key varchar(255) not null
		constraint responses_pkey
			primary key,
	status bigint not null,
	headers text not null,
	response text not null,
	created_at timestamp with time zone
);
