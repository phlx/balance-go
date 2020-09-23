#!/bin/bash
for i in {1..10}; do
  sleep 0.1 && curl -X POST 'http://localhost:3000/v1/take' \
    -H "X-Idempotency-Key: $(uuidgen)" -H 'Content-Type: application/json' -H 'X-Sleep: 1000' \
    --data-raw '{"user_id":99999999,"amount":50,"reason":"taken"}' && echo ""
done
