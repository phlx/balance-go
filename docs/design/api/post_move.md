request

    POST /v1/move

    {
        "from_user_id":123,
        "to_user_id":321,
        "amount":35.45,
        "reason":"fun",
        "idempotence_key":"56d6a679-d4d0-4c8a-bfde-6c620928bb96"
    }

response

    200 OK
    {
        "transaction_from_id": 22,
        "transaction_from_time": "2020-09-20T13:47:41.2050577+03:00",
        "transaction_to_id": 23,
        "transaction_to_time": "2020-09-20T13:47:41.2060457+03:00"
    }

    400 Bad Request
    {"error":{"code":"incorrect_user","detail":"Incorrect value in from_user_id"}}

    500 Internal Server Error
    {"error":{"code":"internal_error","detail":"Error without staktrace"},"debug":"<stacktrace if debug allowed>"}
