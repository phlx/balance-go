request

    POST /v1/take
    
    {"user_id":123,"amount":35.45,"reason":"Service charges","idempotence_key":"56d6a679-d4d0-4c8a-bfde-6c620928bb96"}
    
response

    200 OK
    {
        "transaction": 25,
        "time": "2020-09-20T17:48:42.493"
    }
    
    400 Bad Request
    {"error":{"code":"insufficient_funds","detail":"User has insufficient funds"}}
    
    500 Internal Server Error
    {"error":{"code":"internal_error","detail":"Error without staktrace"},"debug":"<stacktrace if debug allowed>"}
