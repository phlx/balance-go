request

    GET /v1/transactions?user_id=123&sort[time]=desc&limit=2&page=8

response

    200 OK
    {
      "transactions": [
        {id:"5",time:"2020-09-19T02:27:08.572","amount":-5,"from_user_id":null,"reason":"Service charges"},
        {id:"4",time:"2020-09-19T02:27:08.482","amount":5,"from_user_id":321,"reason":"Owed you"}
      ],
      "pages": 9
    }
    
    400 Bad Request
    {"error":{"code":"incorrect_user","detail":"Incorrect value in from_user_id"}}
    
    500 Internal Server Error
    {"error":{"code":"internal_error","detail":"Error without staktrace"},"debug":"<stacktrace if debug allowed>"}
