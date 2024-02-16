# Plan

- First complete the exercise - only write tests after. This is because the exercise has incremental things to do, so what do we test for until it isn't completed (?)

### Testing

- Should be easier to test compared to client, as we are straightaway testing handlers, which are simply functions
- What all would we like to test?
    - random requests don't pattern match + return 404 + correct 404 message
    - for 200, 500 handlers
        - correct headers, if any specific ones are set
        - response status code
    - for base handler
        - correct headers
        - response status code
        - returns correct data i.e. html + it is escaped
    - handle authenticated
        - correct headers in all cases
        - no creds gives a 401 and a www-authenticated header
        - invalid creds fail and gives 401
        - valid creds give 200 + the correct html + correct content type header
- Mock request:
    - Perhaps NewRequest() in HTTP test could help?
- Mock w http.ResponseWriter:
    - Has only 3 methods: Write(), Header() and WriteHeader(), I think this can be mocked quickly + easily?
    - Actually - ResponseRecorder in net/httptest seems perfect for this purpose!
        - Create an httptest.NewRequest()
        - Create a new recorder using httptest.NewRecorder()
        - Pass both to handler
        - get http.Response using Result() method
        - inspect this and match with all our conditions defined above!!


### Thoughts post impl

- Exercise to re-send post request html content. Is using two buffers the way to go? (Answer: No, something like io.Copy makes things much simpler)
- Exercise to bake-in query params is interesting in terms of escaping the html. Teaches how to use query params and also revises the concepts of buffers and writers (interfaces in general)
- Implemented test for 200 handler, for 500 and 404, will look very very similar. Is my usage of a buffer and io.Copy() correct in the test?
