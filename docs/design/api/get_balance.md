request

    GET /v1/balance?user_id=123&currency=USD
    
response

    200 OK
    {"balance":1.23,"currency":"USD"}
    
    400 Bad Request
    {"errors":[{"code":"incorrect_user","detail":"Incorrect value in from_user_id"}]}
    
    500 Internal Server Error
    {"errors":[{"code":"internal_error","detail":"Error without staktrace"},"debug":"<stacktrace if debug allowed>"]}
