#!/bin/bash
curl -X POST 'http://localhost:3000/v1/give' \
  -H "X-Idempotency-Key: $(uuidgen)" -H 'Content-Type: application/json' \
  --data-raw '{"user_id":99999999,"amount":100,"reason":"given"}'
echo ""
