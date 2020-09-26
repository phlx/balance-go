request

    POST /v1/give
    
    {"user_id":123,"amount":35.45,"initiator_id":321,"reason":"The Great Bonus","idempotence_key":"56d6a679-d4d0-4c8a-bfde-6c620928bb96"}

response

    200 OK
    {
        "transaction": 24,
        "time": "2020-09-20T17:46:42.272"
    }
    
    400 Bad Request
    {"error":{"code":"incorrect_amount","detail":"Amount must be numeric value greater than zero"}}
    
    500 Internal Server Error
    {"error":{"code":"internal_error","detail":"Error without staktrace"},"debug":"<stacktrace if debug allowed>"}
