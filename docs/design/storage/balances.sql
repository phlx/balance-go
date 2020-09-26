create table balances
(
	user_id bigserial not null
		constraint balances_pkey
			primary key,
	balance numeric(20,2) not null,
	updated_at timestamp with time zone
);
