transactions

    id              serial bigint
    user_id         bigint
    amount          numeric(20,2)
    initiator_id    bigint null
    reason          text
    created_at      timestamp
