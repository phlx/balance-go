transactions

    id              serial bigint
    user_id         bigint
    amount          double precision
    initiator_id    bigint null
    reason          text
    created_at      timestamp
