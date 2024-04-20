# Implementation

Will just be describing the implementation details here since I feel I am fairly comfortable with the material.


## Server serving in-memory JSON

- Simple server with a single endpoint
- Query params:
    - To determine indent
- To serve we'd need to Marshal data into JSON, instead of using marshal we can use the Encoder
    - Initialize a new encoder as NewEncoder(w http.ResponseWriter)
    - encoder.Encode((our struct here))
    - we can use struct tags to make sure the data is in snake_case
    - for indentation, we can use encode.SetIndent() depending on the query param


## Server serving Postgres

- Mostly involves understanding how pgx works, it has an iterator style interface for fetching data
- Once that is done, use a post request, use a decoder to get the underlying JSON, then use pgx again to insert it



## Overaching ideas

- Working with JSON data
    - struct tags, which is a Go thing
    - the idea of marshal/unmarshal encode/decode which are on a high-level ideas useful for working with any data type - so if you're working with things other than JSON, you will have a similar interface
- Working with env variables for sensitive/changeable things
- Error handling + controlling exposure of client to internal errors
- Contexts (lesson didn't dive into it but encouraged exploring)
- Databases and how talking to them is part of a very common developer workflow
    - for which you usually use an ORM, so you'll usually have to understand the ORM's interface 


