create table transactions
(
	id bigserial not null
		constraint transactions_pkey
			primary key,
	user_id bigint not null,
	amount numeric(20,2) not null,
	initiator_id bigint,
	reason text,
	created_at timestamp with time zone
);
