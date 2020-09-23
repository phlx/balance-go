CREATE TABLE IF NOT EXISTS "balances"
(
    "id"         bigserial,
    "user_id"    bigint           NOT NULL UNIQUE,
    "balance"    double precision NOT NULL,
    "updated_at" timestamptz,
    PRIMARY KEY ("id"),
    UNIQUE ("user_id")
);
