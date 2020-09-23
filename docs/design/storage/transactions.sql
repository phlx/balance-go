CREATE TABLE IF NOT EXISTS "transactions"
(
    "id"           bigserial,
    "user_id"      bigint           NOT NULL,
    "amount"       double precision NOT NULL,
    "initiator_id" bigint,
    "reason"       text,
    "created_at"   timestamptz,
    PRIMARY KEY ("id")
);
